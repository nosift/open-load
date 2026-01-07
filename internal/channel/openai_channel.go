package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	app_errors "gpt-load/internal/errors"
	"gpt-load/internal/models"
	"gpt-load/internal/utils"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	Register("openai", newOpenAIChannel)
}

type OpenAIChannel struct {
	*BaseChannel
}

func newOpenAIChannel(f *Factory, group *models.Group) (ChannelProxy, error) {
	base, err := f.newBaseChannel("openai", group)
	if err != nil {
		return nil, err
	}

	return &OpenAIChannel{
		BaseChannel: base,
	}, nil
}

func isOpenRouterURL(u *url.URL) bool {
	if u == nil {
		return false
	}
	return strings.Contains(strings.ToLower(u.Hostname()), "openrouter.ai")
}

func applyOpenRouterHeaders(req *http.Request, group *models.Group) {
	if req == nil || req.URL == nil || !isOpenRouterURL(req.URL) {
		return
	}

	// OpenRouter recommends these headers and may enforce them for some (e.g. free) routes.
	if req.Header.Get("HTTP-Referer") == "" {
		referer := req.Header.Get("Origin")
		if referer == "" {
			referer = req.Header.Get("Referer")
		}
		if referer == "" && group != nil {
			referer = group.EffectiveConfig.AppUrl
		}
		if referer != "" {
			req.Header.Set("HTTP-Referer", referer)
		}
	}

	if req.Header.Get("X-Title") == "" && group != nil {
		title := group.DisplayName
		if title == "" {
			title = group.Name
		}
		if title != "" {
			req.Header.Set("X-Title", title)
		}
	}
}

// ModifyRequest sets the Authorization header for the OpenAI service.
func (ch *OpenAIChannel) ModifyRequest(req *http.Request, apiKey *models.APIKey, group *models.Group) {
	req.Header.Set("Authorization", "Bearer "+apiKey.KeyValue)
	applyOpenRouterHeaders(req, group)
}

// IsStreamRequest checks if the request is for a streaming response using the pre-read body.
func (ch *OpenAIChannel) IsStreamRequest(c *gin.Context, bodyBytes []byte) bool {
	if strings.Contains(c.GetHeader("Accept"), "text/event-stream") {
		return true
	}

	if c.Query("stream") == "true" {
		return true
	}

	type streamPayload struct {
		Stream bool `json:"stream"`
	}
	var p streamPayload
	if err := json.Unmarshal(bodyBytes, &p); err == nil {
		return p.Stream
	}

	return false
}

func (ch *OpenAIChannel) ExtractModel(c *gin.Context, bodyBytes []byte) string {
	type modelPayload struct {
		Model string `json:"model"`
	}
	var p modelPayload
	if err := json.Unmarshal(bodyBytes, &p); err == nil {
		return p.Model
	}
	return ""
}

// isOfficialOpenAI checks if the group is using official OpenAI API
func (ch *OpenAIChannel) isOfficialOpenAI(group *models.Group) bool {
	upstreamURL := ch.getUpstreamURL()
	if upstreamURL == nil {
		return false
	}

	// Check if the hostname is official OpenAI domain
	hostname := strings.ToLower(upstreamURL.Hostname())
	return strings.Contains(hostname, "api.openai.com") ||
	       strings.Contains(hostname, "openai.azure.com")
}

// ValidateKey checks if the given API key is valid by making a chat completion request.
func (ch *OpenAIChannel) ValidateKey(ctx context.Context, apiKey *models.APIKey, group *models.Group) (bool, error) {
	upstreamURL := ch.getUpstreamURL()
	if upstreamURL == nil {
		return false, fmt.Errorf("no upstream URL configured for channel %s", ch.Name)
	}

	// Parse validation endpoint to extract path and query parameters
	endpointURL, err := url.Parse(ch.ValidationEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to parse validation endpoint: %w", err)
	}

	// Build final URL with path and query parameters
	finalURL := *upstreamURL
	endpointPath := endpointURL.Path
	// Compatibility: avoid duplicated API version prefix when the upstream base URL already includes it.
	// e.g. base=https://host/api/v1 + endpoint=/v1/chat/completions => /api/v1/chat/completions
	basePath := strings.TrimRight(finalURL.Path, "/")
	if strings.HasSuffix(basePath, "/v1") && strings.HasPrefix(endpointPath, "/v1") {
		endpointPath = strings.TrimPrefix(endpointPath, "/v1")
		if endpointPath == "" {
			endpointPath = "/"
		}
	}
	finalURL.Path = strings.TrimRight(finalURL.Path, "/") + endpointPath
	finalURL.RawQuery = endpointURL.RawQuery
	reqURL := finalURL.String()

	// Use a minimal, low-cost payload for validation
	payload := gin.H{
		"model": ch.TestModel,
		"messages": []gin.H{
			{"role": "user", "content": "hi"},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal validation payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("failed to create validation request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey.KeyValue)
	req.Header.Set("Content-Type", "application/json")
	applyOpenRouterHeaders(req, group)

	// Apply custom header rules if available
	if len(group.HeaderRuleList) > 0 {
		headerCtx := utils.NewHeaderVariableContext(group, apiKey)
		utils.ApplyHeaderRules(req, group.HeaderRuleList, headerCtx)
	}

	resp, err := ch.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send validation request: %w", err)
	}
	defer resp.Body.Close()

	// Note: Organization verification status is NOT determined during key validation.
	// It's only marked as verified when an actual API call to a premium model succeeds.
	// This is because validation uses a test model, which doesn't require org verification.

	// Any 2xx status code indicates the key is valid.
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	// For non-200 responses, parse the body to provide a more specific error reason.
	errorBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("key is invalid (status %d), but failed to read error body: %w", resp.StatusCode, err)
	}

	// Use the new parser to extract a clean error message.
	parsedError := app_errors.ParseUpstreamError(errorBody)

	return false, fmt.Errorf("[status %d] %s", resp.StatusCode, parsedError)
}

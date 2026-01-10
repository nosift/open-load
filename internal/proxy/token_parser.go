package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// TokenUsage represents token usage from API response
type TokenUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

// UsageResponse represents common API response with usage
type UsageResponse struct {
	Usage *TokenUsage `json:"usage"`
}

// parseTokensFromResponse parses token usage from a non-streaming response
func parseTokensFromResponse(body []byte) *TokenUsage {
	var resp UsageResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil
	}
	return resp.Usage
}

// parseTokensFromStreamChunk parses token usage from a stream chunk (SSE format)
func parseTokensFromStreamChunk(data []byte) *TokenUsage {
	// SSE format: data: {...}\n\n
	dataStr := string(data)
	lines := strings.Split(dataStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		jsonStr := strings.TrimPrefix(line, "data:")
		jsonStr = strings.TrimSpace(jsonStr)

		// Skip [DONE] marker
		if jsonStr == "[DONE]" {
			continue
		}

		var resp UsageResponse
		if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
			continue
		}

		if resp.Usage != nil {
			return resp.Usage
		}
	}

	return nil
}

func (ps *ProxyServer) handleStreamingResponseWithTokens(c *gin.Context, resp *http.Response) *TokenUsage {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		logrus.Error("Streaming unsupported by the writer, falling back to normal response")
		return ps.handleNormalResponseWithTokens(c, resp)
	}

	var lastUsage *TokenUsage
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			// Try to parse tokens from each chunk
			if usage := parseTokensFromStreamChunk(line); usage != nil {
				lastUsage = usage
			}

			if _, writeErr := c.Writer.Write(line); writeErr != nil {
				logUpstreamError("writing stream to client", writeErr)
				return lastUsage
			}
			flusher.Flush()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			logUpstreamError("reading from upstream", err)
			return lastUsage
		}
	}

	return lastUsage
}

func (ps *ProxyServer) handleNormalResponseWithTokens(c *gin.Context, resp *http.Response) *TokenUsage {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logUpstreamError("reading response body", err)
		return nil
	}

	// Parse tokens from response
	usage := parseTokensFromResponse(body)

	// Write body to client
	if _, err := io.Copy(c.Writer, bytes.NewReader(body)); err != nil {
		logUpstreamError("copying response body", err)
	}

	return usage
}

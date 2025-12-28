export type ModelBadge = {
  label: string;
  short: string;
  className: string;
};

export const getModelBadge = (model: string, fallbackLabel = "Model"): ModelBadge => {
  const normalized = (model || "").toLowerCase();

  if (
    normalized.includes("gpt") ||
    normalized.includes("openai") ||
    normalized.startsWith("o1") ||
    normalized.startsWith("o3") ||
    normalized.includes("text-embedding") ||
    normalized.includes("whisper")
  ) {
    return { label: "OpenAI", short: "OA", className: "model-badge-openai" };
  }

  if (normalized.includes("claude") || normalized.includes("anthropic")) {
    return { label: "Anthropic", short: "A", className: "model-badge-anthropic" };
  }

  if (normalized.includes("gemini") || normalized.includes("palm")) {
    return { label: "Gemini", short: "G", className: "model-badge-gemini" };
  }

  if (normalized.includes("deepseek")) {
    return { label: "DeepSeek", short: "D", className: "model-badge-deepseek" };
  }

  if (normalized.includes("qwen")) {
    return { label: "Qwen", short: "Q", className: "model-badge-qwen" };
  }

  return { label: fallbackLabel, short: "M", className: "model-badge-default" };
};

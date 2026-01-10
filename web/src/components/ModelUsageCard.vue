<script setup lang="ts">
import type { DashboardStatsResponse, ModelUsageItem } from "@/types/models";
import { getModelBadge } from "@/utils/model-badge";
import { NCard, NEmpty, NTooltip } from "naive-ui";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

interface Props {
  stats: DashboardStatsResponse | null;
}

const props = defineProps<Props>();

const range = ref<"24h" | "7d">("24h");

const usage24h = computed(() => props.stats?.model_usage_24h ?? []);
const usage7d = computed(() => props.stats?.model_usage_7d ?? []);
const usage = computed(() => (range.value === "24h" ? usage24h.value : usage7d.value));
const hasStats = computed(() => props.stats !== null);

const formatCount = (value: number | undefined) => {
  if (!value) return "0";
  if (value >= 1000000) {
    return `${(value / 1000000).toFixed(1)}M`;
  }
  if (value >= 1000) {
    return `${(value / 1000).toFixed(1)}K`;
  }
  return value.toString();
};

const getDisplayName = (item: ModelUsageItem) => {
  if (item.model && item.model.trim()) {
    return item.model;
  }
  return t("dashboard.unknownModel");
};
</script>

<template>
  <n-card :bordered="false" class="model-usage-card">
    <div class="model-usage-header">
      <div>
        <div class="model-usage-title">{{ t("dashboard.modelUsage") }}</div>
        <div class="model-usage-subtitle">{{ t("dashboard.requestStatistics") }}</div>
      </div>
      <div class="header-controls">
        <div class="cf-toggle-container">
          <button
            class="cf-toggle-btn"
            :class="{ active: range === '24h' }"
            @click="range = '24h'"
          >
            {{ t("dashboard.modelUsage24h") }}
          </button>
          <div class="cf-divider"></div>
          <button
            class="cf-toggle-btn"
            :class="{ active: range === '7d' }"
            @click="range = '7d'"
          >
            {{ t("dashboard.modelUsage7d") }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="usage.length" class="model-usage-grid">
      <div
        v-for="(item, index) in usage"
        :key="`${item.model || 'unknown'}-${index}`"
        class="model-usage-item"
      >
        <span
          :class="['model-badge', getModelBadge(item.model, t('dashboard.modelUsage')).className]"
          :title="getModelBadge(item.model, t('dashboard.modelUsage')).label"
        >
          {{ getModelBadge(item.model, t("dashboard.modelUsage")).short }}
        </span>
        <div class="model-usage-info">
          <div class="model-usage-name">{{ getDisplayName(item) }}</div>
          <div class="model-usage-stats">
            <n-tooltip trigger="hover">
              <template #trigger>
                <span class="stat-item">
                  <span class="stat-value">{{ formatCount(item.request_count) }}</span>
                  <span class="stat-label">{{ t("dashboard.requests") }}</span>
                </span>
              </template>
              {{ t("dashboard.requests") }}: {{ item.request_count?.toLocaleString() ?? 0 }}
            </n-tooltip>
            <span class="stat-divider">|</span>
            <n-tooltip trigger="hover">
              <template #trigger>
                <span class="stat-item">
                  <span class="stat-value">{{ formatCount(item.total_tokens) }}</span>
                  <span class="stat-label">Tokens</span>
                </span>
              </template>
              {{ t("dashboard.promptTokens") }}: {{ item.prompt_tokens?.toLocaleString() ?? 0 }}<br>
              {{ t("dashboard.completionTokens") }}: {{ item.completion_tokens?.toLocaleString() ?? 0 }}<br>
              {{ t("dashboard.totalTokens") }}: {{ item.total_tokens?.toLocaleString() ?? 0 }}
            </n-tooltip>
          </div>
        </div>
      </div>
    </div>
    <n-empty v-else-if="hasStats" :description="t('dashboard.noModelData')" />
    <div v-else class="model-usage-loading">{{ t("common.loading") }}</div>
  </n-card>
</template>

<style scoped>
.model-usage-card {
  background: var(--card-bg);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(16px);
}

.model-usage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.header-controls {
  display: flex;
  align-items: center;
}

.cf-toggle-container {
  display: flex;
  align-items: center;
  border: 1px solid var(--border-color);
  border-radius: 999px;
  background: var(--bg-tertiary);
  padding: 4px;
  gap: 6px;
}

:root.dark .cf-toggle-container {
  background: rgba(255, 255, 255, 0.06);
  border-color: var(--border-color);
}

.cf-toggle-btn {
  padding: 6px 14px;
  font-size: 12px;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  cursor: pointer;
  border-radius: 999px;
  transition: all 0.2s ease;
  font-weight: 600;
}

.cf-toggle-btn:hover {
  color: var(--text-primary);
  background: rgba(0, 0, 0, 0.04);
}

:root.dark .cf-toggle-btn:hover {
  background: rgba(255, 255, 255, 0.12);
}

.cf-toggle-btn.active {
  background: var(--card-bg-solid);
  color: var(--text-primary);
  font-weight: 600;
  box-shadow: var(--shadow-sm);
}

:root.dark .cf-toggle-btn.active {
  background: rgba(255, 255, 255, 0.18);
  color: var(--text-primary);
}

.cf-divider {
  display: none;
}

:root.dark .cf-divider {
  background: rgba(255, 255, 255, 0.2);
}

.model-usage-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-primary);
}

.model-usage-subtitle {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.model-usage-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.model-usage-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 80px;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.model-usage-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color-light);
  border-radius: var(--border-radius-lg);
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
  box-shadow: var(--shadow-sm);
}

.model-usage-item:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--border-hover);
}

.model-usage-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
  flex: 1;
}

.model-usage-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 0.9rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.model-usage-stats {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.8rem;
}

.stat-item {
  display: flex;
  align-items: baseline;
  gap: 4px;
  cursor: help;
}

.stat-value {
  font-weight: 700;
  color: var(--text-primary);
  font-size: 0.95rem;
  font-feature-settings: "tnum";
}

.stat-label {
  color: var(--text-secondary);
  font-size: 0.75rem;
}

.stat-divider {
  color: var(--border-color);
  font-size: 0.7rem;
}

.model-badge {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  color: white;
  letter-spacing: 0.2px;
  flex-shrink: 0;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.16);
}

.model-badge-openai {
  background: linear-gradient(135deg, #111827 0%, #374151 100%);
}

.model-badge-gemini {
  background: linear-gradient(135deg, #1d4ed8 0%, #0ea5e9 100%);
}

.model-badge-anthropic {
  background: linear-gradient(135deg, #d97706 0%, #f59e0b 100%);
}

.model-badge-deepseek {
  background: linear-gradient(135deg, #0ea5e9 0%, #22c55e 100%);
}

.model-badge-qwen {
  background: linear-gradient(135deg, #ef4444 0%, #f97316 100%);
}

.model-badge-default {
  background: linear-gradient(135deg, #64748b 0%, #94a3b8 100%);
}

@media (max-width: 768px) {
  .model-usage-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>

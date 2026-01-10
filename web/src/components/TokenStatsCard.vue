<script setup lang="ts">
import type { DashboardStatsResponse } from "@/types/models";
import { NCard, NGrid, NGridItem, NSpin } from "naive-ui";
import { computed } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

interface Props {
  stats: DashboardStatsResponse | null;
  loading?: boolean;
}

const props = defineProps<Props>();

const tokenStats = computed(() => props.stats?.token_stats);

const formatTokens = (value: number | undefined) => {
  if (!value) return "0";
  if (value >= 1000000000) {
    return `${(value / 1000000000).toFixed(2)}B`;
  }
  if (value >= 1000000) {
    return `${(value / 1000000).toFixed(2)}M`;
  }
  if (value >= 1000) {
    return `${(value / 1000).toFixed(1)}K`;
  }
  return value.toString();
};

// 使用选定日期范围的数据 (24h 字段现在代表选定范围)
const promptTokens = computed(() => tokenStats.value?.prompt_tokens_24h);
const completionTokens = computed(() => tokenStats.value?.completion_tokens_24h);
const totalTokens = computed(() => tokenStats.value?.total_tokens_24h);
</script>

<template>
  <n-card :bordered="false" class="token-stats-card">
    <n-spin :show="loading">
      <div class="token-stats-header">
        <div>
          <div class="token-stats-title">{{ t("dashboard.tokenUsage") }}</div>
          <div class="token-stats-subtitle">{{ t("dashboard.tokenStatistics") }}</div>
        </div>
        <div class="header-controls">
          <slot name="header-extra"></slot>
        </div>
      </div>

      <n-grid cols="3" :x-gap="16" :y-gap="16">
        <n-grid-item>
          <div class="token-stat-item">
            <div class="token-stat-icon prompt-icon">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M12 2L2 7l10 5 10-5-10-5z" fill="none" />
                <path d="M2 17l10 5 10-5" fill="none" />
                <path d="M2 12l10 5 10-5" fill="none" />
              </svg>
            </div>
            <div class="token-stat-info">
              <div class="token-stat-value">{{ formatTokens(promptTokens) }}</div>
              <div class="token-stat-label">{{ t("dashboard.promptTokens") }}</div>
            </div>
          </div>
        </n-grid-item>

        <n-grid-item>
          <div class="token-stat-item">
            <div class="token-stat-icon completion-icon">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" fill="none" />
              </svg>
            </div>
            <div class="token-stat-info">
              <div class="token-stat-value">{{ formatTokens(completionTokens) }}</div>
              <div class="token-stat-label">{{ t("dashboard.completionTokens") }}</div>
            </div>
          </div>
        </n-grid-item>

        <n-grid-item>
          <div class="token-stat-item">
            <div class="token-stat-icon total-icon">
              <svg viewBox="0 0 24 24" aria-hidden="true">
                <circle cx="12" cy="12" r="10" fill="none" />
                <path d="M12 6v6l4 2" fill="none" />
              </svg>
            </div>
            <div class="token-stat-info">
              <div class="token-stat-value">{{ formatTokens(totalTokens) }}</div>
              <div class="token-stat-label">{{ t("dashboard.totalTokens") }}</div>
            </div>
          </div>
        </n-grid-item>
      </n-grid>
    </n-spin>
  </n-card>
</template>

<style scoped>
.token-stats-card {
  background: var(--card-bg);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(16px);
}

.token-stats-header {
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

.token-stats-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--text-primary);
}

.token-stats-subtitle {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.token-stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color-light);
  border-radius: var(--border-radius-lg);
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
  box-shadow: var(--shadow-sm);
}

.token-stat-item:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--border-hover);
}

.token-stat-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--border-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-primary);
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color-light);
  flex-shrink: 0;
}

.token-stat-icon svg {
  width: 20px;
  height: 20px;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.token-stat-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.token-stat-value {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
  font-feature-settings: "tnum";
}

.token-stat-label {
  font-size: 0.8rem;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .token-stats-header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>

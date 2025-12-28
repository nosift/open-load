<script setup lang="ts">
import { getDashboardStats } from "@/api/dashboard";
import type { DashboardStatsResponse, ModelUsageItem } from "@/types/models";
import { NDatePicker, NEmpty } from "naive-ui";
import { computed, ref, watch } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const stats = ref<DashboardStatsResponse | null>(null);
const dateRange = ref<[number, number]>([
  Date.now() - 6 * 24 * 60 * 60 * 1000,
  Date.now(),
]);

const usage = computed(() => stats.value?.model_usage_24h ?? []);
const sortedUsage = computed(() =>
  [...usage.value].sort((a, b) => b.request_count - a.request_count)
);
const totalRequests = computed(() =>
  sortedUsage.value.reduce((sum, item) => sum + item.request_count, 0)
);
const formatBeijingDate = (value: number) => {
  const parts = new Intl.DateTimeFormat("zh-CN", {
    timeZone: "Asia/Shanghai",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).formatToParts(new Date(value));
  const getPart = (type: string) => parts.find((part) => part.type === type)?.value ?? "";
  return `${getPart("year")}-${getPart("month")}-${getPart("day")}`;
};

const startDateLabel = computed(() => formatBeijingDate(dateRange.value[0]));
const endDateLabel = computed(() => formatBeijingDate(dateRange.value[1]));

const formatCount = (value: number) => {
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

const getPercent = (value: number) => {
  if (!totalRequests.value) {
    return 0;
  }
  return Math.round((value / totalRequests.value) * 1000) / 10;
};

const getCount = (value?: number) => value ?? 0;

const hoveredLabel = ref<string | null>(null);
const tooltipPos = ref({ x: 0, y: 0 });
const chartWrapRef = ref<HTMLElement | null>(null);

const palette = [
  "#22c1ff",
  "#b3ccff",
  "#b78bff",
  "#f6e36d",
  "#7c5cff",
  "#55d27e",
  "#f79b5d",
  "#f0637a",
];

const chartRadius = 86;
const chartStroke = 22;
const chartSize = 260;

const chartSegments = computed(() => {
  const total = totalRequests.value;
  if (!total) {
    return [];
  }
  const circumference = 2 * Math.PI * chartRadius;
  const gap = circumference * 0.01;
  let offset = 0;

  return sortedUsage.value.map((item, index) => {
    const ratio = item.request_count / total;
    const fullLength = ratio * circumference;
    const length = Math.max(0, fullLength - gap);
    const segment = {
      label: getDisplayName(item),
      value: item.request_count,
      percent: getPercent(item.request_count),
      success: getCount(item.success_count),
      failure: getCount(item.failure_count),
      retry: getCount(item.retry_count),
      color: palette[index % palette.length],
      dasharray: `${length} ${circumference - length}`,
      dashoffset: -offset,
    };
    offset += fullLength;
    return segment;
  });
});

const hoveredSegment = computed(() =>
  hoveredLabel.value
    ? chartSegments.value.find((segment) => segment.label === hoveredLabel.value) ?? null
    : null
);

const updateTooltipPos = (event: MouseEvent) => {
  const wrap = chartWrapRef.value;
  if (!wrap) {
    return;
  }
  const rect = wrap.getBoundingClientRect();
  tooltipPos.value = {
    x: event.clientX - rect.left,
    y: event.clientY - rect.top,
  };
};

const onSegmentEnter = (label: string, event: MouseEvent) => {
  hoveredLabel.value = label;
  updateTooltipPos(event);
};

const onSegmentMove = (event: MouseEvent) => {
  if (!hoveredLabel.value) {
    return;
  }
  updateTooltipPos(event);
};

const onSegmentLeave = () => {
  hoveredLabel.value = null;
};

const fetchStats = async () => {
  try {
    const params = {
      start: startDateLabel.value,
      end: endDateLabel.value,
      tz: "Asia/Shanghai",
    };
    const response = await getDashboardStats(params);
    stats.value = response.data;
  } catch (error) {
    console.error("Failed to load model usage:", error);
  }
};

watch(dateRange, fetchStats, { immediate: true });
</script>

<template>
  <div class="model-usage-page">
    <div class="analysis-shell">
      <div class="analysis-header">
        <div class="analysis-title">
          <span class="analysis-icon">
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <circle cx="12" cy="12" r="8" fill="none" />
              <path d="M12 12V4a8 8 0 0 1 8 8h-8z" fill="none" />
            </svg>
          </span>
          <span>{{ t("modelUsage.analysisTitle") }}</span>
        </div>
      </div>

      <div class="analysis-body">
        <div class="analysis-panel">
          <div class="panel-header">
            <div>
              <div class="panel-title">{{ t("modelUsage.requestShareTitle") }}</div>
              <div class="panel-subtitle">
                {{ t("modelUsage.total") }}: {{ formatCount(totalRequests) }}
              </div>
            </div>
            <div class="range-toggle">
              <span class="range-label">{{ t("modelUsage.timezoneLabel") }}</span>
              <n-date-picker
                v-model:value="dateRange"
                type="daterange"
                size="small"
                :clearable="false"
              />
            </div>
          </div>

          <div v-if="chartSegments.length" class="panel-content">
            <div class="legend">
              <div v-for="segment in chartSegments" :key="segment.label" class="legend-item">
                <span class="legend-dot" :style="{ background: segment.color }"></span>
                <span class="legend-name">{{ segment.label }}</span>
                <span class="legend-meta">
                  {{ formatCount(segment.value) }} · {{ segment.percent }}%
                </span>
              </div>
            </div>
            <div class="chart-wrap" ref="chartWrapRef">
              <div
                v-if="hoveredSegment"
                class="chart-tooltip"
                :style="{
                  left: `${tooltipPos.x}px`,
                  top: `${tooltipPos.y}px`,
                }"
              >
                <div class="tooltip-title">{{ hoveredSegment.label }}</div>
                <div class="tooltip-row">
                  <span class="tooltip-name">{{ t("modelUsage.successCount") }}</span>
                  <span class="tooltip-count">{{ formatCount(hoveredSegment.success) }}</span>
                </div>
                <div class="tooltip-row">
                  <span class="tooltip-name">{{ t("modelUsage.failureCount") }}</span>
                  <span class="tooltip-count">{{ formatCount(hoveredSegment.failure) }}</span>
                </div>
                <div class="tooltip-row">
                  <span class="tooltip-name">{{ t("modelUsage.retryCount") }}</span>
                  <span class="tooltip-count">{{ formatCount(hoveredSegment.retry) }}</span>
                </div>
              </div>
              <svg
                class="donut-chart"
                :width="chartSize"
                :height="chartSize"
                :viewBox="`0 0 ${chartSize} ${chartSize}`"
              >
                <circle
                  class="chart-track"
                  :cx="chartSize / 2"
                  :cy="chartSize / 2"
                  :r="chartRadius"
                  :stroke-width="chartStroke"
                />
                <g :transform="`rotate(-90 ${chartSize / 2} ${chartSize / 2})`">
                  <circle
                    v-for="segment in chartSegments"
                    :key="segment.label"
                    class="chart-segment"
                    :cx="chartSize / 2"
                    :cy="chartSize / 2"
                    :r="chartRadius"
                    :stroke="segment.color"
                    :stroke-width="chartStroke"
                    :stroke-dasharray="segment.dasharray"
                    :stroke-dashoffset="segment.dashoffset"
                    @mouseenter="onSegmentEnter(segment.label, $event)"
                    @mousemove="onSegmentMove"
                    @mouseleave="onSegmentLeave"
                  />
                </g>
              </svg>
            </div>
          </div>
          <n-empty v-else-if="stats" :description="t('dashboard.noModelData')" />
          <div v-else class="usage-loading">{{ t("common.loading") }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.model-usage-page {
  display: flex;
  flex-direction: column;
}

.analysis-shell {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(16px);
}

.analysis-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color-light);
}

.analysis-title {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--text-primary);
}

.analysis-icon {
  width: 28px;
  height: 28px;
  border-radius: 10px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color-light);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-primary);
}

.analysis-icon svg {
  width: 16px;
  height: 16px;
  stroke: currentColor;
  stroke-width: 2.2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.analysis-body {
  padding: 28px 32px 32px;
}

.analysis-panel {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.panel-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
}

.panel-title {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--text-primary);
}

.panel-subtitle {
  margin-top: 6px;
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.range-toggle {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 6px 10px;
  border-radius: 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
}

.range-label {
  color: var(--text-secondary);
  font-size: 0.85rem;
}

.panel-content {
  display: grid;
  grid-template-columns: minmax(220px, 280px) minmax(240px, 1fr);
  gap: 32px;
  align-items: center;
}

.legend {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-primary);
  font-size: 0.95rem;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  box-shadow: 0 0 0 4px rgba(0, 0, 0, 0.03);
}

.legend-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.legend-meta {
  margin-left: auto;
  font-size: 0.85rem;
  color: var(--text-secondary);
  white-space: nowrap;
}

.chart-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.donut-chart {
  overflow: visible;
}

.chart-track {
  fill: none;
  stroke: var(--bg-tertiary);
}

.chart-segment {
  fill: none;
  stroke-linecap: round;
  transition: stroke-dashoffset 0.3s ease, stroke-width 0.2s ease, transform 0.2s ease,
    filter 0.2s ease;
  transform-origin: center;
  transform-box: fill-box;
  cursor: pointer;
}

.chart-segment:hover {
  stroke-width: 26px;
  transform: scale(1.03);
  filter: drop-shadow(0 6px 12px rgba(0, 0, 0, 0.08));
}

.chart-tooltip {
  position: absolute;
  z-index: 2;
  padding: 14px 16px;
  min-width: 220px;
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: var(--shadow-md);
  pointer-events: none;
  transform: translate(16px, -50%);
}

.tooltip-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.tooltip-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  color: var(--text-secondary);
}

.tooltip-name {
  flex: 1;
  font-size: 0.95rem;
  color: var(--text-secondary);
}

.tooltip-count {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.usage-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 220px;
  color: var(--text-secondary);
}

@media (max-width: 900px) {
  .analysis-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .panel-content {
    grid-template-columns: 1fr;
  }

  .chart-wrap {
    justify-content: flex-start;
  }
}
</style>

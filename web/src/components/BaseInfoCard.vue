<script setup lang="ts">
import type { DashboardStatsResponse } from "@/types/models";
import { NCard, NGrid, NGridItem, NSpace, NTag, NTooltip } from "naive-ui";
import { computed, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

// Props
interface Props {
  stats: DashboardStatsResponse | null;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

// 使用计算属性代替ref
const stats = computed(() => props.stats);
const animatedValues = ref<Record<string, number>>({});

// 格式化数值显示
const formatValue = (value: number, type: "count" | "rate" = "count"): string => {
  if (type === "rate") {
    return `${value.toFixed(1)}%`;
  }
  if (value >= 1000) {
    return `${(value / 1000).toFixed(1)}K`;
  }
  return value.toString();
};

// 格式化趋势显示
const formatTrend = (trend: number): string => {
  const sign = trend >= 0 ? "+" : "";
  return `${sign}${trend.toFixed(1)}%`;
};

// 监听stats变化并更新动画值
const updateAnimatedValues = () => {
  if (stats.value) {
    setTimeout(() => {
      animatedValues.value = {
        key_count:
          (stats.value?.key_count?.value ?? 0) /
          ((stats.value?.key_count?.value ?? 1) + (stats.value?.key_count?.sub_value ?? 1)),
        rpm: Math.min(100 + (stats.value?.rpm?.trend ?? 0), 100) / 100,
        request_count: Math.min(100 + (stats.value?.request_count?.trend ?? 0), 100) / 100,
        error_rate: (100 - (stats.value?.error_rate?.value ?? 0)) / 100,
      };
    }, 0);
  }
};

// 监听stats变化
onMounted(() => {
  updateAnimatedValues();
});
</script>

<template>
  <div class="stats-container">
    <n-space vertical size="medium">
      <n-grid cols="2 s:4" :x-gap="20" :y-gap="20" responsive="screen">
        <!-- 密钥数量 -->
        <n-grid-item span="1">
          <n-card :bordered="false" class="stat-card" style="animation-delay: 0s">
            <div class="stat-header">
              <div class="stat-icon key-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <rect x="6" y="11" width="12" height="9" rx="2" fill="none" />
                  <path d="M8 11V7a4 4 0 0 1 8 0v4" fill="none" />
                </svg>
              </div>
              <n-tooltip v-if="stats?.key_count.sub_value" trigger="hover">
                <template #trigger>
                  <n-tag type="error" size="small" class="stat-trend">
                    {{ stats.key_count.sub_value }}
                  </n-tag>
                </template>
                {{ stats.key_count.sub_value_tip }}
              </n-tooltip>
            </div>

            <div class="stat-content">
              <div class="stat-value">
                {{ stats?.key_count?.value ?? 0 }}
              </div>
              <div class="stat-title">{{ t("dashboard.totalKeys") }}</div>
            </div>

            <div class="stat-bar">
              <div
                class="stat-bar-fill key-bar"
                :style="{
                  width: `${(animatedValues.key_count ?? 0) * 100}%`,
                }"
              />
            </div>
          </n-card>
        </n-grid-item>

        <!-- RPM (10分钟) -->
        <n-grid-item span="1">
          <n-card :bordered="false" class="stat-card" style="animation-delay: 0.05s">
            <div class="stat-header">
              <div class="stat-icon rpm-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <circle cx="12" cy="12" r="9" fill="none" />
                  <path d="M12 7v5l3 3" fill="none" />
                </svg>
              </div>
              <n-tag
                v-if="stats?.rpm && stats.rpm.trend !== undefined"
                :type="stats?.rpm.trend_is_growth ? 'success' : 'error'"
                size="small"
                class="stat-trend"
              >
                {{ stats ? formatTrend(stats.rpm.trend) : "--" }}
              </n-tag>
            </div>

            <div class="stat-content">
              <div class="stat-value">
                {{ stats?.rpm?.value.toFixed(1) ?? 0 }}
              </div>
              <div class="stat-title">{{ t("dashboard.rpm10Min") }}</div>
            </div>

            <div class="stat-bar">
              <div
                class="stat-bar-fill rpm-bar"
                :style="{
                  width: `${(animatedValues.rpm ?? 0) * 100}%`,
                }"
              />
            </div>
          </n-card>
        </n-grid-item>

        <!-- 24小时请求 -->
        <n-grid-item span="1">
          <n-card :bordered="false" class="stat-card" style="animation-delay: 0.1s">
            <div class="stat-header">
              <div class="stat-icon request-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M3 17l6-6 4 4 7-8" fill="none" />
                  <path d="M16 7h5v5" fill="none" />
                </svg>
              </div>
              <n-tag
                v-if="stats?.request_count && stats.request_count.trend !== undefined"
                :type="stats?.request_count.trend_is_growth ? 'success' : 'error'"
                size="small"
                class="stat-trend"
              >
                {{ stats ? formatTrend(stats.request_count.trend) : "--" }}
              </n-tag>
            </div>

            <div class="stat-content">
              <div class="stat-value">
                {{ stats ? formatValue(stats.request_count.value) : "--" }}
              </div>
              <div class="stat-title">{{ t("dashboard.requests24h") }}</div>
            </div>

            <div class="stat-bar">
              <div
                class="stat-bar-fill request-bar"
                :style="{
                  width: `${(animatedValues.request_count ?? 0) * 100}%`,
                }"
              />
            </div>
          </n-card>
        </n-grid-item>

        <!-- 24小时错误率 -->
        <n-grid-item span="1">
          <n-card :bordered="false" class="stat-card" style="animation-delay: 0.15s">
            <div class="stat-header">
              <div class="stat-icon error-icon">
                <svg viewBox="0 0 24 24" aria-hidden="true">
                  <path d="M12 3l7 4v5c0 5-3.5 9-7 10-3.5-1-7-5-7-10V7l7-4z" fill="none" />
                  <path d="M12 8v4" fill="none" />
                  <path d="M12 16h0" fill="none" />
                </svg>
              </div>
              <n-tag
                v-if="stats?.error_rate.trend !== 0"
                :type="stats?.error_rate.trend_is_growth ? 'success' : 'error'"
                size="small"
                class="stat-trend"
              >
                {{ stats ? formatTrend(stats.error_rate.trend) : "--" }}
              </n-tag>
            </div>

            <div class="stat-content">
              <div class="stat-value">
                {{ stats ? formatValue(stats.error_rate.value ?? 0, "rate") : "--" }}
              </div>
              <div class="stat-title">{{ t("dashboard.errorRate24h") }}</div>
            </div>

            <div class="stat-bar">
              <div
                class="stat-bar-fill error-bar"
                :style="{
                  width: `${(animatedValues.error_rate ?? 0) * 100}%`,
                }"
              />
            </div>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-space>
  </div>
</template>

<style scoped>
.stats-container {
  width: 100%;
  animation: fadeInUp 0.2s ease-out;
  margin-bottom: 16px;
}

.stat-card {
  background: var(--card-bg);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  position: relative;
  overflow: hidden;
  transition: all 0.2s ease;
  height: 100%; /* 确保卡片等高 */
  backdrop-filter: blur(16px);
}

.stat-card:hover {
  border-color: var(--border-hover);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px; /* 增加头部间距 */
}

/* 统一图标样式，改为单色/描边风格 */
.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--border-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-primary);
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color-light);
  box-shadow: inset 0 1px 1px rgba(255, 255, 255, 0.6);
}

.stat-icon svg {
  width: 18px;
  height: 18px;
  stroke: currentColor;
  stroke-width: 2.2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.key-icon,
.rpm-icon,
.request-icon,
.error-icon {
  /* 移除特定颜色背景，使用通用样式 */
  background: var(--bg-tertiary);
}

.stat-trend {
  font-weight: 500;
  /* 移除趋势标签前的箭头伪元素，因为NTag已经包含内容 */
}

.stat-trend:before {
  display: none;
}

.stat-content {
  margin-bottom: 16px;
}

.stat-value {
  font-size: 1.9rem;
  font-weight: 600;
  line-height: 1.2;
  color: var(--text-primary);
  margin-bottom: 4px;
  font-feature-settings: "tnum"; /* 表格数字对齐 */
}

.stat-title {
  font-size: 0.9rem;
  color: var(--text-secondary);
  font-weight: 400;
}

/* 简化进度条 */
.stat-bar {
  width: 100%;
  height: 3px;
  background: var(--border-color-light);
  border-radius: 1px;
  overflow: hidden;
  position: relative;
}

.stat-bar-fill {
  height: 100%;
  border-radius: 1px;
  transition: width 0.5s ease-out;
  background: var(--primary-color);
}

/* 移除特定进度条颜色，保持黑白灰 */
.key-bar, .rpm-bar, .request-bar, .error-bar {
  background: var(--primary-color);
}

/* 针对错误率进度条可以保留警示色，或者也统一 */
.error-bar {
  /* background: var(--error-color); */
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式网格 */
:deep(.n-grid-item) {
  min-width: 0;
}
</style>

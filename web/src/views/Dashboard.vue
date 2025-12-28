<script setup lang="ts">
import { getDashboardStats } from "@/api/dashboard";
import BaseInfoCard from "@/components/BaseInfoCard.vue";
import EncryptionMismatchAlert from "@/components/EncryptionMismatchAlert.vue";
import LineChart from "@/components/LineChart.vue";
import ModelUsageCard from "@/components/ModelUsageCard.vue";
import SecurityAlert from "@/components/SecurityAlert.vue";
import type { DashboardStatsResponse } from "@/types/models";
import { NSpace } from "naive-ui";
import { onMounted, ref } from "vue";

const dashboardStats = ref<DashboardStatsResponse | null>(null);

onMounted(async () => {
  try {
    const response = await getDashboardStats();
    dashboardStats.value = response.data;
  } catch (error) {
    console.error("Failed to load dashboard stats:", error);
  }
});
</script>

<template>
  <div class="dashboard-container">
    <n-space vertical size="large">
      <!-- 顶部状态栏：加密错误和安全警告 -->
      <n-space vertical>
        <encryption-mismatch-alert />
        
        <security-alert
          v-if="dashboardStats?.security_warnings"
          :warnings="dashboardStats.security_warnings"
        />
      </n-space>

      <!-- 核心指标卡片 -->
      <base-info-card :stats="dashboardStats" />
      
      <!-- 图表区域 -->
      <div class="dashboard-section chart-section">
        <h3 class="section-title">Request Trend</h3>
        <line-chart class="dashboard-chart" />
      </div>

      <!-- 模型使用情况 -->
      <div class="dashboard-section">
        <h3 class="section-title">Model Usage</h3>
        <model-usage-card :stats="dashboardStats" />
      </div>
    </n-space>
  </div>
</template>

<style scoped>
.dashboard-container {
  max-width: 100%;
}

.dashboard-section {
  margin-top: 8px;
}

.chart-section {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-lg);
  padding: 28px;
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(16px);
}

.section-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 16px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.dashboard-chart {
  /* 移除之前的动画，保持稳重 */
}
</style>

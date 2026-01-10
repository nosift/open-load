<script setup lang="ts">
import { getDashboardStats } from "@/api/dashboard";
import BaseInfoCard from "@/components/BaseInfoCard.vue";
import EncryptionMismatchAlert from "@/components/EncryptionMismatchAlert.vue";
import ModelUsageCard from "@/components/ModelUsageCard.vue";
import SecurityAlert from "@/components/SecurityAlert.vue";
import TokenStatsCard from "@/components/TokenStatsCard.vue";
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

      <!-- Token 统计 -->
      <token-stats-card :stats="dashboardStats" />

      <!-- 模型使用情况 -->
      <model-usage-card :stats="dashboardStats" />
    </n-space>
  </div>
</template>

<style scoped>
.dashboard-container {
  max-width: 100%;
}
</style>

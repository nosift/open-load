<script setup lang="ts">
import { getDashboardStats } from "@/api/dashboard";
import BaseInfoCard from "@/components/BaseInfoCard.vue";
import EncryptionMismatchAlert from "@/components/EncryptionMismatchAlert.vue";
import ModelUsageCard from "@/components/ModelUsageCard.vue";
import SecurityAlert from "@/components/SecurityAlert.vue";
import TokenStatsCard from "@/components/TokenStatsCard.vue";
import type { DashboardStatsResponse } from "@/types/models";
import { NDatePicker, NSpace } from "naive-ui";
import { onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const dashboardStats = ref<DashboardStatsResponse | null>(null);
const loading = ref(false);

// 日期范围：默认为今天
const getToday = () => {
  const now = new Date();
  now.setHours(0, 0, 0, 0);
  return now.getTime();
};

const dateRange = ref<[number, number]>([getToday(), getToday()]);

const formatDate = (timestamp: number): string => {
  const date = new Date(timestamp);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

const fetchStats = async () => {
  loading.value = true;
  try {
    const [startTs, endTs] = dateRange.value;
    const response = await getDashboardStats({
      start: formatDate(startTs),
      end: formatDate(endTs),
    });
    dashboardStats.value = response.data;
  } catch (error) {
    console.error("Failed to load dashboard stats:", error);
  } finally {
    loading.value = false;
  }
};

// 监听日期范围变化
watch(dateRange, () => {
  fetchStats();
});

onMounted(() => {
  fetchStats();
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
      <token-stats-card :stats="dashboardStats" :loading="loading">
        <template #header-extra>
          <n-date-picker
            v-model:value="dateRange"
            type="daterange"
            size="small"
            :shortcuts="{
              [t('dashboard.today')]: () => {
                const today = getToday();
                return [today, today];
              },
              [t('dashboard.yesterday')]: () => {
                const yesterday = getToday() - 24 * 60 * 60 * 1000;
                return [yesterday, yesterday];
              },
              [t('dashboard.last7days')]: () => {
                const today = getToday();
                return [today - 6 * 24 * 60 * 60 * 1000, today];
              },
              [t('dashboard.last30days')]: () => {
                const today = getToday();
                return [today - 29 * 24 * 60 * 60 * 1000, today];
              },
            }"
            clearable
            class="date-picker"
          />
        </template>
      </token-stats-card>

      <!-- 模型使用情况 -->
      <model-usage-card :stats="dashboardStats" :loading="loading" />
    </n-space>
  </div>
</template>

<style scoped>
.dashboard-container {
  max-width: 100%;
}

.date-picker {
  width: 240px;
}

.date-picker :deep(.n-input) {
  border-radius: 999px;
}
</style>

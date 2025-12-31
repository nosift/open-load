import type { ChartData, DashboardStatsResponse, Group } from "@/types/models";
import http from "@/utils/http";

const getClientTimeZone = () => {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone;
  } catch {
    return undefined;
  }
};

/**
 * 获取仪表盘基础统计数据
 */
export const getDashboardStats = (params?: {
  date?: string;
  start?: string;
  end?: string;
  tz?: string;
}) => {
  const tz = params?.tz ?? getClientTimeZone();
  return http.get<DashboardStatsResponse>("/dashboard/stats", {
    params: {
      ...params,
      ...(tz ? { tz } : {}),
    },
  });
};

/**
 * 获取仪表盘图表数据
 * @param groupId 可选的分组ID
 */
export const getDashboardChart = (groupId?: number) => {
  const tz = getClientTimeZone();
  return http.get<ChartData>("/dashboard/chart", {
    params: {
      ...(groupId ? { groupId } : {}),
      ...(tz ? { tz } : {}),
    },
  });
};

/**
 * 获取用于筛选的分组列表
 */
export const getGroupList = () => {
  return http.get<Group[]>("/groups/list");
};

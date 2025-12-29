<script setup lang="ts">
import type { Group } from "@/types/models";
import { getGroupDisplayName } from "@/utils/display";
import { Add, LinkOutline, Search } from "@vicons/ionicons5";
import { NButton, NCard, NEmpty, NInput, NSpin, NTag } from "naive-ui";
import { computed, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import AggregateGroupModal from "./AggregateGroupModal.vue";
import GroupFormModal from "./GroupFormModal.vue";

const { t } = useI18n();

interface Props {
  groups: Group[];
  selectedGroup: Group | null;
  loading?: boolean;
}

interface Emits {
  (e: "group-select", group: Group): void;
  (e: "refresh"): void;
  (e: "refresh-and-select", groupId: number): void;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
});

const emit = defineEmits<Emits>();

const searchText = ref("");
const showGroupModal = ref(false);
// 存储分组项 DOM 元素的引用
const groupItemRefs = ref(new Map());
const showAggregateGroupModal = ref(false);

function getGroupIconText(group: Group): string {
  if (group.group_type === "aggregate") {
    return "AG";
  }

  switch (group.channel_type) {
    case "openai":
      return "AI";
    case "gemini":
      return "GE";
    case "anthropic":
      return "AN";
    default:
      return "--";
  }
}

const filteredGroups = computed(() => {
  if (!searchText.value.trim()) {
    return props.groups;
  }
  const search = searchText.value.toLowerCase().trim();
  return props.groups.filter(
    group =>
      group.name.toLowerCase().includes(search) ||
      group.display_name?.toLowerCase().includes(search)
  );
});

// 监听选中项 ID 的变化，并自动滚动到该项
watch(
  () => props.selectedGroup?.id,
  id => {
    if (!id || props.groups.length === 0) {
      return;
    }

    const element = groupItemRefs.value.get(id);
    if (element) {
      element.scrollIntoView({
        behavior: "smooth", // 平滑滚动
        block: "nearest", // 将元素滚动到最近的边缘
      });
    }
  },
  {
    flush: "post", // 确保在 DOM 更新后执行回调
    immediate: true, // 立即执行一次以处理初始加载
  }
);

function handleGroupClick(group: Group) {
  emit("group-select", group);
}

// 获取渠道类型的标签颜色
function getChannelTagType(channelType: string) {
  switch (channelType) {
    case "openai":
      return "success";
    case "gemini":
      return "info";
    case "anthropic":
      return "warning";
    default:
      return "default";
  }
}

function openCreateGroupModal() {
  showGroupModal.value = true;
}

function openCreateAggregateGroupModal() {
  showAggregateGroupModal.value = true;
}

function handleGroupCreated(group: Group) {
  showGroupModal.value = false;
  showAggregateGroupModal.value = false;
  if (group?.id) {
    emit("refresh-and-select", group.id);
  }
}
</script>

<template>
  <div class="group-list-container">
    <n-card class="group-list-card modern-card" :bordered="false" size="small">
      <!-- 搜索框 -->
      <div class="search-section">
        <n-input
          v-model:value="searchText"
          :placeholder="t('keys.searchGroupPlaceholder')"
          size="small"
          clearable
        >
          <template #prefix>
            <n-icon :component="Search" />
          </template>
        </n-input>
      </div>

      <!-- 分组列表 -->
      <div class="groups-section">
        <n-spin :show="loading" size="small">
          <div v-if="filteredGroups.length === 0 && !loading" class="empty-container">
            <n-empty
              size="small"
              :description="searchText ? t('keys.noMatchingGroups') : t('keys.noGroups')"
            />
          </div>
          <div v-else class="groups-list">
            <div
              v-for="group in filteredGroups"
              :key="group.id"
              class="group-item"
              :class="{
                active: selectedGroup?.id === group.id,
                aggregate: group.group_type === 'aggregate',
              }"
              @click="handleGroupClick(group)"
              :ref="
                el => {
                  if (el) groupItemRefs.set(group.id, el);
                }
              "
            >
              <div class="group-icon">
                <span class="group-icon-text">{{ getGroupIconText(group) }}</span>
              </div>
              <div class="group-content">
                <div class="group-name">{{ getGroupDisplayName(group) }}</div>
                <div class="group-meta">
                  <n-tag size="tiny" :type="getChannelTagType(group.channel_type)">
                    {{ group.channel_type }}
                  </n-tag>
                  <n-tag v-if="group.group_type === 'aggregate'" size="tiny" type="warning" round>
                    {{ t("keys.aggregateGroup") }}
                  </n-tag>
                  <span v-if="group.group_type !== 'aggregate'" class="group-id">
                    #{{ group.name }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </n-spin>
      </div>

      <!-- 添加分组按钮 -->
      <div class="add-section">
        <n-button class="action-btn-custom success" size="small" block @click="openCreateGroupModal">
          <template #icon>
            <n-icon :component="Add" />
          </template>
          {{ t("keys.createGroup") }}
        </n-button>
        <n-button class="action-btn-custom info" size="small" block @click="openCreateAggregateGroupModal">
          <template #icon>
            <n-icon :component="LinkOutline" />
          </template>
          {{ t("keys.createAggregateGroup") }}
        </n-button>
      </div>
    </n-card>
    <group-form-modal v-model:show="showGroupModal" @success="handleGroupCreated" />
    <aggregate-group-modal
      v-model:show="showAggregateGroupModal"
      :groups="groups"
      @success="handleGroupCreated"
    />
  </div>
</template>

<style scoped>
:deep(.n-card__content) {
  height: auto;
}

.groups-section::-webkit-scrollbar {
  width: 1px;
  height: 1px;
}

.group-list-container {
  height: auto;
}

.group-list-card {
  height: auto;
  display: flex;
  flex-direction: column;
  background: var(--sidebar-bg);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-sm);
  backdrop-filter: blur(16px);
}

.group-list-card:hover {
  transform: none;
  box-shadow: none;
}

.search-section {
  height: 44px;
  padding: 0 10px;
}

.groups-section {
  flex: 0 0 auto;
  max-height: calc(100vh - 360px);
  overflow: auto;
  padding: 0 10px;
}

.group-icon-text {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.empty-container {
  padding: 20px 0;
}

.groups-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 100%;
}

.group-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid transparent; /* 默认无边框 */
  font-size: 13px;
  color: var(--text-secondary);
  background: transparent;
  box-sizing: border-box;
  position: relative;
}

/* 聚合分组样式 - 简化 */
.group-item.aggregate {
  border-style: none;
  background: transparent;
}

:root.dark .group-item.aggregate {
  background: transparent;
  border-color: transparent;
}

.group-item:hover,
.group-item.aggregate:hover {
  background: var(--hover-bg);
  color: var(--text-primary);
  border-color: transparent;
  box-shadow: var(--shadow-sm);
}

:root.dark .group-item:hover {
  background: var(--hover-bg);
}

.group-item.active,
:root.dark .group-item.active,
:root.dark .group-item.aggregate.active {
  background: var(--card-bg-solid);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  font-weight: 500;
}

.group-icon {
  font-size: 14px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  border-radius: 6px;
  flex-shrink: 0;
  box-sizing: border-box;
  color: var(--text-tertiary); /* 图标颜色变淡 */
}

.group-item.active .group-icon {
  color: var(--text-primary); /* 选中时加深 */
}

.group-content {
  flex: 1;
  min-width: 0;
}

.group-name {
  font-weight: 500;
  font-size: 14px;
  line-height: 1.2;
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
  flex-wrap: wrap;
}

.group-id {
  opacity: 0.6;
  color: var(--text-tertiary);
  font-size: 10px;
}

.group-item.active .group-id {
  opacity: 0.8;
  color: var(--text-secondary);
}

.add-section {
  border-top: 1px solid var(--border-color);
  padding: 12px 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 滚动条样式 */
.groups-list::-webkit-scrollbar {
  width: 4px;
}

.groups-list::-webkit-scrollbar-track {
  background: transparent;
}

.groups-list::-webkit-scrollbar-thumb {
  background: var(--scrollbar-bg);
  border-radius: 2px;
}

.groups-list::-webkit-scrollbar-thumb:hover {
  background: var(--border-color);
}

/* 暗黑模式特殊样式 */
:root.dark .group-item {
  border-color: transparent;
}

:root.dark .group-icon {
  background: transparent;
  border: none;
}

:root.dark .search-section :deep(.n-input) {
  --n-border: 1px solid var(--border-color);
  --n-border-hover: 1px solid var(--border-hover);
  --n-border-focus: 1px solid var(--primary-color);
  background: rgba(255, 255, 255, 0.05);
}

/* 标签样式优化 - 极简 */
:root.dark .group-meta :deep(.n-tag) {
  background: rgba(255, 255, 255, 0.1);
  border: none;
}

:root.dark .group-item.active .group-meta :deep(.n-tag) {
  background: rgba(255, 255, 255, 0.15);
}

/* 自定义操作按钮样式 */
.action-btn-custom {
  font-weight: 500;
  border-radius: var(--border-radius-md);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  background: var(--card-bg-solid);
  transition: all 0.2s ease;
  box-shadow: none;
}

:root.dark .action-btn-custom {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.2);
  color: white;
}

/* 主要操作按钮 (蓝色) */
.action-btn-custom.success,
.action-btn-custom.info {
  background-color: var(--primary-soft-bg) !important;
  color: var(--primary-soft-text) !important;
  border-color: var(--primary-soft-border) !important;
}

.action-btn-custom.success:hover,
.action-btn-custom.info:hover {
  background-color: var(--primary-soft-bg-hover) !important;
  border-color: var(--primary-soft-border-hover) !important;
}

:root.dark .action-btn-custom.success,
:root.dark .action-btn-custom.info {
  background-color: var(--primary-soft-bg) !important;
  border-color: var(--primary-soft-border) !important;
  color: var(--primary-soft-text) !important;
}

:deep(.action-btn-custom .n-button__state-border),
:deep(.action-btn-custom .n-button__border) {
  box-shadow: none !important;
  border: none !important;
}
</style>

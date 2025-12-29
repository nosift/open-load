<script setup lang="ts">
import { keysApi } from "@/api/keys";
import type { APIKey, Group, KeyStatus } from "@/types/models";
import { appState, triggerSyncOperationRefresh } from "@/utils/app-state";
import { copy } from "@/utils/clipboard";
import { getGroupDisplayName, maskKey } from "@/utils/display";
import {
  AddCircleOutline,
  AlertCircleOutline,
  CheckmarkCircle,
  CopyOutline,
  EyeOffOutline,
  EyeOutline,
  Pencil,
  RemoveCircleOutline,
  Search,
  TimeOutline,
} from "@vicons/ionicons5";
import {
  NButton,
  NDropdown,
  NEmpty,
  NIcon,
  NInput,
  NModal,
  NSpin,
  useDialog,
  type MessageReactive,
} from "naive-ui";
import { h, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import KeyCreateDialog from "./KeyCreateDialog.vue";
import KeyDeleteDialog from "./KeyDeleteDialog.vue";

const { t } = useI18n();

interface KeyRow extends APIKey {
  is_visible: boolean;
}

interface Props {
  selectedGroup: Group | null;
}

const props = defineProps<Props>();

const keys = ref<KeyRow[]>([]);
const loading = ref(false);
const searchText = ref("");
const statusFilter = ref<"all" | "active" | "invalid">("active");
const recentMinutes = ref<number | null>(null);
const currentPage = ref(1);
const pageSize = ref(12);
const total = ref(0);
const totalPages = ref(0);
const keyStats = ref<{ total_keys: number; active_keys: number; invalid_keys: number } | null>(null);
const dialog = useDialog();
const confirmInput = ref("");

// 状态过滤选项
const statusOptions = [
  { label: t("keys.valid"), value: "active" },
  { label: t("keys.invalid"), value: "invalid" },
  { label: t("common.all"), value: "all" },
];

// 更多操作下拉菜单选项
const moreOptions = [
  { label: t("keys.exportAllKeys"), key: "copyAll" },
  { label: t("keys.exportValidKeys"), key: "copyValid" },
  { label: t("keys.exportInvalidKeys"), key: "copyInvalid" },
  { type: "divider" },
  { label: t("keys.restoreAllInvalidKeys"), key: "restoreAll" },
  {
    label: t("keys.clearAllInvalidKeys"),
    key: "clearInvalid",
    props: { style: { color: "#d03050" } },
  },
  {
    label: t("keys.clearAllKeys"),
    key: "clearAll",
    props: { style: { color: "red", fontWeight: "bold" } },
  },
  { type: "divider" },
  { label: t("keys.validateAllKeys"), key: "validateAll" },
  { label: t("keys.validateValidKeys"), key: "validateActive" },
  { label: t("keys.validateInvalidKeys"), key: "validateInvalid" },
];

let testingMsg: MessageReactive | null = null;
const isDeling = ref(false);
const isRestoring = ref(false);

const createDialogShow = ref(false);
const deleteDialogShow = ref(false);

// 备注编辑相关
const notesDialogShow = ref(false);
const editingKey = ref<KeyRow | null>(null);
const editingNotes = ref("");
const viewKeyDialogShow = ref(false);
const viewingKey = ref<KeyRow | null>(null);

watch(
  () => props.selectedGroup,
  async newGroup => {
    if (newGroup) {
      const prevPage = currentPage.value;
      const prevStatus = statusFilter.value;
      resetPage();
      const pageChanged = prevPage !== currentPage.value;
      const statusChanged = prevStatus !== statusFilter.value;
      if (!pageChanged && !statusChanged) {
        await loadKeys();
      }
    }
  },
  { immediate: true }
);

watch([currentPage, pageSize], async () => {
  await loadKeys();
});

watch(statusFilter, async () => {
  if (currentPage.value !== 1) {
    currentPage.value = 1;
  } else {
    await loadKeys();
  }
});

// 监听任务完成事件，自动刷新密钥列表
watch(
  () => appState.groupDataRefreshTrigger,
  () => {
    // 检查是否需要刷新当前分组的密钥列表
    if (appState.lastCompletedTask && props.selectedGroup) {
      // 通过分组名称匹配
      const isCurrentGroup = appState.lastCompletedTask.groupName === props.selectedGroup.name;

      const shouldRefresh =
        appState.lastCompletedTask.taskType === "KEY_VALIDATION" ||
        appState.lastCompletedTask.taskType === "KEY_IMPORT" ||
        appState.lastCompletedTask.taskType === "KEY_DELETE";

      if (isCurrentGroup && shouldRefresh) {
        // 刷新当前分组的密钥列表
        loadKeys();
      }
    }
  }
);

// 处理搜索输入的防抖
function handleSearchInput() {
  if (currentPage.value !== 1) {
    currentPage.value = 1;
  } else {
    loadKeys();
  }
}

function toggleRecentAdded() {
  recentMinutes.value = recentMinutes.value === 3 ? null : 3;
  currentPage.value = 1;
  loadKeys();
}

// 处理更多操作菜单
function handleMoreAction(key: string) {
  switch (key) {
    case "copyAll":
      copyAllKeys();
      break;
    case "copyValid":
      copyValidKeys();
      break;
    case "copyInvalid":
      copyInvalidKeys();
      break;
    case "restoreAll":
      restoreAllInvalid();
      break;
    case "validateAll":
      validateKeys("all");
      break;
    case "validateActive":
      validateKeys("active");
      break;
    case "validateInvalid":
      validateKeys("invalid");
      break;
    case "clearInvalid":
      clearAllInvalid();
      break;
    case "clearAll":
      clearAll();
      break;
  }
}

async function loadKeys() {
  if (!props.selectedGroup?.id) {
    return;
  }

  try {
    loading.value = true;
    const result = await keysApi.getGroupKeys({
      group_id: props.selectedGroup.id,
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value === "all" ? undefined : (statusFilter.value as KeyStatus),
      key_value: searchText.value.trim() || undefined,
      recent_minutes: recentMinutes.value ?? undefined,
    });
    keys.value = result.items as KeyRow[];
    total.value = result.pagination.total_items;
    totalPages.value = result.pagination.total_pages;
    await loadKeyStats();
  } finally {
    loading.value = false;
  }
}

async function loadKeyStats() {
  if (!props.selectedGroup?.id) {
    keyStats.value = null;
    return;
  }

  try {
    const stats = await keysApi.getGroupStats(props.selectedGroup.id);
    keyStats.value = stats?.key_stats ?? null;
  } catch (error) {
    console.error("Load key stats failed", error);
    keyStats.value = null;
  }
}

// 处理批量删除成功后的刷新
async function handleBatchDeleteSuccess() {
  await loadKeys();
  // 触发同步操作刷新
  if (props.selectedGroup) {
    triggerSyncOperationRefresh(props.selectedGroup.name, "BATCH_DELETE");
  }
}

async function copyKey(key: KeyRow) {
  const success = await copy(key.key_value);
  if (success) {
    window.$message.success(t("keys.keyCopied"));
  } else {
    window.$message.error(t("keys.copyFailed"));
  }
}

function copyViewingKey() {
  if (!viewingKey.value) {
    return;
  }
  copyKey(viewingKey.value);
}

async function testKey(_key: KeyRow) {
  if (!props.selectedGroup?.id || !_key.key_value || testingMsg) {
    return;
  }

  testingMsg = window.$message.info(t("keys.testingKey"), {
    duration: 0,
  });

  try {
    const response = await keysApi.testKeys(props.selectedGroup.id, _key.key_value);
    const curValid = response.results?.[0] || {};
    if (curValid.is_valid) {
      window.$message.success(
        t("keys.testSuccess", { duration: formatDuration(response.total_duration) })
      );
    } else {
      window.$message.error(curValid.error || t("keys.testFailed"), {
        keepAliveOnHover: true,
        duration: 5000,
        closable: true,
      });
    }
    await loadKeys();
    // 触发同步操作刷新
    triggerSyncOperationRefresh(props.selectedGroup.name, "TEST_SINGLE");
  } catch (_error) {
    console.error("Test failed");
  } finally {
    testingMsg?.destroy();
    testingMsg = null;
  }
}

function formatDuration(ms: number): string {
  if (ms < 0) {
    return "0ms";
  }

  const minutes = Math.floor(ms / 60000);
  const seconds = Math.floor((ms % 60000) / 1000);
  const milliseconds = ms % 1000;

  let result = "";
  if (minutes > 0) {
    result += `${minutes}m`;
  }
  if (seconds > 0) {
    result += `${seconds}s`;
  }
  if (milliseconds > 0 || result === "") {
    result += `${milliseconds}ms`;
  }

  return result;
}

function toggleKeyVisibility(key: KeyRow) {
  viewingKey.value = key;
  viewKeyDialogShow.value = true;
}

// 获取要显示的值（备注优先，否则显示密钥）
function getDisplayValue(key: KeyRow): string {
  if (key.notes && !key.is_visible) {
    return key.notes;
  }
  return key.is_visible ? key.key_value : maskKey(key.key_value);
}

// 编辑密钥备注
function editKeyNotes(key: KeyRow) {
  editingKey.value = key;
  editingNotes.value = key.notes || "";
  notesDialogShow.value = true;
}

// 保存备注
async function saveKeyNotes() {
  if (!editingKey.value) {
    return;
  }

  try {
    const trimmed = editingNotes.value.trim();
    await keysApi.updateKeyNotes(editingKey.value.id, trimmed);
    editingKey.value.notes = trimmed;
    window.$message.success(t("keys.notesUpdated"));
    notesDialogShow.value = false;
  } catch (error) {
    console.error("Update notes failed", error);
  }
}

async function restoreKey(key: KeyRow) {
  if (!props.selectedGroup?.id || !key.key_value || isRestoring.value) {
    return;
  }

  const d = dialog.warning({
    title: t("keys.restoreKey"),
    content: t("keys.confirmRestoreKey", { key: maskKey(key.key_value) }),
    positiveText: t("common.confirm"),
    negativeText: t("common.cancel"),
    onPositiveClick: async () => {
      if (!props.selectedGroup?.id) {
        return;
      }

      isRestoring.value = true;
      d.loading = true;

      try {
        await keysApi.restoreKeys(props.selectedGroup.id, key.key_value);
        await loadKeys();
        // 触发同步操作刷新
        triggerSyncOperationRefresh(props.selectedGroup.name, "RESTORE_SINGLE");
      } catch (_error) {
        console.error("Restore failed");
      } finally {
        d.loading = false;
        isRestoring.value = false;
      }
    },
  });
}

async function deleteKey(key: KeyRow) {
  if (!props.selectedGroup?.id || !key.key_value || isDeling.value) {
    return;
  }

  const d = dialog.warning({
    title: t("keys.deleteKey"),
    content: t("keys.confirmDeleteKey", { key: maskKey(key.key_value) }),
    positiveText: t("common.confirm"),
    negativeText: t("common.cancel"),
    onPositiveClick: async () => {
      if (!props.selectedGroup?.id) {
        return;
      }

      d.loading = true;
      isDeling.value = true;

      try {
        await keysApi.deleteKeys(props.selectedGroup.id, key.key_value);
        await loadKeys();
        // 触发同步操作刷新
        triggerSyncOperationRefresh(props.selectedGroup.name, "DELETE_SINGLE");
      } catch (_error) {
        console.error("Delete failed");
      } finally {
        d.loading = false;
        isDeling.value = false;
      }
    },
  });
}

function formatRelativeTime(date: string) {
  if (!date) {
    return t("keys.never");
  }
  const now = new Date();
  const target = new Date(date);
  const diffSeconds = Math.floor((now.getTime() - target.getTime()) / 1000);
  const diffMinutes = Math.floor(diffSeconds / 60);
  const diffHours = Math.floor(diffMinutes / 60);
  const diffDays = Math.floor(diffHours / 24);

  if (diffDays > 0) {
    return t("keys.daysAgo", { days: diffDays });
  }
  if (diffHours > 0) {
    return t("keys.hoursAgo", { hours: diffHours });
  }
  if (diffMinutes > 0) {
    return t("keys.minutesAgo", { minutes: diffMinutes });
  }
  if (diffSeconds > 0) {
    return t("keys.secondsAgo", { seconds: diffSeconds });
  }
  return t("keys.justNow");
}

function formatAbsoluteTime(date?: string) {
  if (!date) {
    return "-";
  }
  const parsed = new Date(date);
  if (Number.isNaN(parsed.getTime())) {
    return "-";
  }
  return parsed.toLocaleString(undefined, { hour12: false }).replace(/\//g, "-");
}

function getStatusClass(status: KeyStatus): string {
  switch (status) {
    case "active":
      return "status-valid";
    case "invalid":
      return "status-invalid";
    default:
      return "status-unknown";
  }
}

async function copyAllKeys() {
  if (!props.selectedGroup?.id) {
    return;
  }

  keysApi.exportKeys(props.selectedGroup.id, "all");
}

async function copyValidKeys() {
  if (!props.selectedGroup?.id) {
    return;
  }

  keysApi.exportKeys(props.selectedGroup.id, "active");
}

async function copyInvalidKeys() {
  if (!props.selectedGroup?.id) {
    return;
  }

  keysApi.exportKeys(props.selectedGroup.id, "invalid");
}

async function restoreAllInvalid() {
  if (!props.selectedGroup?.id || isRestoring.value) {
    return;
  }

  const d = dialog.warning({
    title: t("keys.restoreKeys"),
    content: t("keys.confirmRestoreAllInvalid"),
    positiveText: t("common.confirm"),
    negativeText: t("common.cancel"),
    onPositiveClick: async () => {
      if (!props.selectedGroup?.id) {
        return;
      }

      isRestoring.value = true;
      d.loading = true;
      try {
        await keysApi.restoreAllInvalidKeys(props.selectedGroup.id);
        await loadKeys();
        // 触发同步操作刷新
        triggerSyncOperationRefresh(props.selectedGroup.name, "RESTORE_ALL_INVALID");
      } catch (_error) {
        console.error("Restore failed");
      } finally {
        d.loading = false;
        isRestoring.value = false;
      }
    },
  });
}

async function validateKeys(status: "all" | "active" | "invalid") {
  if (!props.selectedGroup?.id || testingMsg) {
    return;
  }

  let statusText = t("common.all");
  if (status === "active") {
    statusText = t("keys.valid");
  } else if (status === "invalid") {
    statusText = t("keys.invalid");
  }

  testingMsg = window.$message.info(t("keys.validatingKeysMsg", { type: statusText }), {
    duration: 0,
  });

  try {
    await keysApi.validateGroupKeys(props.selectedGroup.id, status === "all" ? undefined : status);
    localStorage.removeItem("last_closed_task");
    appState.taskPollingTrigger++;
  } catch (_error) {
    console.error("Test failed");
  } finally {
    testingMsg?.destroy();
    testingMsg = null;
  }
}

async function clearAllInvalid() {
  if (!props.selectedGroup?.id || isDeling.value) {
    return;
  }

  const d = dialog.warning({
    title: t("keys.clearKeys"),
    content: t("keys.confirmClearInvalidKeys"),
    positiveText: t("common.confirm"),
    negativeText: t("common.cancel"),
    onPositiveClick: async () => {
      if (!props.selectedGroup?.id) {
        return;
      }

      isDeling.value = true;
      d.loading = true;
      try {
        const { data } = await keysApi.clearAllInvalidKeys(props.selectedGroup.id);
        window.$message.success(data?.message || t("keys.clearSuccess"));
        await loadKeys();
        // 触发同步操作刷新
        triggerSyncOperationRefresh(props.selectedGroup.name, "CLEAR_ALL_INVALID");
      } catch (_error) {
        console.error("Delete failed");
      } finally {
        d.loading = false;
        isDeling.value = false;
      }
    },
  });
}

async function clearAll() {
  if (!props.selectedGroup?.id || isDeling.value) {
    return;
  }

  dialog.warning({
    title: t("keys.clearAllKeys"),
    content: t("keys.confirmClearAllKeys"),
    positiveText: t("common.confirm"),
    negativeText: t("common.cancel"),
    onPositiveClick: () => {
      confirmInput.value = ""; // Reset before opening second dialog
      dialog.create({
        title: t("keys.enterGroupNameToConfirm"),
        content: () =>
          h("div", null, [
            h("p", null, [
              t("keys.dangerousOperationWarning1"),
              h("strong", null, t("common.all")),
              t("keys.dangerousOperationWarning2"),
              h("strong", { style: { color: "#d03050" } }, props.selectedGroup?.name),
              t("keys.toConfirm"),
            ]),
            h(NInput, {
              value: confirmInput.value,
              "onUpdate:value": v => {
                confirmInput.value = v;
              },
              placeholder: t("keys.enterGroupName"),
            }),
          ]),
        positiveText: t("keys.confirmClear"),
        negativeText: t("common.cancel"),
        onPositiveClick: async () => {
          if (confirmInput.value !== props.selectedGroup?.name) {
            window.$message.error(t("keys.incorrectGroupName"));
            return false; // Prevent dialog from closing
          }

          if (!props.selectedGroup?.id) {
            return;
          }

          isDeling.value = true;
          try {
            await keysApi.clearAllKeys(props.selectedGroup.id);
            window.$message.success(t("keys.clearAllKeysSuccess"));
            await loadKeys();
            // Trigger sync operation refresh
            triggerSyncOperationRefresh(props.selectedGroup.name, "CLEAR_ALL");
          } catch (_error) {
            console.error("Clear all failed", _error);
          } finally {
            isDeling.value = false;
          }
        },
      });
    },
  });
}

function changePage(page: number) {
  currentPage.value = page;
}

function changePageSize(size: number) {
  pageSize.value = size;
  currentPage.value = 1;
}

function resetPage() {
  currentPage.value = 1;
  searchText.value = "";
  statusFilter.value = "active";
}

// 获取各状态的数量
function getStatusCount(status: string): number {
  if (keyStats.value) {
    if (status === "active") {
      return keyStats.value.active_keys;
    }
    if (status === "invalid") {
      return keyStats.value.invalid_keys;
    }
    return keyStats.value.total_keys;
  }

  if (status === "all") {
    return total.value;
  }
  return 0;
}
</script>

<template>
  <div class="key-table-container">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <n-button class="action-btn-custom primary" size="small" @click="createDialogShow = true">
          <template #icon>
            <n-icon :component="AddCircleOutline" />
          </template>
          {{ t("keys.addKey") }}
        </n-button>
        <n-button class="action-btn-custom danger" size="small" @click="deleteDialogShow = true">
          <template #icon>
            <n-icon :component="RemoveCircleOutline" />
          </template>
          {{ t("keys.deleteKey") }}
        </n-button>
        <n-dropdown :options="moreOptions" trigger="click" @select="handleMoreAction">
          <n-button class="action-btn-custom" size="small">
            <template #icon>
              <span style="font-size: 14px">⋯</span>
            </template>
          </n-button>
        </n-dropdown>
      </div>
      <div class="toolbar-right">
        <n-input-group class="search-group-custom">
          <n-input
            v-model:value="searchText"
            :placeholder="t('keys.keyExactMatch')"
            size="small"
            class="search-input-custom"
            clearable
            @keyup.enter="handleSearchInput"
          >
            <template #prefix>
              <n-icon :component="Search" />
            </template>
          </n-input>
          <n-button
            size="small"
            class="search-btn-custom"
            :disabled="loading"
            @click="handleSearchInput"
          >
            {{ t("common.search") }}
          </n-button>
        </n-input-group>
        <button
          class="filter-toggle-btn"
          :class="{ active: recentMinutes === 3 }"
          @click="toggleRecentAdded"
        >
          <n-icon :component="TimeOutline" class="toggle-icon" />
          {{ t("keys.recentAdded3Min") }}
        </button>
      </div>
    </div>

    <!-- 标签页式状态筛选 -->
    <div class="status-tabs">
      <button
        v-for="option in statusOptions"
        :key="option.value"
        class="status-tab"
        :class="{ active: statusFilter === option.value }"
        @click="statusFilter = option.value as any"
      >
        {{ option.label }}
        <span class="count">{{ getStatusCount(option.value) }}</span>
      </button>
    </div>

    <!-- 密钥卡片网格 -->
    <div class="keys-grid-container">
      <n-spin :show="loading">
        <div v-if="keys.length === 0 && !loading" class="empty-container">
          <n-empty :description="t('keys.noMatchingKeys')" />
        </div>
        <div v-else class="keys-grid">
          <div
            v-for="key in keys"
            :key="key.id"
            class="key-card"
            :class="getStatusClass(key.status)"
          >
            <!-- 主要信息行：Key + 快速操作 -->
            <div class="key-main">
              <div class="key-section">
                <n-tag v-if="key.status === 'active'" type="success" :bordered="false" round>
                  <template #icon>
                    <n-icon :component="CheckmarkCircle" />
                  </template>
                  {{ t("keys.validShort") }}
                </n-tag>
                <n-tag v-else :bordered="false" round>
                  <template #icon>
                    <n-icon :component="AlertCircleOutline" />
                  </template>
                  {{ t("keys.invalidShort") }}
                </n-tag>
                <n-input class="key-text" :value="getDisplayValue(key)" readonly size="small" />
                <div class="quick-actions">
                  <n-button
                    size="tiny"
                    text
                    @click="editKeyNotes(key)"
                    :title="t('keys.editNotes')"
                  >
                    <template #icon>
                      <n-icon :component="Pencil" />
                    </template>
                  </n-button>
                  <n-button
                    size="tiny"
                    text
                    @click="toggleKeyVisibility(key)"
                    :title="t('keys.showKeys')"
                  >
                    <template #icon>
                      <n-icon :component="key.is_visible ? EyeOffOutline : EyeOutline" />
                    </template>
                  </n-button>
                  <n-button size="tiny" text @click="copyKey(key)" :title="t('common.copy')">
                    <template #icon>
                      <n-icon :component="CopyOutline" />
                    </template>
                  </n-button>
                </div>
              </div>
            </div>

            <!-- 统计信息 + 操作按钮行 -->
            <div class="key-bottom">
              <div class="key-stats">
                <span class="stat-item">
                  {{ t("keys.requestsShort") }}
                  <strong>{{ key.request_count }}</strong>
                </span>
                <span class="stat-item">
                  {{ t("keys.failuresShort") }}
                  <strong>{{ key.failure_count }}</strong>
                </span>
              </div>
              <n-button-group class="key-actions">
                <n-button
                  round
                  tertiary
                  type="info"
                  size="tiny"
                  @click="testKey(key)"
                  :title="t('keys.testKey')"
                >
                  {{ t("keys.testShort") }}
                </n-button>
                <n-button
                  v-if="key.status !== 'active'"
                  tertiary
                  size="tiny"
                  @click="restoreKey(key)"
                  :title="t('keys.restoreKey')"
                  type="warning"
                >
                  {{ t("keys.restoreShort") }}
                </n-button>
                <n-button
                  round
                  tertiary
                  size="tiny"
                  type="error"
                  @click="deleteKey(key)"
                  :title="t('keys.deleteKey')"
                >
                  {{ t("common.deleteShort") }}
                </n-button>
              </n-button-group>
            </div>
            <div class="key-meta">
              <span
                class="meta-item"
                :title="key.last_used_at ? formatAbsoluteTime(key.last_used_at) : ''"
              >
                <span class="meta-label">{{ t("keys.lastUsed") }}</span>
                <span class="meta-value">
                  {{ key.last_used_at ? formatRelativeTime(key.last_used_at) : t("keys.unused") }}
                </span>
              </span>
              <span class="meta-item" :title="formatAbsoluteTime(key.created_at)">
                <span class="meta-label">{{ t("keys.createdAt") }}</span>
                <span class="meta-value">
                  {{ key.created_at ? formatRelativeTime(key.created_at) : "-" }}
                </span>
              </span>
              <span class="meta-item" :title="formatAbsoluteTime(key.updated_at)">
                <span class="meta-label">{{ t("keys.updatedAt") }}</span>
                <span class="meta-value">
                  {{ key.updated_at ? formatRelativeTime(key.updated_at) : "-" }}
                </span>
              </span>
            </div>
          </div>
        </div>
      </n-spin>
    </div>

    <!-- 分页 -->
    <div class="pagination-container">
      <div class="pagination-info">
        <span>{{ t("keys.totalRecords", { total }) }}</span>
        <n-select
          v-model:value="pageSize"
          :options="[
            { label: t('keys.recordsPerPage', { count: 12 }), value: 12 },
            { label: t('keys.recordsPerPage', { count: 24 }), value: 24 },
            { label: t('keys.recordsPerPage', { count: 60 }), value: 60 },
            { label: t('keys.recordsPerPage', { count: 120 }), value: 120 },
          ]"
          size="small"
          style="width: 100px; margin-left: 12px"
          @update:value="changePageSize"
        />
      </div>
      <div class="pagination-controls">
        <n-button size="small" :disabled="currentPage <= 1" @click="changePage(currentPage - 1)">
          {{ t("common.previousPage") }}
        </n-button>
        <span class="page-info">
          {{ t("keys.pageInfo", { current: currentPage, total: totalPages }) }}
        </span>
        <n-button
          size="small"
          :disabled="currentPage >= totalPages"
          @click="changePage(currentPage + 1)"
        >
          {{ t("common.nextPage") }}
        </n-button>
      </div>
    </div>

    <key-create-dialog
      v-if="selectedGroup?.id"
      v-model:show="createDialogShow"
      :group-id="selectedGroup.id"
      :group-name="getGroupDisplayName(selectedGroup!)"
      @success="loadKeys"
    />

    <key-delete-dialog
      v-if="selectedGroup?.id"
      v-model:show="deleteDialogShow"
      :group-id="selectedGroup.id"
      :group-name="getGroupDisplayName(selectedGroup!)"
      @success="handleBatchDeleteSuccess"
    />
  </div>

  <!-- 备注编辑对话框 -->
  <n-modal v-model:show="notesDialogShow" preset="dialog" :title="t('keys.editKeyNotes')">
    <n-input
      v-model:value="editingNotes"
      type="textarea"
      :placeholder="t('keys.enterNotes')"
      :rows="3"
      maxlength="255"
      show-count
    />
    <template #action>
      <n-button @click="notesDialogShow = false">{{ t("common.cancel") }}</n-button>
      <n-button type="primary" @click="saveKeyNotes">{{ t("common.save") }}</n-button>
    </template>
  </n-modal>

  <!-- 密钥查看窗口 -->
  <n-modal v-model:show="viewKeyDialogShow" preset="dialog" :title="t('keys.keyValue')">
    <div class="key-viewer">
      <n-input
        :value="viewingKey?.key_value || ''"
        type="textarea"
        readonly
        :rows="3"
        class="key-viewer-input"
      />
    </div>
    <template #action>
      <n-button @click="viewKeyDialogShow = false">{{ t("common.close") }}</n-button>
      <n-button type="primary" @click="copyViewingKey">{{ t("common.copy") }}</n-button>
    </template>
  </n-modal>
</template>

<style scoped>
.key-table-container {
  background: var(--bg-secondary);
  border-radius: var(--border-radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  height: auto;
  display: flex;
  flex-direction: column;
  backdrop-filter: blur(16px);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: var(--card-bg-solid);
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
  gap: 16px;
}

.toolbar-left {
  display: flex;
  gap: 10px;
  flex-shrink: 0;
}

.toolbar-right {
  display: flex;
  gap: 10px;
  align-items: center;
  flex: 1;
  justify-content: flex-end;
  min-width: 0;
}

/* 操作按钮样式 */
.action-btn-custom {
  font-weight: 500;
  border-radius: var(--border-radius-md);
  border: 1px solid var(--border-color);
  background: var(--card-bg-solid);
  color: var(--text-primary);
  transition: all 0.2s ease;
  padding: 6px 14px;
  font-size: 13px;
  box-shadow: none;
}

.action-btn-custom:hover {
  background: var(--hover-bg);
  border-color: var(--border-hover);
}

.action-btn-custom.primary {
  background: var(--primary-soft-bg);
  color: var(--primary-soft-text);
  border: 1px solid var(--primary-soft-border);
}

.action-btn-custom.primary:hover {
  background: var(--primary-soft-bg-hover);
  border-color: var(--primary-soft-border-hover);
}

.action-btn-custom.danger {
  color: var(--error-color);
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color);
}

.action-btn-custom.danger:hover {
  background: var(--error-bg);
  border-color: var(--error-color);
}

:root.dark .action-btn-custom {
  background: rgba(255, 255, 255, 0.04);
  border-color: var(--border-color);
  color: var(--text-primary);
}

:root.dark .action-btn-custom:hover {
  background: rgba(255, 255, 255, 0.08);
}

:root.dark .action-btn-custom.primary {
  background: var(--primary-soft-bg);
  border-color: var(--primary-soft-border);
  color: var(--primary-soft-text);
}

:deep(.action-btn-custom .n-button__state-border),
:deep(.action-btn-custom .n-button__border) {
  box-shadow: none !important;
  border: none !important;
}

/* 搜索框样式 */
.search-group-custom {
  display: flex;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  overflow: hidden;
  background: var(--card-bg-solid);
}

.search-input-custom {
  border: none !important;
  background: transparent;
}

.search-input-custom :deep(.n-input__border),
.search-input-custom :deep(.n-input__state-border) {
  border: none !important;
}

.search-btn-custom {
  background: transparent;
  border: none;
  border-left: 1px solid var(--border-color);
  color: var(--text-secondary);
  font-weight: 500;
  border-radius: 0;
  box-shadow: none;
}

.search-btn-custom:hover {
  background: var(--hover-bg);
  color: var(--text-primary);
}

.search-btn-custom:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px var(--primary-color-suppl);
}

:deep(.search-btn-custom .n-button__state-border),
:deep(.search-btn-custom .n-button__border) {
  box-shadow: none !important;
  border: none !important;
}

/* 筛选切换按钮 */
.filter-toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  font-size: 13px;
  color: var(--text-secondary);
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 500;
}

.filter-toggle-btn:hover {
  background: var(--hover-bg);
  border-color: var(--border-hover);
}

.filter-toggle-btn.active {
  background: var(--primary-color-suppl);
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.toggle-icon {
  font-size: 14px;
}

/* 标签页式状态筛选 */
.status-tabs {
  display: flex;
  gap: 0;
  background: transparent;
  border-bottom: 1px solid var(--border-color);
  padding: 0 20px;
}

.status-tab {
  position: relative;
  padding: 12px 20px;
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.status-tab:hover {
  color: var(--text-primary);
  background: var(--hover-bg);
}

.status-tab.active {
  color: var(--text-primary);
  border-bottom-color: var(--primary-color);
  background: transparent;
}

.status-tab .count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: 600;
  border-radius: 10px;
}

.status-tab:not(.active) .count {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

.status-tab.active .count {
  background: var(--primary-color-suppl);
  color: var(--primary-color);
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.more-actions {
  position: relative;
}

.more-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  box-shadow: var(--shadow-lg);
  min-width: 180px;
  z-index: 1000;
  overflow: hidden;
}

.menu-item {
  display: block;
  width: 100%;
  padding: 8px 12px;
  border: none;
  background: none;
  text-align: left;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
  transition: background-color 0.2s;
}

.menu-item:hover {
  background: var(--hover-bg);
}

.menu-item.danger {
  color: var(--error-color);
}

.menu-item.danger:hover {
  background: var(--error-bg);
}

.menu-divider {
  height: 1px;
  background: var(--border-color);
  margin: 4px 0;
}

.btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--primary-color);
  color: #ffffff;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-color-hover);
}

.btn-secondary {
  background: var(--text-tertiary);
  color: #ffffff;
}

.btn-secondary:hover:not(:disabled) {
  background: var(--text-secondary);
}

.more-icon {
  font-size: 16px;
  font-weight: bold;
}

.filter-select,
.search-input,
.page-size-select {
  padding: 4px 8px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  font-size: 12px;
  color: var(--text-primary);
  background: var(--card-bg-solid);
}

.search-input {
  width: 180px;
}

.filter-select:focus,
.search-input:focus,
.page-size-select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-color-suppl);
}


/* 密钥卡片网格 */
.keys-grid-container {
  flex: 0 0 auto;
  overflow: visible;
  padding: 24px;
  background: transparent;
}

.empty-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 65px;
  padding: 16px 0;
}

.keys-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.key-card {
  background: var(--card-bg-solid);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-lg);
  padding: 16px;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-shadow: var(--shadow-sm);
}

.key-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
  border-color: var(--border-hover);
}

/* 状态相关样式 */
.key-card.status-valid {
  border-left: 3px solid var(--success-color);
}

.key-card.status-invalid {
  border-left: 3px solid var(--error-color);
  opacity: 0.95;
}

/* 主要信息行 */
.key-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.key-section {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

/* 底部统计和按钮行 */
.key-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.key-stats {
  display: flex;
  gap: 8px;
  font-size: 12px;
  overflow: hidden;
  color: var(--text-secondary);
  flex: 1;
  min-width: 0;
}

.key-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 12px;
  font-size: 11px;
  color: var(--text-secondary);
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.meta-label {
  font-weight: 600;
  color: var(--text-secondary);
}

.meta-value {
  color: var(--text-primary);
}

.stat-item {
  white-space: nowrap;
  color: var(--text-secondary);
}

.stat-item strong {
  color: var(--text-primary);
  font-weight: 600;
}

.key-actions {
  flex-shrink: 0;
}

.key-actions :deep(.n-button) {
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  padding: 4px 10px;
}

.key-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-secondary);
  flex: 1;
  min-width: 0;
}

.stat-item {
  white-space: nowrap;
  color: var(--text-secondary);
}

.stat-item strong {
  color: var(--text-primary);
  font-weight: 600;
  margin-left: 4px;
}

/* 底部统计和按钮行 */
.key-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  padding-top: 4px;
  border-top: 1px solid var(--border-color);
}

.key-actions {
  flex-shrink: 0;
}

.key-actions :deep(.n-button) {
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  padding: 4px 10px;
}

.key-stats {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--text-secondary);
  flex: 1;
  min-width: 0;
}

.stat-item {
  white-space: nowrap;
  color: var(--text-secondary);
}

.stat-item strong {
  color: var(--text-primary);
  font-weight: 600;
  margin-left: 4px;
}

/* 主要信息行 */
.key-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.key-section {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.key-text {
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
  font-weight: 400;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  font-size: 13px;
  color: var(--text-primary);
  background: var(--code-bg);
}

:deep(.n-input__input-el) {
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
  font-size: 13px;
  color: var(--text-primary);
}

:deep(.n-tag) {
  border-radius: 6px;
  font-weight: 500;
  font-size: 12px;
}

.quick-actions {
  display: flex;
  gap: 2px;
  flex-shrink: 0;
}

.quick-actions :deep(.n-button) {
  border-radius: 6px;
  transition: all 0.2s ease;
}

.quick-actions :deep(.n-button:hover) {
  background: var(--hover-bg);
}

/* 密钥查看弹窗 */
.key-viewer {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.key-viewer-input :deep(.n-input__input-el) {
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, Courier, monospace;
  font-size: 13px;
  color: var(--text-primary);
}

/* 分页 */
.pagination-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: var(--card-bg-solid);
  border-top: 1px solid var(--border-color);
  flex-shrink: 0;
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: var(--text-secondary);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.pagination-controls :deep(.n-button) {
  border-radius: 8px;
  font-weight: 500;
}

.page-info {
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
}

/* 元数据信息 */
.key-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 12px;
  font-size: 11px;
  color: var(--text-secondary);
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.meta-label {
  font-weight: 600;
  color: var(--text-secondary);
}

.meta-value {
  color: var(--text-primary);
}

/* 响应式 */
@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .toolbar-left,
  .toolbar-right {
    width: 100%;
    justify-content: space-between;
  }

  .status-tabs {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  .keys-grid {
    grid-template-columns: 1fr;
  }

  .pagination-container {
    flex-direction: column;
    gap: 12px;
  }
}
</style>

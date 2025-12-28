<script setup lang="ts">
import { type MenuOption } from "naive-ui";
import { computed, h, watch } from "vue";
import { RouterLink, useRoute } from "vue-router";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const props = defineProps({
  mode: {
    type: String,
    default: "vertical", // 默认为垂直
  },
});

const emit = defineEmits(["close"]);

const iconBase = {
  viewBox: "0 0 24 24",
  fill: "none",
  stroke: "currentColor",
  "stroke-width": "2.2",
  "stroke-linecap": "round",
  "stroke-linejoin": "round",
};

const navIcons: Record<string, () => ReturnType<typeof h>> = {
  dashboard: () =>
    h("svg", iconBase, [
      h("path", { d: "M4 20V10" }),
      h("path", { d: "M10 20V4" }),
      h("path", { d: "M16 20V14" }),
      h("path", { d: "M2 20H22" }),
    ]),
  models: () =>
    h("svg", iconBase, [
      h("circle", { cx: "12", cy: "12", r: "8" }),
      h("path", { d: "M12 12V4a8 8 0 0 1 8 8h-8z" }),
    ]),
  keys: () =>
    h("svg", iconBase, [
      h("rect", { x: "6", y: "11", width: "12", height: "9", rx: "2" }),
      h("path", { d: "M8 11V7a4 4 0 0 1 8 0v4" }),
    ]),
  logs: () =>
    h("svg", iconBase, [
      h("path", { d: "M9 3h6a2 2 0 0 1 2 2v2H7V5a2 2 0 0 1 2-2z" }),
      h("path", { d: "M7 7H5a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-2" }),
      h("path", { d: "M9 12h6" }),
      h("path", { d: "M9 16h6" }),
    ]),
  settings: () =>
    h("svg", iconBase, [
      h("path", { d: "M4 6H20" }),
      h("path", { d: "M4 12H20" }),
      h("path", { d: "M4 18H20" }),
      h("circle", { cx: "8", cy: "6", r: "2" }),
      h("circle", { cx: "16", cy: "12", r: "2" }),
      h("circle", { cx: "12", cy: "18", r: "2" }),
    ]),
};

const mainMenuOptions = computed<MenuOption[]>(() => {
  const options: MenuOption[] = [
    renderMenuItem("dashboard", t("nav.dashboard"), "dashboard"),
    renderMenuItem("keys", t("nav.keys"), "keys"),
    renderMenuItem("logs", t("nav.logs"), "logs"),
    renderMenuItem("settings", t("nav.settings"), "settings"),
  ];

  return options;
});

const secondaryMenuOptions = computed<MenuOption[]>(() => [
  renderMenuItem("models", t("nav.modelUsage"), "models"),
]);

const route = useRoute();
const activeMenu = computed(() => route.name);

watch(activeMenu, () => {
  if (props.mode === "vertical") {
    emit("close");
  }
});

function renderMenuItem(key: string, label: string, iconKey: string): MenuOption {
  const icon = navIcons[iconKey] ? navIcons[iconKey]() : null;
  return {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: key,
          },
          class: "nav-menu-item",
        },
        {
          default: () => [
            h("span", { class: "nav-item-icon" }, icon ? [icon] : []),
            h("span", { class: "nav-item-text" }, label),
          ],
        }
      ),
    key,
  };
}
</script>

<template>
  <div class="nav-container">
    <n-menu :mode="mode" :options="mainMenuOptions" :value="activeMenu" class="modern-menu" />
    <div class="nav-separator" />
    <n-menu
      :mode="mode"
      :options="secondaryMenuOptions"
      :value="activeMenu"
      class="modern-menu nav-secondary"
    />
  </div>
</template>

<style scoped>
.nav-container {
  width: 100%;
  padding: 8px;
  display: flex;
  flex-direction: column;
}

.nav-separator {
  height: 1px;
  background: var(--border-color);
  margin: 12px 12px 10px;
  opacity: 0.7;
}

:deep(.nav-menu-item) {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  color: var(--text-secondary);
  width: 100%;
  font-weight: 500;
  font-size: 0.95rem;
}

:deep(.n-menu-item) {
  border-radius: 12px;
  margin: 4px 0;
  background: transparent;
}

/* 垂直模式特定样式 */
:deep(.n-menu--vertical .n-menu-item-content) {
  padding: 10px 16px !important;
  min-height: 44px;
  height: auto;
  line-height: 1.2;
}

:deep(.n-menu--vertical .n-menu-item) {
  margin: 4px 0;
  height: auto;
}

:deep(.n-menu-item:hover) {
  background: transparent;
  color: var(--text-primary);
}

:deep(.n-menu-item-content) {
  border-radius: 12px;
  transition: background 0.2s ease, color 0.2s ease;
  background: transparent;
}

:deep(.n-menu-item-content:hover) {
  background: var(--hover-bg);
}

:deep(.n-menu-item--selected) {
  background: transparent;
  color: var(--text-primary);
  font-weight: 600;
}

:deep(.n-menu-item--selected .n-menu-item-content) {
  background: var(--primary-color-suppl);
}

:deep(.n-menu-item::before),
:deep(.n-menu-item::after),
:deep(.n-menu-item-content::before),
:deep(.n-menu-item-content::after) {
  background: transparent !important;
  box-shadow: none !important;
}

:deep(.n-menu-item--selected::before),
:deep(.n-menu-item--selected::after) {
  display: none !important;
}

/* 选中状态的文字颜色 */
:deep(.n-menu-item--selected .nav-menu-item) {
  color: var(--text-primary);
}

:deep(.nav-item-icon) {
  width: 18px;
  height: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  opacity: 1;
  color: var(--text-primary);
}

:deep(.nav-item-icon svg) {
  width: 18px;
  height: 18px;
}

:deep(.n-menu-item--selected .nav-item-icon) {
  opacity: 1;
  color: var(--text-primary);
}

</style>

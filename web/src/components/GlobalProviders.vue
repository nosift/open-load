<script setup lang="ts">
import { appState } from "@/utils/app-state";
import { actualTheme } from "@/utils/theme";
import { getLocale } from "@/locales";
import {
  darkTheme,
  NConfigProvider,
  NDialogProvider,
  NLoadingBarProvider,
  NMessageProvider,
  useLoadingBar,
  useMessage,
  type GlobalTheme,
  type GlobalThemeOverrides,
  zhCN,
  enUS,
  jaJP,
  dateZhCN,
  dateEnUS,
  dateJaJP,
} from "naive-ui";
import { computed, defineComponent, watch } from "vue";

// 自定义主题配置 - 根据主题动态调整
const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const baseOverrides: GlobalThemeOverrides = {
    common: {
      primaryColor: "#09090b", // 极简黑
      primaryColorHover: "#27272a",
      primaryColorPressed: "#000000",
      primaryColorSuppl: "rgba(9, 9, 11, 0.05)",
      borderRadius: "6px", // 小圆角
      borderRadiusSmall: "4px",
      fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
      // 亮色模式下的背景色
      bodyColor: "#ffffff",
      cardColor: "#ffffff",
      modalColor: "#ffffff",
      popoverColor: "#ffffff",
      tableColor: "#ffffff",
      textColorBase: "#09090b",
    },
    Card: {
      paddingMedium: "24px",
      borderRadius: "8px",
      borderColor: "#e4e4e7", // 显式边框
    },
    Button: {
      fontWeight: "500",
      heightMedium: "36px",
      heightLarge: "44px",
      borderRadiusMedium: "6px",
      border: "1px solid #e4e4e7",
    },
    Input: {
      heightMedium: "36px",
      heightLarge: "44px",
      borderRadius: "6px",
      border: "1px solid #e4e4e7",
      borderHover: "1px solid #a1a1aa",
      borderFocus: "1px solid #09090b",
    },
    Menu: {
      itemHeight: "36px",
      borderRadius: "6px",
    },
    LoadingBar: {
      colorLoading: "#09090b",
      colorError: "#ef4444",
      height: "2px",
    },
  };

  // 暗黑模式下的特殊覆盖
  if (actualTheme.value === "dark") {
    return {
      ...baseOverrides,
      common: {
        ...baseOverrides.common,
        primaryColor: "#ffffff",
        primaryColorHover: "#e4e4e7",
        primaryColorPressed: "#fafafa",
        primaryColorSuppl: "rgba(255, 255, 255, 0.1)",
        
        bodyColor: "#09090b",
        cardColor: "#09090b",
        modalColor: "#09090b",
        popoverColor: "#09090b",
        tableColor: "#09090b",
        inputColor: "#18181b",
        
        textColorBase: "#fafafa",
        textColor1: "#fafafa",
        textColor2: "#a1a1aa",
        textColor3: "#52525b",
        
        borderColor: "#27272a",
        dividerColor: "#27272a",
      },
      Card: {
        ...baseOverrides.Card,
        color: "#09090b",
        borderColor: "#27272a",
        textColor: "#fafafa",
      },
      Input: {
        ...baseOverrides.Input,
        color: "#09090b",
        textColor: "#fafafa",
        border: "1px solid #27272a",
        borderHover: "1px solid #52525b",
        borderFocus: "1px solid #fafafa",
        colorFocus: "#09090b",
      },
      Button: {
        ...baseOverrides.Button,
        border: "1px solid #27272a",
      },
      DataTable: {
        tdColor: "#09090b",
        thColor: "#09090b", // 表头与背景同色，或略有区分
        thTextColor: "#fafafa",
        tdTextColor: "#fafafa",
        borderColor: "#27272a",
      },
      Message: {
        color: "#18181b",
        textColor: "#fafafa",
        iconColor: "#fafafa",
        borderRadius: "6px",
      },
      LoadingBar: {
        colorLoading: "#ffffff",
        colorError: "#ef4444",
        height: "2px",
      },
    };
  }

  return baseOverrides;
});

// 根据当前主题动态返回主题对象
const theme = computed<GlobalTheme | undefined>(() => {
  return actualTheme.value === "dark" ? darkTheme : undefined;
});

// 根据当前语言返回对应的 locale 配置
const locale = computed(() => {
  const currentLocale = getLocale();
  switch (currentLocale) {
    case "zh-CN":
      return zhCN;
    case "en-US":
      return enUS;
    case "ja-JP":
      return jaJP;
    default:
      return zhCN;
  }
});

// 根据当前语言返回对应的日期 locale 配置
const dateLocale = computed(() => {
  const currentLocale = getLocale();
  switch (currentLocale) {
    case "zh-CN":
      return dateZhCN;
    case "en-US":
      return dateEnUS;
    case "ja-JP":
      return dateJaJP;
    default:
      return dateZhCN;
  }
});

function useGlobalMessage() {
  window.$message = useMessage();
}

const LoadingBar = defineComponent({
  setup() {
    const loadingBar = useLoadingBar();
    watch(
      () => appState.loading,
      loading => {
        if (loading) {
          loadingBar.start();
        } else {
          loadingBar.finish();
        }
      }
    );
    return () => null;
  },
});

const Message = defineComponent({
  setup() {
    useGlobalMessage();
    return () => null;
  },
});
</script>

<template>
  <n-config-provider
    :theme="theme"
    :theme-overrides="themeOverrides"
    :locale="locale"
    :date-locale="dateLocale"
  >
    <n-loading-bar-provider>
      <n-message-provider placement="top-right">
        <n-dialog-provider>
          <slot />
          <loading-bar />
          <message />
        </n-dialog-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

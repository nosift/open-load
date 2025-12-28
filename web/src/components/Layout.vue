<script setup lang="ts">
import AppFooter from "@/components/AppFooter.vue";
import GlobalTaskProgressBar from "@/components/GlobalTaskProgressBar.vue";
import LanguageSelector from "@/components/LanguageSelector.vue";
import Logout from "@/components/Logout.vue";
import NavBar from "@/components/NavBar.vue";
import ThemeToggle from "@/components/ThemeToggle.vue";
import { useMediaQuery } from "@vueuse/core";
import { computed, ref, watch } from "vue";
import { useRoute } from "vue-router";

const isMenuOpen = ref(false);
const isMobile = useMediaQuery("(max-width: 768px)");
const route = useRoute();
const showFooter = computed(() => route.name !== "keys");

watch(isMobile, value => {
  if (!value) {
    isMenuOpen.value = false;
  }
});

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value;
};
</script>

<template>
  <div class="layout-container">
    <!-- Header -->
    <header class="app-header">
      <div class="header-left">
        <div class="brand">
          <span class="brand-name">
            <span class="brand-prefix">GPT</span>
            <span class="brand-suffix">Load</span>
          </span>
        </div>
      </div>
      <div class="header-right">
        <language-selector />
        <theme-toggle />
        <logout v-if="!isMobile" />
        <n-button v-if="isMobile" text @click="toggleMenu" class="menu-btn">
          <svg viewBox="0 0 24 24" width="24" height="24">
            <path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z" />
          </svg>
        </n-button>
      </div>
    </header>

    <div class="main-body">
      <!-- Sidebar -->
      <aside v-if="!isMobile" class="app-sidebar">
        <nav-bar mode="vertical" />
      </aside>

      <!-- Mobile Drawer -->
      <n-drawer v-model:show="isMenuOpen" :width="260" placement="left">
        <n-drawer-content
          title="GPT Load"
          body-content-style="padding: 0; display: flex; flex-direction: column; height: 100%;"
        >
          <div style="flex: 1; overflow-y: auto">
            <nav-bar mode="vertical" @close="isMenuOpen = false" />
          </div>
          <div class="mobile-actions">
            <logout />
          </div>
        </n-drawer-content>
      </n-drawer>

      <!-- Main Content -->
      <main class="app-content">
        <div class="content-wrapper">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
        <app-footer v-if="showFooter" />
      </main>
    </div>

    <!-- 全局任务进度条 -->
    <global-task-progress-bar />
  </div>
</template>

<style scoped>
.layout-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: transparent;
}

.app-header {
  height: 64px;
  background-color: var(--header-bg);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  flex-shrink: 0;
  z-index: 1000;
  backdrop-filter: saturate(180%) blur(18px);
  box-shadow: var(--shadow-sm);
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-name {
  font-size: 18px;
  font-weight: 600;
  letter-spacing: -0.01em;
  color: var(--text-primary);
  position: relative;
  padding-bottom: 6px;
  display: inline-flex;
  align-items: baseline;
  gap: 8px;
}

.brand-prefix {
  font-weight: 700;
  letter-spacing: 0.22em;
  text-transform: uppercase;
}

.brand-suffix {
  font-weight: 500;
  letter-spacing: 0.04em;
}

.brand-name::after {
  content: "";
  position: absolute;
  left: 0;
  right: -16px;
  bottom: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--primary-color), transparent);
  border-radius: 999px;
  opacity: 0.85;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.menu-btn {
  border-radius: var(--border-radius-md);
  color: var(--text-secondary);
}

.menu-btn:hover {
  color: var(--text-primary);
  background: var(--hover-bg);
}

.main-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.app-sidebar {
  width: 240px;
  background-color: var(--sidebar-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow-y: auto;
  padding: 12px 8px;
  backdrop-filter: saturate(160%) blur(16px);
}

.app-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  background-color: transparent;
}

.content-wrapper {
  flex: 1;
  padding: 32px 40px;
  max-width: 1320px;
  margin: 0 auto;
  width: 100%;
}

.mobile-actions {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .app-header {
    padding: 0 16px;
  }
  
  .content-wrapper {
    padding: 20px;
  }
}
</style>

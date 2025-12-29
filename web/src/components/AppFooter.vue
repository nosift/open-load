<script setup lang="ts">
import { versionService, type VersionInfo } from "@/services/version";
import {
  BugOutline,
  CheckmarkCircleOutline,
  DocumentTextOutline,
  LogoGithub,
  PeopleOutline,
  TimeOutline,
  WarningOutline,
} from "@vicons/ionicons5";
import { NIcon, NTooltip } from "naive-ui";
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const versionInfo = ref<VersionInfo>({
  currentVersion: "0.1.0",
  latestVersion: null,
  isLatest: false,
  hasUpdate: false,
  releaseUrl: null,
  lastCheckTime: 0,
  status: "checking",
});

const isChecking = ref(false);

// 版本状态配置
const statusConfig = {
  checking: {
    color: "#0066cc",
    icon: TimeOutline,
    text: t("footer.checking"),
  },
  latest: {
    color: "#18a058",
    icon: CheckmarkCircleOutline,
    text: t("footer.latestVersion"),
  },
  "update-available": {
    color: "#f0a020",
    icon: WarningOutline,
    text: t("footer.updateAvailable"),
  },
  error: {
    color: "#d03050",
    icon: WarningOutline,
    text: t("footer.checkFailed"),
  },
};

const formatVersion = (version: string): string => {
  return version.startsWith("v") ? version : `v${version}`;
};

const checkVersion = async () => {
  if (isChecking.value) {
    return;
  }

  isChecking.value = true;
  try {
    const result = await versionService.checkForUpdates();
    versionInfo.value = result;
  } catch (error) {
    console.warn("Version check failed:", error);
  } finally {
    isChecking.value = false;
  }
};

const FALLBACK_RELEASE_URL = "https://github.com/nosift/gpt-load/releases";

const handleVersionClick = () => {
  if (isChecking.value) {
    return;
  }

  const targetUrl = versionInfo.value.releaseUrl || FALLBACK_RELEASE_URL;
  window.open(targetUrl, "_blank", "noopener,noreferrer");
};

onMounted(() => {
  checkVersion();
});
</script>

<template>
  <footer class="app-footer">
    <div class="footer-container">
      <!-- 主要信息区 -->
      <div class="footer-main">
        <span class="project-info">
          <a href="https://github.com/nosift/gpt-load" target="_blank" rel="noopener noreferrer">
            <span class="footer-brand">
              <span class="brand-prefix">Open</span>
              <span class="brand-suffix">Load</span>
            </span>
          </a>
        </span>

        <n-divider vertical />

        <!-- 版本信息 -->
        <div
          class="version-container"
          :class="{
            'version-clickable': !isChecking,
            'version-checking': isChecking,
          }"
          @click="handleVersionClick"
        >
          <n-icon
            v-if="statusConfig[versionInfo.status].icon"
            :component="statusConfig[versionInfo.status].icon"
            :color="statusConfig[versionInfo.status].color"
            :size="14"
            class="version-icon"
          />
          <span class="version-text">
            {{ formatVersion(versionInfo.currentVersion) }}
            -
            <span :style="{ color: statusConfig[versionInfo.status].color }">
              {{ statusConfig[versionInfo.status].text }}
              <template v-if="versionInfo.status === 'update-available'">
                [{{ formatVersion(versionInfo.latestVersion || "") }}]
              </template>
            </span>
          </span>
        </div>

        <n-divider vertical />

        <!-- 链接区 -->
        <div class="links-container">
          <n-tooltip trigger="hover" placement="top">
            <template #trigger>
              <a
                href="https://github.com/nosift/gpt-load#readme"
                target="_blank"
                rel="noopener noreferrer"
                class="footer-link"
              >
                <n-icon :component="DocumentTextOutline" :size="14" class="link-icon" />
                <span>{{ t("footer.docs") }}</span>
              </a>
            </template>
            {{ t("footer.officialDocs") }}
          </n-tooltip>

          <n-tooltip trigger="hover" placement="top">
            <template #trigger>
              <a
                href="https://github.com/nosift/gpt-load"
                target="_blank"
                rel="noopener noreferrer"
                class="footer-link"
              >
                <n-icon :component="LogoGithub" :size="14" class="link-icon" />
                <span>GitHub</span>
              </a>
            </template>
            {{ t("footer.viewSource") }}
          </n-tooltip>

          <n-tooltip trigger="hover" placement="top">
            <template #trigger>
              <a
                href="https://github.com/nosift/gpt-load/issues"
                target="_blank"
                rel="noopener noreferrer"
                class="footer-link"
              >
                <n-icon :component="BugOutline" :size="14" class="link-icon" />
                <span>{{ t("footer.feedback") }}</span>
              </a>
            </template>
            {{ t("footer.reportIssue") }}
          </n-tooltip>

          <n-tooltip trigger="hover" placement="top">
            <template #trigger>
              <a
                href="https://github.com/nosift/gpt-load/graphs/contributors"
                target="_blank"
                rel="noopener noreferrer"
                class="footer-link"
              >
                <n-icon :component="PeopleOutline" :size="14" class="link-icon" />
                <span>{{ t("footer.contributors") }}</span>
              </a>
            </template>
            {{ t("footer.viewContributors") }}
          </n-tooltip>

        </div>

        <n-divider vertical />

        <!-- 版权信息 -->
        <div class="copyright-container">
          <span class="copyright-text">
            (c) 2025 by
            <a
              href="https://github.com/nosift"
              target="_blank"
              rel="noopener noreferrer"
              class="author-link"
            >
              nosift
            </a>
          </span>
          <span class="license-text">MIT License</span>
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.app-footer {
  background: var(--footer-bg);
  backdrop-filter: blur(20px);
  border-top: 1px solid var(--border-color-light);
  padding: 12px 24px;
  font-size: 14px;
  min-height: 52px;
  box-shadow: 0 -8px 24px rgba(0, 0, 0, 0.04);
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-main {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  line-height: 1.4;
}

.project-info {
  color: var(--text-secondary);
  font-weight: 500;
}

.project-info a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 600;
}

.project-info a:hover {
  text-decoration: underline;
}

.footer-brand {
  display: inline-flex;
  align-items: baseline;
  gap: 8px;
  letter-spacing: -0.01em;
  position: relative;
  padding-bottom: 6px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.footer-brand::after {
  content: "";
  position: absolute;
  left: 0;
  right: -16px;
  bottom: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--primary-color), transparent);
  opacity: 0.85;
  border-radius: 999px;
}

.footer-brand .brand-prefix {
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.footer-brand .brand-suffix {
  font-weight: 500;
  letter-spacing: 0.04em;
}

/* 版本信息区域 */
.version-container {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.version-icon {
  display: flex;
  align-items: center;
}

.version-text {
  font-weight: 500;
  font-size: 13px;
  color: var(--text-secondary);
  white-space: nowrap;
}

.version-clickable {
  cursor: pointer;
}

.version-clickable:hover {
  background: var(--warning-bg);
  transform: translateY(-1px);
}

.version-checking {
  opacity: 0.7;
}

/* 链接区域 */
.links-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.footer-link {
  display: flex;
  align-items: center;
  gap: 4px;
  color: var(--text-secondary);
  text-decoration: none;
  padding: 4px 6px;
  border-radius: 4px;
  transition: all 0.2s ease;
  font-size: 13px;
  white-space: nowrap;
}

.footer-link:hover {
  color: var(--primary-color);
  background: var(--primary-color-suppl);
  transform: translateY(-1px);
}

.link-icon {
  display: flex;
  align-items: center;
}

/* 版权信息区域 */
.copyright-container {
  display: flex;
  align-items: center;
  gap: 8px;
}

.copyright-text {
  color: var(--text-tertiary);
  font-size: 12px;
}

.license-text {
  color: var(--text-tertiary);
  font-size: 12px;
}

.author-link {
  font-weight: 600;
  color: var(--primary-color);
  text-decoration: none;
}

.author-link:hover {
  text-decoration: underline !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-footer {
    padding: 10px 16px;
    height: auto;
  }

  .footer-main {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }

  .footer-main :deep(.n-divider) {
    display: none;
  }

  .links-container {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .footer-main {
    gap: 6px;
  }

  .links-container {
    flex-wrap: wrap;
    justify-content: center;
    gap: 12px;
  }

  .project-info {
    font-size: 12px;
  }

  .footer-link {
    font-size: 12px;
  }
}
</style>

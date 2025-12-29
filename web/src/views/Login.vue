<script setup lang="ts">
import { useAuthService } from "@/services/auth";
import { LockClosedOutline } from "@vicons/ionicons5";
import { NButton, NCard, NForm, NFormItem, NIcon, NInput } from "naive-ui";
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import AppFooter from "@/components/AppFooter.vue";
import LanguageSelector from "@/components/LanguageSelector.vue";
import ThemeToggle from "@/components/ThemeToggle.vue";

const { t } = useI18n();
const router = useRouter();
const { login } = useAuthService();

const authKey = ref("");
const loading = ref(false);

async function handleLogin() {
  if (!authKey.value) {
    return;
  }

  loading.value = true;
  try {
    const success = await login(authKey.value);
    if (success) {
      window.$message?.success(t("login.loginSuccess"));
      router.push({ name: "dashboard" });
    } else {
      window.$message?.error(t("login.invalidKey"));
    }
  } catch (error) {
    console.error("Login error:", error);
    window.$message?.error(t("common.error"));
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-header">
      <div class="logo-area">
        <h1 class="brand-name">
          <span class="brand-prefix">Open</span>
          <span class="brand-suffix">Load</span>
        </h1>
      </div>
      <div class="actions">
        <language-selector />
        <theme-toggle />
      </div>
    </div>

    <div class="login-content">
      <div class="login-spotlight" aria-hidden="true"></div>
      <div class="login-card-wrapper">
        <n-card :bordered="true" class="login-card" size="large">
          <div class="card-header">
            <h2>{{ t("login.welcome") }}</h2>
            <p class="subtitle">{{ t("login.welcomeDesc") }}</p>
          </div>

          <n-form @submit.prevent="handleLogin">
            <n-form-item>
              <n-input
                v-model:value="authKey"
                type="password"
                show-password-on="click"
                :placeholder="t('login.authKeyPlaceholder')"
                :loading="loading"
                @keydown.enter.prevent="handleLogin"
                size="large"
              >
                <template #prefix>
                  <n-icon :component="LockClosedOutline" />
                </template>
              </n-input>
            </n-form-item>
            
            <n-button
              type="primary"
              block
              size="large"
              :loading="loading"
              :disabled="!authKey"
              @click="handleLogin"
              class="login-button"
            >
              {{ t("login.loginButton") }}
            </n-button>
          </n-form>
        </n-card>
      </div>
    </div>
    
    <app-footer />
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-gradient);
  overflow-y: auto;
}

.login-header {
  padding: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--header-bg);
  border-bottom: 1px solid var(--border-color);
  backdrop-filter: saturate(180%) blur(18px);
  box-shadow: var(--shadow-sm);
}

.logo-area {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-name {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
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
  letter-spacing: 0.12em;
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

.actions {
  display: flex;
  gap: 12px;
}

.login-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  position: relative;
  overflow: hidden;
}

.login-spotlight {
  position: absolute;
  width: 720px;
  height: 720px;
  border-radius: 50%;
  background:
    radial-gradient(circle at center, rgba(255, 255, 255, 0.8), transparent 65%),
    radial-gradient(circle at 30% 30%, rgba(0, 113, 227, 0.12), transparent 60%);
  top: -240px;
  right: -140px;
  filter: blur(12px);
  opacity: 0.8;
  pointer-events: none;
}

.login-card-wrapper {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
}

.login-card {
  /* 移除 Naive UI 默认的边框，使用我们自定义的更细腻的边框 */
  /* 这里依赖全局变量的覆盖，但也可以强制指定 */
  box-shadow: var(--shadow-md);
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-lg);
  backdrop-filter: blur(16px);
}

.card-header {
  text-align: center;
  margin-bottom: 28px;
}

.card-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 8px;
  color: var(--text-primary);
}

.subtitle {
  color: var(--text-secondary);
  font-size: 0.875rem;
  margin: 0;
}

.login-button {
  font-weight: 500;
}
</style>

import axios from "axios";

export interface GitHubRelease {
  tag_name: string;
  html_url: string;
  published_at: string;
  name: string;
}

export interface GitHubTag {
  name: string;
}

export interface VersionInfo {
  currentVersion: string;
  latestVersion: string | null;
  isLatest: boolean;
  hasUpdate: boolean;
  releaseUrl: string | null;
  lastCheckTime: number;
  status: "checking" | "latest" | "update-available" | "error";
}

const CACHE_KEY = "open-load-version-info";
const CACHE_DURATION = 30 * 60 * 1000;

class VersionService {
  private currentVersion: string;

  constructor() {
    this.currentVersion = import.meta.env.VITE_VERSION || "1.0.0";
  }

  private parseSemver(version: string): { major: number; minor: number; patch: number } | null {
    const normalized = version.trim();
    const match = normalized.match(/^v?(\d+)\.(\d+)\.(\d+)/);
    if (!match) {
      return null;
    }

    const major = Number(match[1]);
    const minor = Number(match[2]);
    const patch = Number(match[3]);

    if (Number.isNaN(major) || Number.isNaN(minor) || Number.isNaN(patch)) {
      return null;
    }

    return { major, minor, patch };
  }

  /**
   * 获取缓存的版本信息
   */
  private getCachedVersionInfo(): VersionInfo | null {
    try {
      const cached = localStorage.getItem(CACHE_KEY);
      if (!cached) {
        return null;
      }

      const versionInfo: VersionInfo = JSON.parse(cached);
      const now = Date.now();

      // 检查缓存是否过期
      if (now - versionInfo.lastCheckTime > CACHE_DURATION) {
        return null;
      }

      // 检查缓存中的版本号是否与当前应用版本号一致
      if (versionInfo.currentVersion !== this.currentVersion) {
        this.clearCache();
        return null;
      }

      return versionInfo;
    } catch (error) {
      console.warn("Failed to parse cached version info:", error);
      localStorage.removeItem(CACHE_KEY);
      return null;
    }
  }

  /**
   * 缓存版本信息
   */
  private setCachedVersionInfo(versionInfo: VersionInfo): void {
    try {
      localStorage.setItem(CACHE_KEY, JSON.stringify(versionInfo));
    } catch (error) {
      console.warn("Failed to cache version info:", error);
    }
  }

  /**
   * 比较版本号 (简单的语义化版本比较)
   */
  private compareVersions(current: string, latest: string): number {
    const currentParsed = this.parseSemver(current);
    const latestParsed = this.parseSemver(latest);

    // 无法解析的版本号（如 main）不参与比较，视为“不需要更新”
    if (!currentParsed || !latestParsed) {
      return 0;
    }

    const currentParts = [currentParsed.major, currentParsed.minor, currentParsed.patch];
    const latestParts = [latestParsed.major, latestParsed.minor, latestParsed.patch];

    for (let i = 0; i < 3; i++) {
      const currentPart = currentParts[i] || 0;
      const latestPart = latestParts[i] || 0;

      if (currentPart < latestPart) {
        return -1;
      }
      if (currentPart > latestPart) {
        return 1;
      }
    }

    return 0;
  }

  /**
   * 从 GitHub API 获取最新版本
   */
  private async fetchLatestRelease(): Promise<GitHubRelease | null> {
    try {
      const response = await axios.get(
        "https://api.github.com/repos/nosift/gpt-load/releases/latest",
        {
          timeout: 10000,
          headers: {
            Accept: "application/vnd.github.v3+json",
          },
        }
      );

      if (response.status === 200 && response.data) {
        return response.data;
      }

      return null;
    } catch (error) {
      console.warn("Failed to fetch latest version from GitHub:", error);
      return null;
    }
  }

  /**
   * 从 GitHub Tags 获取最新 tag（用于未发布 Release 的场景）
   */
  private async fetchLatestTag(): Promise<GitHubTag | null> {
    try {
      const response = await axios.get("https://api.github.com/repos/nosift/gpt-load/tags?per_page=100", {
        timeout: 10000,
        headers: {
          Accept: "application/vnd.github.v3+json",
        },
      });

      if (response.status !== 200 || !Array.isArray(response.data)) {
        return null;
      }

      const tags: GitHubTag[] = response.data;
      const firstSemverTag = tags.find(tag => /^v?\d+\.\d+\.\d+/.test(tag.name));
      return firstSemverTag || tags[0] || null;
    } catch (error) {
      console.warn("Failed to fetch latest tags from GitHub:", error);
      return null;
    }
  }

  /**
   * 检查版本更新
   */
  async checkForUpdates(): Promise<VersionInfo> {
    // 先检查缓存
    const cached = this.getCachedVersionInfo();
    if (cached) {
      return cached;
    }

    // 创建初始状态
    const versionInfo: VersionInfo = {
      currentVersion: this.currentVersion,
      latestVersion: null,
      isLatest: false,
      hasUpdate: false,
      releaseUrl: null,
      lastCheckTime: Date.now(),
      status: "checking",
    };

    try {
      const release = await this.fetchLatestRelease();
      if (release) {
        const comparison = this.compareVersions(this.currentVersion, release.tag_name);

        versionInfo.latestVersion = release.tag_name;
        versionInfo.releaseUrl = release.html_url;
        versionInfo.isLatest = comparison >= 0;
        versionInfo.hasUpdate = comparison < 0;
        versionInfo.status = comparison < 0 ? "update-available" : "latest";

        this.setCachedVersionInfo(versionInfo);
        return versionInfo;
      }

      const tag = await this.fetchLatestTag();
      if (tag?.name) {
        const comparison = this.compareVersions(this.currentVersion, tag.name);

        versionInfo.latestVersion = tag.name;
        versionInfo.releaseUrl = `https://github.com/nosift/gpt-load/releases/tag/${encodeURIComponent(tag.name)}`;
        versionInfo.isLatest = comparison >= 0;
        versionInfo.hasUpdate = comparison < 0;
        versionInfo.status = comparison < 0 ? "update-available" : "latest";

        this.setCachedVersionInfo(versionInfo);
        return versionInfo;
      }

      versionInfo.status = "error";
    } catch (error) {
      console.warn("Version check failed:", error);
      versionInfo.status = "error";
    }

    return versionInfo;
  }

  /**
   * 获取当前版本号
   */
  getCurrentVersion(): string {
    return this.currentVersion;
  }

  /**
   * 清除缓存
   */
  clearCache(): void {
    localStorage.removeItem(CACHE_KEY);
  }
}

export const versionService = new VersionService();

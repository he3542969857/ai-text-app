import { defineConfig, devices } from '@playwright/test'

// E2E 针对线上部署运行(可用 BASE_URL 覆盖)。
export default defineConfig({
  testDir: './tests',
  timeout: 60_000,
  expect: { timeout: 15_000 },
  use: {
    baseURL: process.env.BASE_URL || 'http://36.213.150.205:3000',
    screenshot: 'only-on-failure',
    trace: 'off',
    // 直连公网,不走系统代理
    ignoreHTTPSErrors: true,
  },
  // 使用系统自带 Edge(channel: msedge),免下载浏览器内核
  projects: [{ name: 'edge', use: { ...devices['Desktop Edge'], channel: 'msedge' } }],
})

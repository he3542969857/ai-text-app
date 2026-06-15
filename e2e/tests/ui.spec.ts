import { test, expect } from '@playwright/test'

const shotDir = 'shots'

test('① 列表页展示三个功能', async ({ page }) => {
  await page.goto('/')
  await expect(page.locator('.feature-name', { hasText: '中译英' })).toBeVisible()
  await expect(page.locator('.feature-name', { hasText: '英译中' })).toBeVisible()
  await expect(page.locator('.feature-name', { hasText: '文本总结' })).toBeVisible()
  await expect(page.getByText('更聪明的处理方式')).toBeVisible()
  await page.screenshot({ path: `${shotDir}/01-home.png`, fullPage: true })
})

test('② 翻译:流式输出 + 状态完成', async ({ page }) => {
  await page.goto('/translate')
  await page.locator('textarea').first().fill('Hello, this is an end-to-end test.')
  // 提交按钮含 ✨,与导航的"翻译"区分
  await page.getByRole('button', { name: /翻译\s*✨/ }).click()
  const out = page.locator('.stream-output')
  await expect(out).not.toHaveText('等待输出…', { timeout: 30_000 })
  await expect(page.getByText(/完成|done/).first()).toBeVisible({ timeout: 30_000 })
  await page.screenshot({ path: `${shotDir}/02-translate.png`, fullPage: true })
  expect((await out.innerText()).trim().length).toBeGreaterThan(0)
})

test('③ 深浅主题切换', async ({ page }) => {
  await page.goto('/')
  const html = page.locator('html')
  const before = await html.getAttribute('data-theme')
  await page.locator('.theme-btn').click()
  await expect(async () => {
    expect(await html.getAttribute('data-theme')).not.toBe(before)
  }).toPass()
  await page.screenshot({ path: `${shotDir}/03-theme.png`, fullPage: true })
})

test('④ 总结:字数控制 + 流式输出', async ({ page }) => {
  await page.goto('/summarize')
  await expect(page.getByText('字数上限')).toBeVisible()
  await page.locator('textarea').first().fill(
    '人工智能正在深刻改变软件开发的方式,工程师借助大模型完成更多创造性工作,效率显著提升,协作模式也在演进。',
  )
  await page.getByRole('button', { name: /总结\s*✨/ }).click()
  const out = page.locator('.stream-output')
  await expect(out).not.toHaveText('等待输出…', { timeout: 30_000 })
  await page.screenshot({ path: `${shotDir}/04-summarize.png`, fullPage: true })
  expect((await out.innerText()).trim().length).toBeGreaterThan(0)
})

test('⑤ 历史页渲染记录', async ({ page }) => {
  await page.goto('/history')
  await expect(page.getByText('调用历史')).toBeVisible()
  await expect(page.locator('.n-data-table-tr').first()).toBeVisible({ timeout: 15_000 })
  await page.screenshot({ path: `${shotDir}/05-history.png`, fullPage: true })
})

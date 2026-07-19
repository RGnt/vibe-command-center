const { chromium } = require('playwright');
(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  await page.goto('http://127.0.0.1:8089');
  
  await page.getByRole('button', { name: "Don't have an account? Sign up" }).click();
  await page.getByPlaceholder('name@example.com').fill(`test_debug_${Date.now()}@example.com`);
  await page.getByPlaceholder('••••••••').fill('password123');
  await page.getByRole('button', { name: 'Sign Up' }).click();
  
  await page.waitForTimeout(2000);
  try {
    await page.getByRole('button', { name: 'Show Error' }).click({ timeout: 2000 });
    await page.waitForTimeout(500);
    const content = await page.evaluate(() => document.body.innerText);
    console.log("ERROR OUTPUT:\n", content);
  } catch (e) {
    console.log('No error button found');
  }
  await browser.close();
})();

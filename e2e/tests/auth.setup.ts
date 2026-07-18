import { test as setup, expect } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';

const authFile = path.join(__dirname, '../playwright/.auth/user.json');

setup('authenticate', async ({ page }) => {
  // Ensure the directory exists
  fs.mkdirSync(path.dirname(authFile), { recursive: true });

  await page.goto('/');

  // Wait for the login page to load
  await expect(page.getByRole('heading', { name: 'Welcome Back' })).toBeVisible();

  // Switch to Sign up
  await page.getByRole('button', { name: "Don't have an account? Sign up" }).click();
  await expect(page.getByRole('heading', { name: 'Create Account' })).toBeVisible();

  // Fill out the registration form
  const uniqueEmail = `test_${Date.now()}@example.com`;
  await page.getByPlaceholder('name@example.com').fill(uniqueEmail);
  await page.getByPlaceholder('••••••••').fill('password123');
  
  // Submit
  await page.getByRole('button', { name: 'Sign Up' }).click();

  // Wait until we are redirected or the main app appears (KanbanX header)
  await expect(page.getByText('KanbanX')).toBeVisible({ timeout: 10000 });

  // Save the authenticated state
  await page.context().storageState({ path: authFile });
});

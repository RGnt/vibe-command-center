import { test, expect } from '@playwright/test';

test('app loads and shows title', async ({ page }) => {
  await page.goto('/');
  // Check for the main header text instead of document title
  await expect(page.getByText('KanbanX')).toBeVisible();
});

test('loads general project by default', async ({ page }) => {
  await page.goto('/');
  
  // Verify Sidebar has "General Project"
  await expect(page.getByText('General Project').first()).toBeVisible();

  // Verify Kanban board headers (seeding creates "To Do")
  await expect(page.getByText('To Do', { exact: true })).toBeVisible();
});

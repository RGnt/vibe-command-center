import { test, expect } from '@playwright/test';

test.describe('Wiki Features', () => {
  test('navigate to general wiki and see default content', async ({ page }) => {
    await page.goto('/');

    // Go to General Wiki
    await page.getByText('General Wiki').click();

    // Verify seeded Wiki view loaded
    await expect(page.getByText('Global Guide').first()).toBeVisible();

    // Create a new wiki page
    await page.getByRole('button', { name: 'New Page' }).click();
    
    // Fill out the form
    await page.locator('input').nth(0).fill('System Architecture');
    await page.locator('input').nth(1).fill('Docs');
    await page.locator('input').nth(2).fill('arch');
    await page.locator('textarea').fill('# Architecture\nThis is a test.');
    
    // Save
    await page.getByRole('button', { name: 'Save' }).click();

    // Verify it was created
    await expect(page.getByText('System Architecture').first()).toBeVisible();
    await expect(page.getByText('This is a test.')).toBeVisible();
  });
});

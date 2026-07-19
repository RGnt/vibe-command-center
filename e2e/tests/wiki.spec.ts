import { test, expect } from '@playwright/test';

test.describe('Wiki Features', () => {
  test('navigate to general wiki and see default content', async ({ page }) => {
    await page.goto('/');

    // Go to General Wiki
    await page.getByText('General Wiki').click();

    // Verify seeded Wiki view loaded
    await expect(page.getByText('Welcome to the Wiki').first()).toBeVisible();

    // Create a new wiki page
    await page.getByRole('button', { name: 'New Page' }).click();
    
    // Fill out the form
    await page.getByPlaceholder('Page Title').fill('System Architecture');
    await page.getByPlaceholder('Category (e.g. Engineering)').fill('Docs');
    await page.getByPlaceholder('URL Slug (e.g. my-page)').fill('arch');
    
    // Type into Monaco editor
    await page.locator('.monaco-editor').first().click();
    // Select all text and delete it before typing
    await page.keyboard.press('Control+A');
    await page.keyboard.press('Delete');
    await page.keyboard.type('# Architecture\nThis is a test.');
    
    // Save
    await page.getByRole('button', { name: 'Save Page' }).click();

    // Verify it was created
    await expect(page.getByText('System Architecture').first()).toBeVisible();
    await expect(page.getByText('This is a test.')).toBeVisible();
  });
});

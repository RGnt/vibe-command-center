import { test, expect } from '@playwright/test';

test.describe('Diagram Features', () => {
  test('open diagram editor and save a basic diagram', async ({ page }) => {
    await page.goto('/');

    // Open Diagram Editor
    await page.getByText('Mermaid Editor').click();
    
    // Verify it opened
    const saveBtn = page.getByRole('button', { name: 'Save', exact: true });
    await expect(saveBtn).toBeVisible();

    // Type in the diagram name
    const titleInput = page.locator('input.w-48').first();
    await titleInput.fill('Test Flow');
    
    // Save
    await saveBtn.click();

    // Wait a moment
    await page.waitForTimeout(1000);
    
    // Verify it appears in load dropdown
    const loadSelect = page.locator('select').first();
    await loadSelect.selectOption({ label: 'Test Flow' });
  });
});

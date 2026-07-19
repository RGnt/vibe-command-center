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

  test('split-pane layout is visible with code editor and canvas', async ({ page }) => {
    await page.goto('/mermaid-editor');

    // Left pane: code editor heading
    await expect(page.getByText('Mermaid Source')).toBeVisible();

    // Apply Code button inside left pane
    const applyBtn = page.getByRole('button', { name: 'Apply Code' });
    await expect(applyBtn).toBeVisible();

    // Monaco editor present
    const editor = page.locator('.monaco-editor').first();
    await expect(editor).toBeVisible();
  });

  test('zoom controls are visible on the canvas', async ({ page }) => {
    await page.goto('/mermaid-editor');

    // Zoom buttons should be in the canvas toolbar
    await expect(page.getByRole('button', { name: 'Zoom In' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Zoom Out' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Reset Zoom' })).toBeVisible();
  });

  test('applying code renders the diagram', async ({ page }) => {
    await page.goto('/mermaid-editor');

    // Apply Code and confirm canvas re-renders
    const applyBtn = page.getByRole('button', { name: 'Apply Code' });
    await applyBtn.click();

    // SVG should be in the canvas after apply
    await page.waitForTimeout(1000);
    const svg = page.locator('.mermaid-container svg');
    await expect(svg).toBeVisible();
  });

  test('drag to connect SVG nodes appends edge to code', async ({ page }) => {
    await page.goto('/mermaid-editor');

    // Apply the default code (already has Alice + Bob actors) to render SVG
    const applyBtn = page.getByRole('button', { name: 'Apply Code' });
    await applyBtn.click();
    await page.waitForTimeout(1500);

    // Wait for actors to appear in the SVG
    const firstNode = page.locator('.actor').first();
    await expect(firstNode).toBeVisible({ timeout: 8000 });
    await firstNode.click();

    // Connection handle should appear after selecting a node
    const handle = page.locator('.connection-handle').first();
    await expect(handle).toBeVisible({ timeout: 5000 });

    // Drag from handle to second actor
    const secondNode = page.locator('.actor').nth(1);
    const targetBox = await secondNode.boundingBox();
    if (targetBox) {
      await handle.hover();
      await page.mouse.down();
      await page.mouse.move(
        targetBox.x + targetBox.width / 2,
        targetBox.y + targetBox.height / 2,
        { steps: 10 }
      );
      await page.mouse.up();

      // Apply the updated code and verify '->>' edge was appended
      await applyBtn.click();
      await page.waitForTimeout(500);
      await expect(page.locator('.monaco-editor').first()).toContainText('->>');
    }
  });
});

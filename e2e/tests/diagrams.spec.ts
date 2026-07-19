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

  test('interact with floating toolbar and collapsible code editor', async ({ page }) => {
    await page.goto('/mermaid-editor');

    // 1. Verify floating toolbar is present by looking for "Add Shape"
    const addShapeBtn = page.getByRole('button', { name: 'Add Shape' });
    await expect(addShapeBtn).toBeVisible();

    // 2. Click Add Shape and verify modal appears with shape options
    await addShapeBtn.click();
    const actorOption = page.getByRole('button', { name: 'Actor', exact: true });
    await expect(actorOption).toBeVisible();

    // 3. Select a shape to add to canvas
    await actorOption.click();
    
    // Ensure the node is rendered in the Native Mermaid SVG canvas
    const svgNode = page.locator('.actor').first();
    await expect(svgNode).toBeVisible();

    // 4. Verify code preview panel starts expanded or collapsed, and toggle it
    const collapseCodeBtn = page.getByTestId('collapse-code-btn');
    await expect(collapseCodeBtn).toBeVisible();
    await collapseCodeBtn.click();
    
    // We expect the 'Mermaid Source' text to disappear since it collapsed
    await expect(page.getByText('Mermaid Source')).not.toBeVisible();

    // The collapsed button should appear
    const expandCodeBtn = page.getByTestId('expand-code-btn');
    await expect(expandCodeBtn).toBeVisible();
    
    // Click the button to expand
    await expandCodeBtn.click();
    await expect(page.getByText('Mermaid Source')).toBeVisible();
    
    // Verify Apply Code button is present
    const applyBtn = page.getByRole('button', { name: 'Apply Code' });
    await expect(applyBtn).toBeVisible();
  });

  test('drag to connect SVG nodes', async ({ page }) => {
    await page.goto('/mermaid-editor');

    // Add a first shape
    const addShapeBtn = page.getByRole('button', { name: 'Add Shape' });
    await addShapeBtn.click();
    await page.getByRole('button', { name: 'Actor', exact: true }).click();
    
    // Wait for native node to appear
    const firstNode = page.locator('.actor').first();
    await expect(firstNode).toBeVisible();

    // Click the node to reveal the connection handle
    await firstNode.click();
    const handle = page.locator('.connection-handle').first();
    await expect(handle).toBeVisible();

    // Drag from handle to empty space
    await handle.hover();
    await page.mouse.down();
    await page.mouse.move(500, 500); // move to arbitrary empty space
    await page.mouse.up();

    // Expect the shape palette modal to pop up at the cursor
    const floatingMenu = page.locator('.connection-menu');
    await expect(floatingMenu).toBeVisible();

    // Select a shape to create
    await floatingMenu.getByRole('button', { name: 'Actor', exact: true }).click();
    
    // Wait for layout recalculation - we should now have 2 nodes
    await expect(page.locator('.actor')).toHaveCount(2);

    // Apply code and check if the edge was inserted
    const applyBtn = page.getByRole('button', { name: 'Apply Code' });
    await expect(applyBtn).toBeVisible();
    await applyBtn.click();
    
    // We expect the mermaid code to contain the edge connection '->>'
    const editorContent = page.locator('.monaco-editor').first();
    await expect(editorContent).toContainText('->>');
  });
});

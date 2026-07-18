import { test, expect } from '@playwright/test';

test.describe('Kanban Features', () => {
  test('create a new task', async ({ page }) => {
    await page.goto('/');
    
    // Click the '+ Add Task' button (in the To Do column typically)
    await page.getByText('+ Add Task').first().click();

    // Type in the 'Task title...' input
    const taskInput = page.getByPlaceholder('Task title...');
    await taskInput.fill('My New Test Task');
    
    // Click Create Task
    await page.getByRole('button', { name: 'Create Task' }).click();

    // Verify task appears
    const newTask = page.getByText('My New Test Task');
    await expect(newTask).toBeVisible();
  });
});

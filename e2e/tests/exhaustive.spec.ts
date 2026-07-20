import { test, expect } from '@playwright/test';

const UNIQUE_ID = Date.now();
const EMAIL = `testuser_${UNIQUE_ID}@example.com`;
const PASSWORD = 'password123';

test.describe('Exhaustive E2E Test Suite', () => {

    test.describe('Fail Cases & Corner Cases', () => {
        // Run without the saved authentication state
        test.use({ storageState: { cookies: [], origins: [] } });

        test('Authentication Walls: redirects unauthenticated users', async ({ page }) => {
            // Attempt to access a protected route
            await page.goto('/projects/1');
            
            // Should redirect back to auth page
            await expect(page.getByPlaceholder('name@example.com')).toBeVisible();
            await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible();
        });

        test('Registration with invalid/existing email', async ({ page }) => {
            // Setup an existing user first using API to avoid UI dependency
            const apiContext = await page.context().request;
            await apiContext.post('/api/auth/register', {
                data: {
                    email: `existing_${UNIQUE_ID}@example.com`,
                    password: 'password123'
                }
            });
            // Clear the cookie that was just set by the API
            await page.context().clearCookies();

            await page.goto('/');
            
            // Toggle to register mode
            await page.getByText("Don't have an account? Sign up").click();
            
            // Fill existing email
            await page.getByPlaceholder('name@example.com').fill(`existing_${UNIQUE_ID}@example.com`);
            await page.getByPlaceholder('••••••••').fill('password123');
            await page.getByRole('button', { name: 'Sign Up' }).click();

            // Wait for error message
            await expect(page.getByText('If this email is not already registered', { exact: false })).toBeVisible();
        });

        test('Login with incorrect password', async ({ page }) => {
            await page.goto('/');
            
            await page.getByPlaceholder('name@example.com').fill(`existing_${UNIQUE_ID}@example.com`);
            await page.getByPlaceholder('••••••••').fill('wrongpassword');
            await page.getByRole('button', { name: 'Sign In' }).click();

            await expect(page.getByText('Invalid credentials')).toBeVisible();
        });
    });

    test.describe('Happy Path: Full User Journey', () => {
        // Run without the saved authentication state
        test.use({ storageState: { cookies: [], origins: [] } });

        test('Full workflow from registration to task linkage', async ({ page }) => {
            test.setTimeout(120000); // Give this long test plenty of time

            // 1. User Creation
            await page.goto('/');
            await page.getByText("Don't have an account? Sign up").click();
            await page.getByPlaceholder('name@example.com').fill(EMAIL);
            await page.getByPlaceholder('••••••••').fill(PASSWORD);
            await page.getByRole('button', { name: 'Sign Up' }).click();

            // 2. Login verification (should auto-login or land on dashboard)
            await expect(page.getByText('Vibe Command Center')).toBeVisible();
            await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible();

            // 3. Project Creation
            await page.getByRole('button', { name: 'New Project' }).click();
            
            // Fill the React Modal
            await page.getByPlaceholder('E.g., Website Redesign').fill('Exhaustive Project');
            await page.getByRole('button', { name: 'Create Project' }).click();
            
            // Wait for Exhaustive Project to appear in the sidebar
            await expect(page.getByRole('link', { name: 'Exhaustive Project' })).toBeVisible();

            // 4. Project Wiki Verification
            // Click on the project to expand its links
            await page.getByRole('link', { name: 'Exhaustive Project' }).click();
            await page.getByText('Project Wiki').click();

            // Verify the automatic "Index" page exists
            await expect(page.getByText('Welcome to your new project wiki')).toBeVisible();

            // 5. Project Article Creation
            await page.getByRole('button', { name: 'New Page' }).click();
            
            // In the Wiki Editor, there's a title input and a Monaco editor
            await page.getByPlaceholder('Page Title').fill('Technical Specs');
            await page.getByPlaceholder('Category (e.g. Engineering)').fill('Engineering');
            await page.getByPlaceholder('URL Slug (e.g. my-page)').fill('tech-specs');
            
            // Click on the monaco editor
            await page.locator('.monaco-editor').first().click();
            await page.keyboard.press('Control+A');
            await page.keyboard.press('Delete');
            await page.keyboard.type('# Technical Specifications\n\nThese are the specs.');
            
            await page.getByRole('button', { name: 'Save Page' }).click();
            
            // Verify creation
            await expect(page.getByText('These are the specs.')).toBeVisible();

            // 6. Task Creation & Linkage
            await page.getByRole('link', { name: 'Exhaustive Project' }).click(); // Navigate to Kanban
            await page.getByRole('button', { name: '+ Add Task' }).first().click();
            
            await page.getByPlaceholder('Task title...').fill('Implement Specs');
            await page.getByRole('button', { name: 'Write' }).click();
            await page.getByPlaceholder('Add more details using Markdown').fill('Check out [Tech Specs](/projects/2/wiki/tech-specs)');
            await page.getByRole('button', { name: 'Create Task' }).click();

            // Verify task appears
            await expect(page.getByText('Implement Specs')).toBeVisible();

            // 7. Global Wiki
            await page.getByText('General Wiki').click();
            await page.getByRole('button', { name: 'New Page' }).click();
            await page.getByPlaceholder('Page Title').fill('Global Onboarding');
            await page.getByPlaceholder('Category (e.g. Engineering)').fill('HR');
            await page.getByPlaceholder('URL Slug (e.g. my-page)').fill('hr-onboarding');
            
            await page.locator('.monaco-editor').first().click();
            await page.keyboard.press('Control+A');
            await page.keyboard.press('Delete');
            await page.keyboard.type('Welcome to the company!');
            await page.getByRole('button', { name: 'Save Page' }).click();

            await expect(page.getByText('Welcome to the company!')).toBeVisible();
        });
    });
});

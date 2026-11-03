import { test, expect } from "@playwright/test";

test.describe("Authentication Flows (Login & Register UI)", () => {
  test("halaman login memiliki form input dan validasi yang sesuai", async ({ page }) => {
    await page.goto("/login");

    // Login identifier (Email atau Username)
    const identifierInput = page.getByPlaceholder(/nama@email\.com atau username/i);
    await expect(identifierInput).toBeVisible();

    // Password input
    const passwordInput = page.locator("input[type='password'], input[autocomplete='current-password']").first();
    await expect(passwordInput).toBeVisible();

    // Submit button
    const submitBtn = page.locator("button[type='submit']").first();
    await expect(submitBtn).toBeVisible();

    // Verifikasi link ke pendaftaran
    const registerLink = page.locator("a[href='/register']").first();
    await expect(registerLink).toBeVisible();
  });

  test("halaman register memiliki seluruh field pendaftaran", async ({ page }) => {
    await page.goto("/register");

    const nameInput = page.getByPlaceholder(/Contoh: Budi Santoso/i);
    const usernameInput = page.getByPlaceholder(/budimancing/i);
    const emailInput = page.getByPlaceholder(/nama@email\.com/i);
    const passwordInput = page.locator("input[type='password']").first();

    await expect(nameInput).toBeVisible();
    await expect(usernameInput).toBeVisible();
    await expect(emailInput).toBeVisible();
    await expect(passwordInput).toBeVisible();

    const submitBtn = page.locator("button[type='submit']").first();
    await expect(submitBtn).toBeVisible();
  });
});

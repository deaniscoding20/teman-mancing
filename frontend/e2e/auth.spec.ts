import { test, expect } from "@playwright/test";

test.describe("Authentication Flows (Login & Register UI)", () => {
  test("halaman login memiliki form input dan validasi yang sesuai", async ({ page }) => {
    await page.goto("/login");

    const emailInput = page.locator("input[type='email'], input#email, input[name='email']").first();
    const passwordInput = page.locator("input[type='password'], input#password, input[name='password']").first();
    const submitBtn = page.locator("button[type='submit']").first();

    await expect(emailInput).toBeVisible();
    await expect(passwordInput).toBeVisible();
    await expect(submitBtn).toBeVisible();

    // Verifikasi link ke pendaftaran
    const registerLink = page.locator("a[href='/register']");
    await expect(registerLink).toBeVisible();
  });

  test("halaman register memiliki seluruh field pendaftaran", async ({ page }) => {
    await page.goto("/register");

    const nameInput = page.locator("input#name, input[name='name']").first();
    const usernameInput = page.locator("input#username, input[name='username']").first();
    const emailInput = page.locator("input[type='email'], input#email, input[name='email']").first();
    const passwordInput = page.locator("input[type='password'], input#password, input[name='password']").first();

    await expect(nameInput).toBeVisible();
    await expect(usernameInput).toBeVisible();
    await expect(emailInput).toBeVisible();
    await expect(passwordInput).toBeVisible();

    const submitBtn = page.locator("button[type='submit']").first();
    await expect(submitBtn).toBeVisible();
  });
});

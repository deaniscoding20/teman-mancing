import { test, expect } from "@playwright/test";

test.describe("Landing Page & Navigation", () => {
  test("harus menampilkan hero header dan tagline brand", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveTitle(/Teman Mancing/i);

    // Verifikasi branding Teman Mancing
    const brandHeading = page.locator("header");
    await expect(brandHeading).toContainText("Teman Mancing");

    // Verifikasi navigation links di navbar
    const spotsLink = page.locator("header a[href='/spots']");
    const eventsLink = page.locator("header a[href='/events']");
    const chatLink = page.locator("header a[href='/chat']");

    await expect(spotsLink).toBeVisible();
    await expect(eventsLink).toBeVisible();
    await expect(chatLink).toBeVisible();
  });

  test("tombol CTA Jelajahi Spot mengarahkan ke halaman /spots", async ({ page }) => {
    await page.goto("/");
    const ctaButton = page.locator("a[href='/spots']").first();
    await expect(ctaButton).toBeVisible();
    await ctaButton.click();
    await expect(page).toHaveURL(/\/spots/);
  });
});

import { test, expect } from "@playwright/test";

test.describe("Fishing Spots Exploration", () => {
  test("halaman direktori spot memuat judul dan search bar", async ({ page }) => {
    await page.goto("/spots");

    const pageHeading = page.locator("h1");
    await expect(pageHeading).toContainText(/Spot Mancing/i);

    const searchInput = page.locator("input[type='text']").first();
    await expect(searchInput).toBeVisible();
  });
});

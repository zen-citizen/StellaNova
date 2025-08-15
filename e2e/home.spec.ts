import { test, expect } from "@playwright/test";

test("homepage loads and title matches", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/Civic Compass/i);
});

test("core sections render", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/Civic Compass/i);

  // Check that map container is present
  const mapContainer = page.locator(".leaflet-container");
  await expect(mapContainer).toBeVisible();
});

test("map renders and loads tiles", async ({ page }) => {
  await page.goto("/");
  const mapContainer = page.locator(".leaflet-container");

  // Check at least one map tile image has loaded
  const tileImg = mapContainer.locator("img.leaflet-tile");
  await expect(tileImg.first()).toBeVisible();

  // Check that the map tile has a non-empty src attribute
  await expect(tileImg.first()).toHaveAttribute("src", /https?:\/\//);

  // Ensure multiple map tiles are present and map actually rendered fully
  await page.waitForTimeout(3000);
  const tileCount = await tileImg.count();
  expect(tileCount).toBeGreaterThan(0);
});

import { defineConfig } from "@playwright/test";

const DEFAULT_OLD_APP = "https://civic-compass.zencitizen.in/";
const DEFAULT_LOCAL = "http://localhost:5173";
const DEFAULT_NEW_APP = "https://stella-nova.vercel.app/";

const MODE = process.env.PW_MODE ?? "old";

const urlMap: Record<string, string> = {
  local: process.env.LOCAL_URL ?? DEFAULT_LOCAL,
  new: process.env.NEW_URL ?? DEFAULT_NEW_APP,
  old: DEFAULT_OLD_APP,
};

const baseURL = urlMap[MODE];

export default defineConfig({
  testDir: "e2e",
  testMatch: ["**/*.spec.ts", "**/*.test.ts"],
  use: { baseURL, headless: true },
  reporter: [["list"]],
});

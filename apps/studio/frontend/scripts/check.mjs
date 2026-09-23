import { readFile } from "node:fs/promises";
import path from "node:path";
import process from "node:process";

const rootDir = process.cwd();
const requiredFiles = [
  "package.json",
  "tsconfig.json",
  path.join("src", "app-shell.tsx"),
  path.join("src", "navigation.ts")
];

for (const relativePath of requiredFiles) {
  const absolutePath = path.join(rootDir, relativePath);
  await readFile(absolutePath, "utf8");
}

const shellSource = await readFile(path.join(rootDir, "src", "app-shell.tsx"), "utf8");
if (!shellSource.includes("export const STUDIO_MODULES")) {
  throw new Error("src/app-shell.tsx must export STUDIO_MODULES");
}

console.log("frontend skeleton check: OK");

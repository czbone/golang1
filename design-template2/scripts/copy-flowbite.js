const fs = require("fs");
const path = require("path");

const src = path.join(__dirname, "..", "node_modules", "flowbite", "dist", "flowbite.min.js");
const destDir = path.join(__dirname, "..", "static", "js");
const dest = path.join(destDir, "flowbite.min.js");

fs.mkdirSync(destDir, { recursive: true });
fs.copyFileSync(src, dest);

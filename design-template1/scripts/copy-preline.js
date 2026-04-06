const fs = require("fs");
const path = require("path");

const src = path.join(__dirname, "..", "node_modules", "preline", "dist", "preline.js");
const destDir = path.join(__dirname, "..", "static", "js");
const dest = path.join(destDir, "preline.js");

fs.mkdirSync(destDir, { recursive: true });
fs.copyFileSync(src, dest);

const esbuild = require("esbuild");

esbuild
  .build({
    entryPoints: ["./src/*.js"],
    bundle: true,
    minify: true,
    treeShaking: true,
    outdir: "./dist",
    target: ["chrome58", "firefox57", "safari11"],
  })
  .then(() => console.log("⚡Bundle build complete ⚡"))
  .catch(() => {
    process.exit(1);
  });

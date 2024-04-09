const esbuild = require('esbuild');

esbuild
  .build({
    entryPoints: ['./src/*.js'],
    bundle: true,
    minify: true,
    treeShaking: true,
    outdir: './dist',
    target: ['chrome58', 'firefox57', 'safari11'],
    define: {
      'process.env.CLERK_PUBLIC_KEY': `"${process.env.CLERK_PUBLISHABLE}"`,
    },
  })
  .then(() => console.log('⚡Bundle build complete ⚡'))
  .catch(() => {
    process.exit(1);
  });

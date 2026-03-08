import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import typescript from '@rollup/plugin-typescript';

export default {
  input: 'src/index.ts',
  output: [
    {
      file: 'dist/entscape.esm.js',
      format: 'esm',
      sourcemap: true,
    },
    {
      file: 'dist/entscape.umd.js',
      format: 'umd',
      name: 'entscape',
      sourcemap: true,
      globals: {
        cytoscape: 'cytoscape',
      },
    },
  ],
  external: ['cytoscape'],
  plugins: [
    resolve(),
    commonjs(),
    typescript({
      tsconfig: './tsconfig.json',
      declaration: true,
      declarationDir: 'dist',
    }),
  ],
};

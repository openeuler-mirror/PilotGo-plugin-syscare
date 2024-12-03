/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Mar 4 16:45:10 2024 +0800
 */
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.join(__dirname, 'src'), //在任何模块文件内部，可以使用__dirname变量获取当前模块文件所在目录的完整绝对路径。
    },
  },
  server: {
    port: 8080,
    https: false,
    cors: true,
    proxy: {
      '/plugin': {
        target: 'http://10.44.43.181:7432',
        changeOrigin: true,
        rewrite: (p) => p.replace('/plugin^', '')
      },
    },
  },
  base: '/plugin/syscare/',
  build: {
    outDir: 'dist',
    assetsDir: 'static/', //指定静态资源存放路径
    sourcemap: false, //是否构建source map 文件
    chunkSizeWarningLimit: 1500, //chunk 大小警告的限制，默认500KB
  },
});

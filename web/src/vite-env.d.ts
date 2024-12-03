/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Mar 4 16:45:10 2024 +0800
 */

declare module "*.vue";
interface Window {
    remount: any;
    unmount: any;
  readonly '__MICRO_APP_BASE_ROUTE__': string;
  __MICRO_APP_ENVIRONMENT__: any;
}

interface ImportMeta {
  env: {
    BASE_URL: string;
    // Add other environment variables here
  };
}

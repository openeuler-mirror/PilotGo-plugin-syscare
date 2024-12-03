/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Apr 2 11:07:30 2024 +0800
 */
export interface Patch{
  id: number;
  creatTime: string;
  ip: string;
  status: string;
  stdout: string;
  endTime: string;
  buildLog: string;
  // 构建值
  buildKernel?: string;
  buildDebugInfo?: string;
  patchType: string;
  patchDescription: string;
  patchs: string;
  // 热补丁包: 
  hotPatch: string;
  patchKernel: string;
  taskId: string;
}

// 热补丁环境
export interface BuildEnv {
  version: string;
  rpm: {
    debugInfo: string;
    srcRpm: string;
  }
}

// 添加补丁列表
export interface PatchForm {
  ip: string;
  buildKernelSrc: string;
  buildDebugInfo: string;
  version: string;
  patchVersion: string;
  patchDescription: string;
  patchType: string;
  patchRelease: string;
  patchs: any[],
}
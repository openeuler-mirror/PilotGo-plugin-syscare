/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Apr 2 11:07:30 2024 +0800
 */
// 主机type
export interface Host {
  connect: Boolean;
  cpuModelName: string;
  cpuNum: number;
  id: number;
  ip: string;
  kernelArch: string;
  kernelVersion: string;
  lastheartbeat: string;
  osVersion: string;
  platform: string;
  platformVersion: string;
  upTime: string;
  uuid: string;
}
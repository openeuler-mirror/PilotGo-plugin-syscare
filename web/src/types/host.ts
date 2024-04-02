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
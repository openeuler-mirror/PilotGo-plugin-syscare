/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Apr 2 11:07:30 2024 +0800
 */
import request from "./request";
/* 
/plugin/syscare/agentList
agent列表  get
/plugin/syscare/addAgent
添加主机 post  {“ip”：xxx}
/plugin/syscare/:ip
删除主机  delete   参数：ip

/plugin/syscare/buildEnv?ip=10.44.43.181    get  获取这个agent机器可制作热补丁的环境
*/

// 获取所有机器
export function getAllHost(data?:Object) {
  return request({
    url: "/plugin/syscare/agentList",
    method: "get",
    params:data
  });
}
// 获取制作热补丁环境
export function getBuildEnv(data:{ip: string}) {
  return request({
    url: "/plugin/syscare/buildEnv",
    method: "get",
    params: data
  });
}

// 新增机器
export function addHost(data: {ip:string}) {
  return request({
    url: "/plugin/syscare/addAgent",
    method: "post",
    data,
  });
}


// 删除机器
export function delHost(ip:string ) {
  return request({
    url: "/plugin/syscare/"+ip,
    method: "delete",
  });
}

// 上传rpm文件
export function uploadRpm(data: any) {
  return request({
    url: "/plugin/syscare/upload",
    method: "post",
    data,
  });
}
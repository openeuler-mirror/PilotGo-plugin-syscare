/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Apr 2 11:07:30 2024 +0800
 */
import request from "./request";
// 获取所有补丁
export function getAllPatch(data:Object) {
  return request({
    url: "/plugin/syscare/lists",
    method: "get",
    params:data
  });
}

// 新增补丁
export function addPatch(data: any) {
  return request({
    url: "/plugin/syscare/addWarm",
    method: "post",
    data,
  });
}

// 编辑补丁
export function updatePatch(data: object) {
  return request({
    url: "/plugin/atune/task_new",
    method: "post",
    data,
  });
}

// 删除补丁
export function delPatch(data:{id: number}) {
  return request({
    url: "/plugin/syscare/delete",
    method: "get",
    params:data
  });
}
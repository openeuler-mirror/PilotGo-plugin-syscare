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
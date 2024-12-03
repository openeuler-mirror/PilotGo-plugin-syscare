/* 
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Tue Apr 2 11:07:30 2024 +0800
 */
import { createRouter, createWebHistory } from "vue-router";
import type { RouteRecordRaw } from "vue-router";
import HostList from "@/views/hostList.vue";
import PatchList from "@/views/patchList.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "",
    redirect: "/hostList",
  },
  {
    path: "/hostList",
    component: HostList,
    meta: { title: "机器列表" }
  },
  {
    path: "/patchList",
    component: PatchList,
    meta: { title: "补丁制作" }
  },
];
const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
export function directTo(to: any) {
  router.push(to)
}

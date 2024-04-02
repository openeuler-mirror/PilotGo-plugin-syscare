import { createRouter, createWebHistory,createWebHashHistory } from "vue-router";
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
  history: createWebHashHistory(),
  routes,
});

export default router;
export function directTo(to: any) {
  router.push(to)
}

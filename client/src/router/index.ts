import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import Home from "@/views/Home.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/game",
    name: "game",
    component: () => import(/* webpackChunkName: "about" */ "../views/Game.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;

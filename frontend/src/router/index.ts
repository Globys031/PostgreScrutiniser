import { createRouter, createWebHistory } from "vue-router";
import HomeView from "../views/Home.vue";
import LoginView from "../views/Login.vue";

import { useSessionStore } from "../stores/session";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
      meta: { requiresAuth: true },
    },
    {
      path: "/login",
      name: "login",
      component: LoginView,
      meta: { requiresAuth: false },
    },
    {
      path: "/about",
      name: "about",
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("../views/About.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/docs",
      name: "docs",
      component: () => import("../views/Docs.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/configurations",
      name: "configurations",
      component: () => import("../views/RuntimeConfigs.vue"),
      meta: { requiresAuth: true },
    },
    {
      path: "/backups",
      name: "backups",
      component: () => import("../views/ConfigBackups.vue"),
      meta: { requiresAuth: true },
    },
  ],
});

router.beforeEach((to) => {
  const sessionStore = useSessionStore();
  if (to.meta.requiresAuth && !sessionStore.token) return "/login";
  if (sessionStore.token && to.name == "login") return "/";
});

export default router;

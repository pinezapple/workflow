import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/workflow",
    name: "Workflows",
    component: () =>
      import(
        /* webpackChunkName: "about" */ "../views/workflow/workflow-list.vue"
      ),
  },
  {
    path: "/projects",
    name: "Projects",
    component: () =>
      import(
        /* webpackChunkName: "about" */ "../views/projects/project-list.vue"
      ),
  },
  {
    path: "/analyses",
    name: "Analyses",
    component: () =>
      import(
        /* webpackChunkName: "about" */ "../views/analyses/analytics-list.vue"
      ),
  },
  {
    path: "/biosample",
    name: "Biosample",
    component: () =>
      import(
        /* webpackChunkName: "about" */ "../views/biosample/biosample-list.vue"
      ),
  },
  {
    path: "/demodata",
    name: "Demodata",
    component: () =>
      import(
        /* webpackChunkName: "about" */ "../views/demodata/demodata-list.vue"
      ),
  },
  {
    path: "/upload",
    name: "Upload",
    component: () => import("../views/biodata/upload/upload.vue"),
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;

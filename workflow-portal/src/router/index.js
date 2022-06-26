import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/projects",
    name: "Projects",
    component: () =>
      import(
        /* webpackChunkName: "project-list" */ "../views/projects/project-list.vue"
      ),
  },
  {
    path: "/projects/:projectId/data/:path*",
    name: "Project Detail",
    component: () =>
      import(
        /* webpackChunkName: "project-detail" */ "../views/project/project-detail.vue"
      ),
  },
  {
    path: "/projects/:projectId/workflows/:workflowId*",
    name: "Workflow Edit",
    component: () =>
      import(
        /* webpackChunkName: "workflow-edit" */ "../views/workflow/workflow-edit.vue"
      ),
  },
  {
    path: "/workflows",
    name: "Workflows",
    component: () =>
      import(
        /* webpackChunkName: "workflow-list" */ "../views/workflow/workflow-list.vue"
      ),
  },
  {
    path: "/analyses",
    name: "Analyses",
    component: () =>
      import(
        /* webpackChunkName: "analytic-list" */ "../views/runs/analytics-list.vue"
      ),
  },
  {
    path: "/analyses/:runId",
    name: "Analyse Detail",
    component: () =>
      import(
        /* webpackChunkName: "run-detail" */ "../views/runs/run-detail.vue"
      ),
  },
  {
    path: "/tasks/:taskId",
    name: "Task Details",
    component: () =>
      import(
        /* webpackChunkName: "task-detail" */ "../views/task/task-detail.vue"
      ),
  },
  {
    path: "/log/:id",
    name: "Log Details",
    component: () =>
      import(
        /* webpackChunkName: "log-details" */ "../views/logs/log-detail.vue"
      ),
  },

  {
    path: "/projects/add",
    name: "Project Add",
    component: () =>
      import(
        /* webpackChunkName: "project-edit" */ "../views/project/project-edit.vue"
      ),
  },
  {
    path: "/biosample",
    name: "Biosample",
    component: () =>
      import(
        /* webpackChunkName: "biosample-list" */ "../views/biosample/biosample-list.vue"
      ),
  },
  {
    path: "/demodata",
    name: "Demodata",
    component: () =>
      import(
        /* webpackChunkName: "demodata-list" */ "../views/demodata/demodata-list.vue"
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
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;

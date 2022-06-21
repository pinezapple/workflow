import axios from "axios";
//ELEMENT UI
import ElementPlus from "element-plus";
import "element-plus/lib/theme-chalk/index.css";
import { createApp } from "vue";
import VueAxios from "vue-axios";
import App from "./App.vue";
import router from "./router";
import store from "./store";

createApp(App)
  .use(store)
  .use(router)
  .use(VueAxios, axios)
  .use(ElementPlus)
  .mount("#app");

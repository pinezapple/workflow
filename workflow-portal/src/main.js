import { dom, library } from '@fortawesome/fontawesome-svg-core';
import { fab } from '@fortawesome/free-brands-svg-icons';
import { far } from '@fortawesome/free-regular-svg-icons';
import { fas } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import axios from "axios";
//ELEMENT UI
import ElementPlus from "element-plus";
import "element-plus/lib/theme-chalk/index.css";
import moment from "moment";
import { createApp } from "vue";
import VueAxios from "vue-axios";
import VueApexCharts from "vue3-apexcharts";
import App from "./App.vue";
import router from "./router";
import store from "./store";

library.add(fas);
library.add(fab);
library.add(far);
dom.watch();

const app = createApp(App);
app
  .component('font-awesome-icon', FontAwesomeIcon)
  .use(store)
  .use(router)
  .use(VueAxios, axios)
  .use(ElementPlus)
  .use(VueApexCharts)
  .mount("#app");

app.config.globalProperties.$filters = {
  datetime(value) {
    if (value) {
      return moment(value).format("DD/MM/YYYY hh:mm");
    }
  },

  fileSize(bytes, decimals = 2) {
    if (bytes == 0) return "0 Bytes";
    var k = 1024,
      dm = decimals || 2,
      sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"],
      i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  },
};

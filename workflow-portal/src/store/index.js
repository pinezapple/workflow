// import createPersistedState from "vuex-persistedstate";
import { createStore } from "vuex";
import analyses from "./modules/analyses";
import workflow from "./modules/workflow";

// Create a new store instance.
const store = createStore({
  // plugins: [createPersistedState()],
  modules: {
    workflow,
    analyses,
  },
});

export default store;

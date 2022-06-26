// import createPersistedState from "vuex-persistedstate";
import { createStore } from "vuex";

// Load all modules.
function loadModules() {
  const context = require.context("./modules", false, /([a-z_]+)\.js$/i);

  const modules = context
    .keys()
    .map((key) => ({
      key,
      name: key.match(/([a-z_]+)\.js$/i)[1],
      namespaced: true,
    }))
    .reduce((modules, { key, name }) => {
      context(key).default.namespaced = true;
      return {
        ...modules,
        [name]: context(key).default,
      };
    }, {});

  return { context, modules };
}

const { context, modules } = loadModules();

const store = createStore({
  strict: true,
  // plugins: [createPersistedState()],
  modules,
});

export default store;

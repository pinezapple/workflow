import edda from "../../services/edda";

const state = {
  log: null,
};

const getters = {
  getLog: (state) => state.log,
};

const actions = {
  async getLogDetails({ commit }, id) {
    const response = await edda.getLog(id, 0);
    await commit("setLog", response);
  },
};

const mutations = {
  setLog(state, data) {
    state.log = data;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

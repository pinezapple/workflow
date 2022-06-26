import heimdall from "../../services/heimdall";

const state = {
  analyses: [],
  selectedAnalyse: null,
  total: 0,
  currentPage: 1,
  pageSize: 10,
  isFetch: false,
  error: null,
};

const getters = {
};

const actions = {
  async ChangePageSize({ commit, dispatch }, pageSize) {
    await commit("setPageSize", pageSize);
    await dispatch("GetAnalyses");
  },

  async ChangePage({ commit, dispatch }, page) {
    await commit("setPage", page);
    await dispatch("GetAnalyses");
  },

  async GetAnalyses({ commit, state }) {
    let response = await heimdall.getRuns(state.currentPage, state.pageSize);
    await commit("setAnalyses", response);
  },

  async GetSelectedAnalyse({ commit }, itemId) {
    const data = await heimdall.getRun(itemId);
    await commit("setSeletedAnalyse", data);
  },

  async RunAnalyse({ commit }, payload) {
    const data = await heimdall.createRun(payload);
    console.log("Run data", data);
  }
};

const mutations = {
  setPage(state, page) {
    state.currentPage = page;
  },

  setPageSize(state, pageSize) {
    state.pageSize = pageSize;
  },

  setAnalyses(state, data) {
    state.analyses = data.runs;
    state.total = data.total;
  },

  setSeletedAnalyse(state, analyse) {
    state.selectedAnalyse = analyse;
  },

  clearSelectedAnalyse(state) {
    state.selectedAnalyse = null;
  }
};

export default {
  state,
  getters,
  actions,
  mutations,
};

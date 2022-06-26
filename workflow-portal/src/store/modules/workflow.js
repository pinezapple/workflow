import heimdall from "../../services/heimdall";

const state = {
  workflows: [],
  selectedWorkflow: null,
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
    await dispatch("GetWorkflows");
  },

  async ChangePage({ commit, dispatch }, currentPage) {
    await commit("setCurrentPage", currentPage);
    await dispatch("GetWorkflows");
  },

  async GetWorkflows({ commit, state }) {
    let response = await heimdall.getWorkflows(state.currentPage, state.pageSize);
    await commit("setWorkflows", response);
  },

  async SetSelectedWorkflow({ commit }, id) {
    await commit("setSelectedWorkflow", id);
  },

}

const mutations = {
  setCurrentPage(state, currentPage) {
    state.currentPage = currentPage;
  },

  setPageSize(state, pageSize) {
    state.pageSize = pageSize;
  },

  setWorkflows(state, data) {
    state.workflows = data.workflows;
    state.total = data.total;
  },

  getSeletedWorkflow(state, itemsId) {
    state.indexSelect = itemsId;
  },

  setSelectedWorkflow(state, workflowId) {
    state.selectedWorkflow = state.workflows.find(item => item.id === workflowId);
  },

}

export default {
  state,
  getters,
  actions,
  mutations,
};

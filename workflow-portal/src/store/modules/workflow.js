import request from "../../services/heimdall";

const state = {
  workflowList: {
    tableData: [],
    multipleSelection: [],
  },
  total: null,
  page: 1,
  page_size: 20,
  indexSelect: null,
  itemsSelected: null,
  isFetch: false,
  error: null,
};

const getters = {
  stateWorkflowList: (state) => state.workflowList,
  stateIndexSelected: (state) => state.indexSelect,
  stateItemsSelected: (state) => state.itemsSelected,
  stateIsFetch: (state) => !!state.isFetch,
  stateIsError: (state) => state.error,
  getPage: (state) => state.page,
  getPageSize: (state) => state.page_size,
  getTotal: (state) => state.total,
};

const actions = {
  async ChangePageSize({ commit, dispatch }, page_size) {
    await commit("setPageSize", page_size);
    await dispatch("GetWorkflowList");
  },

  async ChangePage({ commit, dispatch }, page) {
    await commit("setPage", page);
    await dispatch("GetWorkflowList");
  },

  async GetWorkflowList({ commit, state }) {
    let response = await request.getWorkflows(state.page, state.page_size);
    await commit("setWorkflowList", response);
  },

  async GetItemsSelected({ commit }, itemsId) {
    await commit("getSeletedWorkflow", itemsId);
    await commit("getItemsWorkflow", itemsId);
  },
};

const mutations = {
  setPage(state, page) {
    state.page = page;
  },
  setPageSize(state, page_size) {
    state.page_size = page_size;
  },

  setWorkflowList(state, data) {
    state.workflowList.tableData = data.workflows;
    state.total = data.total;
  },

  getSeletedWorkflow(state, itemsId) {
    state.indexSelect = itemsId;
  },

  getItemsWorkflow(state, itemsId) {
    state.itemsSelected = state.workflowList.tableData.find(
      (wf) => wf.uuid === itemsId
    );
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

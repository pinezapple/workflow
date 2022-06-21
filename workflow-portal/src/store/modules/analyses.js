import request from "../../services/heimdall";

const state = {
  analysesList: {
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
  stateAnalysesList: (state) => state.analysesList,
  stateIndexSelected: (state) => state.indexSelect,
  stateItemsAnalysesSelected: (state) => state.itemsSelected,
  stateIsFetch: (state) => !!state.isFetch,
  stateIsError: (state) => state.error,
  getAnalysesPage: (state) => state.page,
  getAnalysesPageSize: (state) => state.page_size,
  getAnalysesTotal: (state) => state.total,
};

const actions = {
  async ChangePageSizeAnalyses({ commit, dispatch }, page_size) {
    await commit("setPageSize", page_size);
    await dispatch("GetAnalyses");
  },

  async ChangePageAnalyses({ commit, dispatch }, page) {
    await commit("setPage", page);
    await dispatch("GetAnalyses");
  },

  async GetAnalyses({ commit, state }) {
    let response = await request.getRuns(state.page, state.page_size);
    await commit("setAnalysesList", response);
  },

  async GetItemsSelectedAnalyses({ commit }, itemsId) {
    await commit("getSeletedAnalyses", itemsId);
    await commit("getItemsAnalyses", itemsId);
  },
};

const mutations = {
  setPage(state, page) {
    state.page = page;
  },
  setPageSize(state, page_size) {
    state.page_size = page_size;
  },

  setAnalysesList(state, data) {
    state.analysesList.tableData = data.runs;
    state.total = data.total;
  },

  getSeletedAnalyses(state, itemsId) {
    state.indexSelect = itemsId;
  },

  getItemsAnalyses(state, itemsId) {
    state.itemsSelected = state.analysesList.tableData.find(
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

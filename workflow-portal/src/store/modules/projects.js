import request from "../../services/heimdall";

const state = {
  projects: [],
  selectedProject: null,
  total: 0,
  currentPage: 1,
  pageSize: 10,
  isFetch: false,
  error: null,
};

const getters = {
};

const actions = {
  async ChangePageSize({ commit, dispatch }, size) {
    await commit("setPageSize", size);
    await dispatch("GetProjectList");
  },

  async ChangePage({ commit, dispatch }, page) {
    await commit("setPage", page);
    await dispatch("GetProjectList");
  },

  async GetProjectList({ commit, state }) {
    let response = await request.getProjects(state.currentPage, state.pageSize);
    await commit("setProjectList", response);
  },

  async SetSelectedProject({ commit }, projectId) {
    await commit("setSeletedProject", projectId);
    await commit("getItemsProject", projectId);
  },

  async CreateProject(context, project) {
    await request.createProject(project);
  },
};

const mutations = {
  setPage(state, page) {
    state.currentPage = page;
  },
  setPageSize(state, page_size) {
    state.pageSize = page_size;
  },

  setProjectList(state, data) {
    state.projects = data.projects;
    state.total = data.total;
  },

  setSeletedProject(state, projectId) {
    state.selectedProject = state.projects.find(project => project.id === projectId);
  },

  clearSelectedProject(state) {
    state.selectedProject = null;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

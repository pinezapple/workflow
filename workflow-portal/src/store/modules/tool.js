import brokkr from "../../services/brokkr";

const state = {
  categories: [],
  tools: [],
  tool: null,
  fetching: false,
  error: null,
};

const getters = {
  categories: (state) => state.categories,
  tools: (state) => state.tools,
  tool: (state) => state.tool,
};

const actions = {
  async GetAllTools({ commit }) {
    const tools = await brokkr.getAllTools();
    await commit("setTools", tools);
  },
  async GetTools({ commit }) {
    const tools = await brokkr.getTools();
    await commit("setTools", tools);
  },
  async GetTool({ commit }, { id }) {
    const tool = await brokkr.getTool(id);
    await commit("setTool", tool);
  },
};

const mutations = {
  setCategories(state, categories) {
    state.categories = categories;
  },
  setTools(state, tools) {
    state.tools = tools.tools;
  },
  setTool(state, tool) {
    state.tool = tool;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

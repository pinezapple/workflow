import heimdall from "../../services/heimdall";

const state = {
  run: null,
  tasks: [],
  selectedTask: null,
  isFetch: false,
  error: null,
};

const getters = {
};

const actions = {
  async GetTask({ commit }, taskId) {
    const task = await heimdall.getTask(taskId);
    await commit("setTask", task);
  },

  async GetRun({ commit }, runId) {
    let response = await heimdall.getRun(runId);
    await commit("setRun", response);
    await commit("setTasks", response.tasks);
  },

};

const mutations = {
  setRun(state, run) {
    state.run = run;
  },

  setTask(state, task) {
    state.selectedTask = task;
  },

  setTasks(state, tasks) {
    state.tasks = tasks;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

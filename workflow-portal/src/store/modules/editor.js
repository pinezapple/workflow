import heimdall from "@/services/heimdall";
import valkyrie from "@/services/valkyrie";
const YAML = require("yaml");
const { Document, YAMLMap, YAMLSeq } = require("yaml");

const state = {
  tools: [],
  toolLinks: [],
  workflowInputs: [],
  projectPage: {},
  showSelectData: false,
  selectedInput: null,
  selectedTool: null,
  displaySampleNameDialog: false,
  workflow: {
    name: "",
    id: null,
  }
};

const getters = {}

const actions = {
  async GetAllProjects({ commit, dispatch }) {
    const response = await heimdall.getProjects();
    await commit("setAllProjects", response);
  },

  async GetProjectFiles({ commit }, { projectId, currentFolder }) {
    const project = await heimdall.getProject(projectId);
    const projectFiles = await valkyrie.getProjectDataFiles(projectId, currentFolder);

    const folders = project.folders?.filter(item => item.path.includes(currentFolder) && item.path !== currentFolder && item.path.lastIndexOf("/") === currentFolder.length - 1).map(item => Object.assign({ project_id: projectId, class: "folder" }, item));
    const files = projectFiles.map(item => Object.assign({ class: "file" }, item));
    const items = files.concat(folders);

    commit("setProjectItems", items);
    commit("setProject", project);
  },

  async SaveWorkflow({ dispatch, commit, state }, { projectId, content, steps }) {
    const workflow = {
      id: state.workflow.id,
      name: state.workflow.name,
      project_id: projectId,
      content: content,
      class: "workflow",
      steps: steps
    }
    commit("setWorkflow", workflow);
    if (state.workflow.id) {
      const data = await heimdall.updateWorkflow(state.workflow.id, state.workflow);
      commit("setWorkflow", data);
      console.log("updated workflow: ", data)
    } else {
      const data = await heimdall.createWorkflow(workflow);
      commit("setWorkflow", data);
      console.log("created workflow: ", data)
    }
  },

  async GetWorkflow({ commit, rootState }, workflowId) {
    const workflow = await heimdall.getWorkflow(workflowId);
    commit("setWorkflow", workflow);
    commit("buildWorkflow", rootState);
  }
}

const mutations = {
  setAllProjects(state, projects) {
    state.projectPage = projects;
  },

  setCurrentFolder(state, currentFolder) {
    state.currentFolder = currentFolder;
  },

  setProjectItems(state, items) {
    state.projectItems = items;
  },

  setProject(state, project) {
    state.project = project
  },

  addEditorTool(state, tool) {
    const ntools = state.tools.length;
    state.tools.push({ tool_index: ntools, ...tool, selected_version: tool.versions[0] });
  },

  removeEditorTool(state, tool_index) {
    state.tools = state.tools.filter(tool => tool.tool_index !== tool_index);
    state.toolLinks = state.toolLinks.filter(link => link.to_tool_index !== tool_index && link.from_tool_index !== tool_index);
  },

  addToolLink(state, link) {
    state.toolLinks.push(link);
  },

  removeToolLink(state, { tool_index, input_id }) {
    state.toolLinks = state.toolLinks.filter(link => !(link.to_input.id === input_id && link.to_tool_index === tool_index));
  },

  showSelectableData(state, { tool, input }) {
    state.showSelectData = true;
    state.selectedTool = tool;
    state.selectedInput = input;
  },

  hideSelectableData(state) {
    state.showSelectData = false;
  },

  showSampleNameDialog(state) {
    state.displaySampleNameDialog = true;
  },

  hideSampleNameDialog(state) {
    state.displaySampleNameDialog = false;
  },

  addWorkflowInput(state, workflowInput) {
    state.workflowInputs.push(workflowInput);
  },

  removeWorkflowInput(state, { tool_index, input_id }) {
    console.log("remove input value: ", tool_index, input_id);
    state.workflowInputs = state.workflowInputs.filter(input => !(input.tool.tool_index === tool_index && input.input.id === input_id));
  },

  setWorkflowName(state, workflowName) {
    state.workflow.name = workflowName;
  },

  setWorkflow(state, workflow) {
    state.workflow = workflow;
  },

  buildWorkflow(state, rootState) {
    const editorTools = [];
    const toolLinks = [];
    const workflow = state.workflow;
    const workflowYaml = YAML.parse(workflow.content);
    const steps = workflowYaml.steps;

    for (const stepName in steps) {
      const toolName = stepName.split("_")[0];
      const toolIndex = stepName.split("_")[1];
      const tool = rootState.tool.tools.find(tool => tool.name === toolName);
      editorTools.push({ tool_index: parseInt(toolIndex), ...tool, selected_version: tool.versions[0] });
    }

    for (const stepName in steps) {
      const toToolName = stepName.split("_")[0];
      const toToolIndex = stepName.split("_")[1];
      const toTool = rootState.tool.tools.find(tool => tool.name === toToolName);
      const toVersion = toTool.versions[0];

      const stepInput = steps[stepName]["in"];
      for (const inputName in stepInput) {
        const inputValue = stepInput[inputName];
        if (typeof (inputValue) === "string") {
          if (inputValue.includes("/")) {
            const inputTool = inputValue.split("/")[0];
            const fromToolOutputName = inputValue.split("/")[1];

            const fromToolName = inputTool.split("_")[0];
            const fromToolIndex = inputTool.split("_")[1];
            const fromTool = rootState.tool.tools.find(tool => tool.name === fromToolName);
            const fromVersion = fromTool.versions[0];
            const fromToolOutput = fromVersion.outputs.find(output => output.abbrev === fromToolOutputName);
            const toToolInput = toVersion.inputs.find(input => input.abbrev === inputName);

            toolLinks.push({
              from_tool_index: parseInt(fromToolIndex),
              from_tool: fromTool,
              from_version: fromVersion,
              from_output: fromToolOutput,
              to_tool_index: parseInt(toToolIndex),
              to_tool: toTool,
              to_version: toVersion,
              to_input: toToolInput
            })
          }
        }
      }
    }

    state.toolLinks = toolLinks;
    state.tools = editorTools;
    Object.assign(state, { toolLinks: toolLinks, tools: editorTools });
  }
}

export default {
  state,
  getters,
  actions,
  mutations,
};

import { buildTree } from "@/utils/tree";
import { cloneDeep } from "lodash";
import heimdall from "../../services/heimdall";
import valkyrie from "../../services/valkyrie";

const state = {
  workflows: [],
  files: [],
  folders: [],
  project: null,
  // total: null,
  currentPage: 1,
  pageSize: 10,
  selectedItems: [],
  currentFolder: "/",
  isFetch: false,
  error: null,
};

const getters = {
  total: (state, getters) => (filters) => getters.projectItems(filters)?.length,
  projectPageItems: (state, getters) => (filters) => getters.projectItems(filters).slice((state.currentPage - 1) * state.pageSize,
    state.currentPage * state.pageSize),
  treeFolders: (state) =>
    buildTree(
      cloneDeep(state.folders),
      state.project?.name,
      state.project?.id
    ),
  projectItems: (state) => (filters) => {
    let result = [];

    // Take only direct children path
    if (!filters.includes("folder")) {
      const folders = state.folders.filter((folder) => {
        return (
          folder.path.indexOf(state.currentFolder) === 0 &&
          folder.path
            .slice(state.currentFolder.length + 1, folder.path.length)
            .indexOf("/") === -1 &&
          (folder.path[state.currentFolder.length] === "/" ||
            state.currentFolder.length === 1)
        );
      });

      // Transform project folders to data table
      for (const folder of folders) {
        result.push({
          id: folder.id,
          name: folder.name,
          class: "folder",
          path: folder.path,
          author: folder.author,
          created_at: folder.created_at,
          updated_at: folder.updated_at,
          project_id: state.project.id,
          project_name: state.project.name,
        });
      }
    }

    if (!filters.includes("workflow")) {
      for (const workflow of state.workflows) {
        if (
          workflow.path === state.currentFolder ||
          state.currentFolder === "/"
        )
          result.push({ ...workflow });
      }
    }

    if (!filters.includes("file")) {
      const project_name = state.project?.name;
      for (const file of state.files) {
        if (
          file.project_path === state.currentFolder ||
          state.currentFolder === "/"
        ) {
          result.push(
            Object.assign({ class: "file", project_name: project_name, }, file)
          );
        }
      }
    }
    return result;
  },
  getProjectItem: (state) => (id) => {
    const ws = state.workflows.filter((w) => w.id === id);
    if (ws.length > 0) return ws[0];

    const fs = state.folders.filter((folder) => folder.id === id);
    if (fs.length > 0) return fs[0];

    const files = state.files.filter((file) => file.id === id);
    if (files.length > 0) return files[0];
  },
};

const actions = {
  async ChangePageSize({ commit }, pageSize) {
    await commit("setPageSize", pageSize);
  },

  async ChangeCurrentPage({ commit }, page) {
    await commit("setCurrentPage", page);
  },

  async GetItemsSelected({ commit }, itemsId) {
    await commit("getSeletedItem", itemsId);
    await commit("getProjectItems", itemsId);
  },


  ChangeCurrentFolder({ commit }, currentFolder) {
    commit("setCurrentFolder", currentFolder);
  },

  async GetProjectWorkflows({ commit }, projectId) {
    commit("setWorkflows", []);
    const workflows = await heimdall.getProjectWorkflows(projectId);
    commit("setWorkflows", workflows);
  },

  async GetProjectDataFiles({ commit }, { projectId, currentFolder }) {
    commit("setFiles", []);
    let dataFiles = await valkyrie.getProjectDataFiles(
      projectId,
      currentFolder
    );
    if (dataFiles === null) dataFiles = [];
    // const files = dataFiles.map((file) => Object.assign({ projectId: projectId }, file));
    commit("setFiles", dataFiles);
  },

  async GetProject({ commit }, projectId) {
    const project = await heimdall.getProject(projectId);
    if (project.folders !== null) {
      for (let folder of project.folders) {
        folder.project_id = project.id;
        folder.project_name = project.name;
        folder.class = "folder";
      }
    }

    commit("setProject", project);
  },

  async AddProjectFolder({ commit, dispatch }, { projectId, folder, path }) {
    await heimdall.addProjectFolder(projectId, folder, path);
    dispatch("GetProject", projectId);
  },

  async DeleteProjectFolder({ commit }, { projectId, folderId }) {
    await heimdall.deleteProjectFolder(projectId, folderId);
    const project = await heimdall.getProject(projectId);
    commit("setProject", project);
  },

  async UpdateProjectPath({ dispatch, state }, { source, target }) {
    console.log("update ", source, target);
    if (source.class === "file" && target.class === "folder") {
      await valkyrie.updateFileFolder(source.id, target.path);
      dispatch("GetProjectDataFiles", {
        projectId: state.project.id,
        currentFolder: state.currentFolder,
      });
    } else if (source.class === "folder" && target.class === "folder") {
      console.log(
        "Update folder from " +
        source.path +
        " to " +
        (target.path + source.path)
      );

      const files = await valkyrie.getProjectDataFiles(state.project.id, source.path);
      console.log("list files: ", files)
      for (const file of files) {
        await valkyrie.updateFileFolder(file.id, target.path + source.path);
      }
      await heimdall.updateFolderPath(
        source.id,
        source.name,
        target.path + source.path,
        state.project.id
      );
      dispatch("GetProject", state.project.id);
    }

    // TODO(tuandn8) Need for case workflow move to folder

    // Finally update project
  },

  async DeleteProjectItems({ dispatch, state }, items) {
    console.log("delete: ", items);
    if (items.length === 0) return;
    const deleteType = items[0].class;

    for (const item of items) {
      if (item.class === "workflow" || item.class === "tool") {
        await heimdall.deleteWorkflow(item.id);
      } else if (item.class === "file") {
        await valkyrie.deleteFile(item.id);
      } else if (item.class === "folder") {
        await heimdall.deleteProjectFolder(state.project.id, item.id);
      }
    }

    if (deleteType === "workflow") {
      dispatch("GetProjectWorkflows", state.project.id);
    } else if (deleteType === "file") {
      dispatch("GetProjectDataFiles", { projectId: state.project.id, currentFolder: state.currentFolder })
    } else if (deleteType === "folder") {
      dispatch("GetProject", state.project.id);
    }
  }
};

const mutations = {
  setCurrentPage(state, page) {
    state.currentPage = page;
  },

  setPageSize(state, pageSize) {
    state.pageSize = pageSize;
  },

  setWorkflows(state, workflows) {
    state.workflows = workflows;
  },

  setCurrentFolder(state, currentFolder) {
    state.currentFolder = currentFolder;
  },

  setFiles(state, files) {
    files = files.map((file) =>
      Object.assign(
        {
          class: "file",
          author: file.owner,
          created_at: file.created_at,
          name: file.name,
          size: file.size,
        },
        file
      )
    );
    state.files = files;
  },

  setProject(state, project) {
    state.project = project;
    if (project.folders === null) state.folders = [];
    else state.folders = project.folders;
  },

  getSeletedItem(state, itemId) {
    state.indexSelect = itemId;
  },

  setSelectedItems(state, items) {
    state.selectedItems = state.workflows.filter((workflow) =>
      items.includes(workflow.id)
    );

    state.selectedItems = state.selectedItems.concat(
      state.files.filter((file) => items.includes(file.id))
    );

    state.selectedItems = state.selectedItems.concat(
      state.folders.filter((folder) => items.includes(folder.id))
    );
  },

  getProjectItems(state, itemsId) {
    state.selectedItems = state.workflows.find(
      (item) => item.id === itemsId
    );

    if (state.selectedItems == null) {
      if (state.project !== null && state.project.folders !== null) {
        state.selectedItems = state.project.folders.find(
          (item) => item.id === itemsId
        );
      }
    }
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

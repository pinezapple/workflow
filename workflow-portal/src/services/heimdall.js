import request from "../utils/handle-request";

export default {
  async getProjects(currentPage, pageSize) {
    const rs = await request.get("projects", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });
    return rs.data;
  },
  async getProject(projectId) {
    const rs = await request.get("projects/" + projectId);
    return rs.data;
  },
  async createProject(project) {
    const rs = await request.post("projects", project);
    return rs.data;
  },
  async updateProject(project) {
    const rs = await request.put("projects/" + project.project_uuid);
    return rs.data;
  },
  async deleteProject(projectId) {
    const rs = await request.delete("projects/" + projectId);
    return rs.data;
  },
  async getProjectWorkflows(projectId) {
    const rs = await request.get(
      "projects/" + projectId + "/workflows"
    );
    return rs.data.workflows ? rs.data.workflows : [];
  },

  async addProjectFolder(projectId, folder, path) {
    const rs = await request.post(
      "projects/" + projectId + "/folders",
      {
        name: folder,
        path: path,
      }
    );
    return rs.data;
  },

  async deleteProjectFolder(projectId, folderId) {
    const rs = await request.delete(
      "projects/" + projectId + "/folders/" + folderId
    );
    return rs.data;
  },

  async updateFolderPath(id, name, path, projectId) {
    const rs = await request.put(
      "projects/" + projectId + "/folders",
      {
        id: id,
        name: name,
        path: path,
      }
    );
    return rs.data;
  },

  async getWorkflows(currentPage, pageSize) {
    const rs = await request.get("workflows", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });

    return rs.data;
  },

  async getWorkflow(workflowId) {
    const rs = await request.get("workflows/" + workflowId);
    return rs.data;
  },

  async createWorkflow(workflow) {
    const rs = await request.post("workflows", workflow);
    return rs.data;
  },

  async updateWorkflow(workflowId, workflow) {
    const rs = await request.put(
      "workflows/" + workflowId,
      workflow
    );
    return rs.data;
  },

  async deleteWorkflow(workflowId) {
    const rs = await request.delete("workflows/" + workflowId);
    return rs.data;
  },

  async getRunsOfWorkflow(workflowId) {
    const rs = await request.get(
      "/workflows/" + workflowId + "/runs"
    );
    return rs.data;
  },

  async getRuns(currentPage, pageSize) {
    const rs = await request.get("runs", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });

    return rs.data;
  },

  async getRun(runId) {
    const rs = await request.get("runs/" + runId);
    return rs.data;
  },

  async createRun(run) {
    const rs = await request.post("/runs", run);
    return rs.data;
  },

  async getRunStatus(runId) {
    const rs = await request.get("/runs" + runId + "/status");
    return rs.data;
  },

  async deleteRun(runId) {
    const rs = await request.delete("/runs/" + runId);
    return rs.data;
  },

  async getTasks(currentPage, pageSize) {
    const rs = await request.get("/runs", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });
    return rs.data;
  },

  async getTask(taskId) {
    const rs = await request.get("/tasks/" + taskId);
    return rs.data;
  },
};

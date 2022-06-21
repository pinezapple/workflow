import request from "../utils/handle-request";

export default {
  async getWorkflows(currentPage, pageSize) {
    const rs = await request.get("heimdall/workflows", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });

    return rs.data;
  },

  async getWorkflow(workflowID) {
    const rs = await request.get("heimdall/workflows/" + workflowID);
    return rs.data;
  },

  async createWorkflow(workflow) {
    const rs = await request.post("heimdall/workflows", workflow);
    return rs.data;
  },

  async updateWorkflow(workflow_uuid, workflow) {
    const rs = await request.put(
      "heimdall/workflows" + workflow_uuid,
      workflow
    );
    return rs.data;
  },

  async deleteWorkflow(workflow_uuid) {
    const rs = await request.delete("heimdall/workflows/" + workflow_uuid);
    return rs.data;
  },

  async getRunsOfWorkflow(workflow_uuid) {
    const rs = await request.get(
      "/heimdall/workflows/" + workflow_uuid + "/runs"
    );
    return rs.data;
  },

  async getRuns(currentPage, pageSize) {
    const rs = await request.get("heimdall/runs", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });

    return rs.data;
  },

  async getRun(runUUID) {
    const rs = await request.get("heimdall/runs/" + runUUID);
    return rs.data;
  },

  async createRun(run) {
    const rs = await request.post("/heimdall/runs", run);
    return rs.data;
  },

  async getRunStatus(run_uuid) {
    const rs = await request.get("/heimdall/runs" + run_uuid + "/status");
    return rs.data;
  },

  async deleteRun(run_uuid) {
    const rs = await request.delete("/heimdall/runs/" + run_uuid);
    return rs.data;
  },

  async getTasks(currentPage, pageSize) {
    const rs = await request.get("/heimdall/runs", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });
    return rs.data;
  },

  async getTask(task_uuid) {
    const rs = await request.get("/heimdall/tasks/" + task_uuid);
    return rs.data;
  },
};

import request from "../utils/handle-request";

export default {
  // *************************** CATEGORY *****************************8
  async getCategories(currentPage, pageSize) {
    const rs = await request.get("brokkr/categories", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });

    return rs.data;
  },

  async getCategory(category_id) {
    const rs = await request.get("brokkr/categories/" + category_id);
    return rs.data;
  },

  async updateCategory(category_id, category) {
    const rs = await request.put("brokkr/categories/" + category_id, category);
    return rs.data;
  },

  async addCategory(category) {
    const rs = await request.post("brokkr/categories", category);
    return rs.data;
  },

  async deleteCategory(category_id) {
    const rs = await request.delete("brokkr/categories/" + category_id);
    return rs.data;
  },

  async getCategoryTools(category_id) {
    const rs = await request.get("brokkr/categories/" + category_id + "/tools");
    return rs.data;
  },

  // **************************** TOOL *************************************
  async getAllTools() {
    const rs = await request.get("brokkr/tools");
    return rs.data;
  },
  
  async getTools(currentPage, pageSize) {
    const rs = await request.get("brokkr/tools", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });

    return rs.data;
  },

  async getTool(tool_id) {
    const rs = await request.get("brokkr/tools/" + tool_id);
    return rs.data;
  },

  async updateTool(tool_id, tool) {
    const rs = await request.put("brokkr/tools/" + tool_id, tool);
    return rs.data;
  },

  async addTool(tool) {
    const rs = await request.post("brokkr/tools", tool);
    return rs.data;
  },

  async deleteTool(tool_id) {
    const rs = await request.delete("brokkr/tools/" + tool_id);
    return rs.data;
  },

  async getToolVersions(tool_id) {
    const rs = await request.get("brokkr/tools/" + tool_id + "/version");
    return rs.data;
  },

  // ***************************** VERSION ***********************************
  async getVersion(version_id) {
    const rs = await request.get("brokkr/versions/" + version_id);
    return rs.data;
  },

  async addVersion(version) {
    const rs = await request.post("brokkr/versions", version);
    return rs.data;
  },

  async updateVersion(version_id, version) {
    const rs = await request.put("brokkr/versions/" + version_id, version);
    return rs.data;
  },

  async deleteVersion(version_id) {
    const rs = await request.delete("brokkr/versions/" + version_id);
    return rs.data;
  },

  async deleteArgument(version_id, argument_id) {
    const rs = await request.delete(
      "brokkr/versions/" + version_id + "/args/" + argument_id
    );
    return rs.data;
  },

  async deleteInput(version_id, input_id) {
    const rs = await request.delete(
      "brokkr/versions/" + version_id + "/inputs/" + input_id
    );
    return rs.data;
  },

  async deleteOutput(version_id, output_id) {
    const rs = await request.delete(
      "brokkr/versions/" + version_id + "/outputs/" + output_id
    );
    return rs.data;
  },

  // *********************************** RUNTIME ***************************
  async getRuntimes(currentPage, pageSize) {
    const rs = await request.get("brokkr/runtimes", {
      params: {
        page_token: currentPage,
        page_size: pageSize,
      },
    });

    return rs.data;
  },

  async getRuntime(runtime_id) {
    const rs = await request.get("brokkr/runtimes/" + runtime_id);
    return rs.data;
  },

  async updateRuntime(runtime_id, runtime) {
    const rs = await request.put("brokkr/runtimes/" + runtime_id, runtime);
    return rs.data;
  },

  async addRuntime(runtime) {
    const rs = await request.post("brokkr/runtimes", runtime);
    return rs.data;
  },

  async deleteRuntime(runtime_id) {
    const rs = await request.delete("brokkr/runtimes/" + runtime_id);
    return rs.data;
  },
};

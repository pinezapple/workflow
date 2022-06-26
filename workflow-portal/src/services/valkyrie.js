import request from "../utils/handle-request";

export default {
  async getProjectDataFiles(projectId, currentFolder) {
    const rs = await request.get("files", {
      params: {
        filter: "project_id=" + projectId + ";project_path=" + currentFolder,
      },
    });

    return rs.data.Data;
  },

  async updateFileFolder(id, path) {
    console.log("update file path: ", id, path);
    const rs = await request.put("files/project_path", {
      path_files: [{ file_id: id, project_path: path }],
    });

    return rs.data;
  },

  async deleteFile(id) {
    console.log("delete file: ", id);
    const rs = await request.delete("files/file/" + id);
    return rs.data;
  }
};

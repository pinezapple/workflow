import request from "../utils/handle-request";

export default {
  async getLog(task_id, index) {
    const rs = await request.get("edda/log", {
      params: {
        task_uuid: task_id,
        index: index,
      },
    });

    return rs.data;
  },
};

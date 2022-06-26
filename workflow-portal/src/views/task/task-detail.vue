<template>
  <div class="task-container">
    <div class="task__icon">
      <i class="el-icon-circle-check" />
    </div>
    <div class="task__name">
      <div>
        <span>{{ task?.id }}</span>
        -
        <span>{{ task?.name }}</span>
      </div>
      <div class="task__state">
        <span>{{ task?.state }}</span>
      </div>
    </div>
  </div>
  <div class="task_details-box">
    <el-row :gutter="10">
      <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
        <task-detail />
      </el-col>
    </el-row>
  </div>
</template>
<script>
import TaskDetail from "@/components/task/TaskDetail";
import { useStore } from "vuex";
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";

export default {
  name: "TaskDetailView",
  components: {
    TaskDetail,
  },
  setup() {
    const store = useStore();
    const route = useRoute();

    const taskId = route.params.taskId;
    onMounted(() => {
      store.dispatch("run/GetTask", taskId);
    });

    return {
      task: computed(() => store.state.run.selectedTask),
    };
  },
};
</script>
<style lang="scss" scoped>
.task-container {
  display: flex;
  border-bottom: solid 1px #e5e5e5;
  border-radius: 4px;
  padding: 15px 25px;
  align-items: center;

  .task__icon {
    font-size: 2em;
    margin-right: 10px;
  }

  .task__name {
    font-size: 1em;

    .task__state {
      font-size: 0.8em;
      margin: 5px 0px;
      float: left;
    }
  }
}
</style>

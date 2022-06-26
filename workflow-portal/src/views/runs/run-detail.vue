<template>
  <div class="run-container">
    <div class="run__icon">
      <i class="el-icon-circle-check" />
    </div>
    <div class="run__name">
      <div>
        <span>{{ runName }}</span>
      </div>
      <div class="run__state">
        <span>{{ runState }}</span>
      </div>
    </div>
  </div>
  <div class="task_details-box">
    <el-row :gutter="10">
      <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
        <run-detail />
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { onMounted, computed } from "vue";
import { useRoute } from "vue-router";
import { useStore } from "vuex";
import RunDetail from "@/components/runs/RunDetail";

export default {
  name: "RunDetailView",
  components: {
    RunDetail,
  },
  setup() {
    const store = useStore();
    const route = useRoute();

    const runId = route.params.runId;
    onMounted(() => {
      store.dispatch("run/GetRun", runId);
    });
    return {
      runName: computed(() => store.state.run.run?.id),
      runState: computed(() => store.state.run.run?.state),
    };
  },
};
</script>

<style lang="scss" scoped>
.run-container {
  display: flex;
  border-bottom: solid 1px #e5e5e5;
  border-radius: 4px;
  padding: 15px 25px;
  align-items: center;

  .run__icon {
    font-size: 2em;
    margin-right: 10px;
  }

  .run__name {
    font-size: 1em;

    .run__state {
      font-size: 0.8em;
      margin: 5px 0px;
      float: left;
    }
  }
}
</style>

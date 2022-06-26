<template>
  <div class="task_details-box">
    <el-row :gutter="10">
      <el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24">
        <LogDetails />
      </el-col>
    </el-row>
  </div>
</template>

<script>
// @ is an alias to /src
import LogDetails from "../../components/log/LogDetails";

export default {
  name: "LogDetailView",
  components: {
    LogDetails,
  },
};
</script>
<style></style>

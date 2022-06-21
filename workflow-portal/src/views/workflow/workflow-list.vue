<template>
  <div class="workflow-list">
    <el-row :gutter="10">
      <el-col :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
        <Table></Table>
      </el-col>
      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
        <LeftTable></LeftTable>
      </el-col>
    </el-row>
  </div>
</template>

<script>
// @ is an alias to /src
import Table from "../../components/workflow/Table";
import LeftTable from "../../components/workflow/LeftTable";

export default {
  name: "WorkflowList",
  components: {
    Table,
    LeftTable,
  },
};
</script>

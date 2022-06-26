<template>
  <el-table
    ref="multipleTable"
    :data="workflows"
    :default-sort="{ prop: 'created_at', order: 'descending' }"
    style="width: 100%"
    @selection-change="handleSelectionChange"
  >
    <el-table-column type="selection" width="55" />
    <el-table-column
      prop="name"
      sortable
      property="name"
      label="Name"
      width="300"
    >
      <template #default="{ row }">
        <router-link
          :to="{
            name: 'Workflow Edit',
            params: {
              workflowId: row.id,
              projectId: row.project_id,
            },
          }"
        >
          {{ row.name }}
        </router-link>
      </template>
    </el-table-column>
    <el-table-column
      prop="summary"
      sortable
      property="summary"
      label="Summary"
      show-overflow-tooltip
    />

    <el-table-column
      prop="description"
      sortable
      property="description"
      label="Description"
      show-overflow-tooltip
    />
    <el-table-column
      prop="author"
      sortable
      property="author"
      label="Author"
      show-overflow-tooltip
    />
    <el-table-column prop="created_at" sortable label="Date" width="200">
      <template #default="scope">
        {{ $filters.datetime(scope.row.created_at) }}
      </template>
    </el-table-column>
  </el-table>
  <el-pagination
    :page-sizes="[10, 20, 50, 100]"
    :page-size="pageSize"
    layout="total, sizes, prev, pager, next"
    :total="total"
    @size-change="setPageSize"
    @current-change="setCurrentPage"
  />
</template>

<script>
import { computed, onMounted } from "@vue/runtime-core";
import { useStore } from "vuex";

export default {
  name: "WorkflowsTable",
  components: {},
  setup() {
    const store = useStore();

    onMounted(() => {
      store.dispatch("workflow/GetWorkflows");
    });

    return {
      workflows: computed(() => store.state.workflow.workflows),
      pageSize: computed(() => store.state.workflow.pageSize),
      total: computed(() => store.state.workflow.total),
      currentPage: computed(() => store.state.workflow.currentPage),
      setCurrentPage: (page) => store.dispatch("workflow/ChangePage", page),
      setPageSize: (pageSize) =>
        store.dispatch("workflow/ChangePageSize", pageSize),
      clearSelectedWorkflow: () =>
        store.commit("workflow/clearSelectedWorkflow"),
      SetSelectedWorkflow: (id) =>
        store.dispatch("workflow/SetSelectedWorkflow", id),
    };
  },

  methods: {
    async handleSelectionChange(val) {
      if (val.length) {
        await this.SetSelectedWorkflow(val[0].id);
      } else {
        await this.SetSelectedWorkflow();
      }
    },
  },
};
</script>

<style>
.icon-more {
  margin-top: 10px;
}
</style>

<template>
  <el-table
    :data="analyses"
    :default-sort="{ prop: 'created_at', order: 'descending' }"
    style="width: 100%"
    @selection-change="handleSelectionChange"
  >
    <el-table-column type="selection" width="55" />
    <el-table-column prop="id" sortable property="id" label="ID" width="300">
      <template #default="{ row }">
        <router-link
          :to="{ name: 'Analyse Detail', params: { runId: row.id } }"
        >
          {{ row.id }}
        </router-link>
      </template>
    </el-table-column>
    <el-table-column
      prop="user"
      sortable
      property="user"
      label="User"
      width="250"
    />
    <el-table-column
      prop="state"
      sortable
      property="state"
      label="State"
      show-overflow-tooltip
    />
    <el-table-column
      prop="start_time"
      sortable
      property="start_time"
      label="Start Time"
      show-overflow-tooltip
    />
    <el-table-column
      prop="end_time"
      sortable
      property="end_time"
      label="End Time"
      show-overflow-tooltip
    />
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
import { computed, onMounted } from "vue";
import { useStore } from "vuex";

export default {
  name: "RunsTable",
  setup() {
    const store = useStore();

    onMounted(() => {
      store.dispatch("analyses/GetAnalyses");
    });

    return {
      analyses: computed(() => store.state.analyses.analyses),
      pageSize: computed(() => store.state.analyses.pageSize),
      currentPage: computed(() => store.state.analyses.currentPage),
      total: computed(() => store.state.analyses.total),
      setCurrentPage: (page) => store.dispatch("analyses/ChangePage", page),
      setPageSize: (pageSize) =>
        store.dispatch("analyses/ChangePageSize", pageSize),
      clearSelectedAnalyse: () => store.commit("analyses/clearSelectedAnalyse"),
      GetSelectedAnalyse: (id) =>
        store.dispatch("analyses/GetSelectedAnalyse", id),
    };
  },
  methods: {
    async handleSelectionChange(val) {
      if (val.length) {
        await this.GetSelectedAnalyse(val[0].id);
      } else {
        await this.clearSelectedAnalyse();
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

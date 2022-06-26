<template>
  <el-breadcrumb separator-class="el-icon-arrow-right" class="select__path">
    <el-breadcrumb-item class="select__item">
      All projects
    </el-breadcrumb-item>
  </el-breadcrumb>
  <el-table
    :data="projects"
    :default-sort="{ prop: 'created_at', order: 'descending' }"
    style="width: 100%; cursor: pointer"
    @row-click="handleClickProject"
  >
    <el-table-column
      prop="name"
      sortable
      property="name"
      label="Name"
      :min-width="25"
    >
      <template #default="scope">
        <div class="project__text">
          <i class="fas fa-briefcase project-item-icon" />
          <span>{{ scope.row.name }}</span>
        </div>
      </template>
    </el-table-column>
    <el-table-column
      sortable
      property="author"
      label="Author"
      show-overflow-tooltip
      :min-width="25"
    />
    <el-table-column prop="created_at" sortable label="Create at" width="200">
      <template #default="scope">
        {{ $filters.datetime(scope.row.created_at) }}
      </template>
    </el-table-column>
  </el-table>
  <el-pagination
    :page-sizes="[5, 10, 20]"
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
  name: "ProjectSelector",
  emits: ["selectProject"],
  setup() {
    const store = useStore();

    onMounted(() => {
      store.dispatch("projects/GetProjectList");
    });

    return {
      pageSize: computed(() => store.state.projects.pageSize),
      currentPage: computed(() => store.state.projects.currentPage),
      total: computed(() => store.state.projects.total),
      projects: computed(() => store.state.projects.projects),
      setCurrentPage: (page) => store.dispatch("projects/ChangePage", page),
      setPageSize: (size) => store.dispatch("projects/ChangePageSize", size),
    };
  },
  methods: {
    handleClickProject(row, column, event) {
      this.$emit("selectProject", { projectId: row.id, projectName: row.name });
    },
  },
};
</script>

<style lang="scss" scoped>
.select__path {
  margin: 0px 10px;
  padding-bottom: 10px;
  border-bottom: 2px solid #e9e9e9;

  .select__item {
    font-weight: 400;
    font-size: 18px;
    text-transform: capitalize;
  }
}
.project__text {
  font-weight: 600;
}
</style>

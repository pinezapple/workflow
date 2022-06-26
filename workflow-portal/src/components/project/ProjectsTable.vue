<template>
  <el-table
    :data="projects"
    :default-sort="{ prop: 'created_at', order: 'descending' }"
    style="width: 100%"
    @selection-change="handleSelectionChange"
  >
    <el-table-column type="selection" width="55" />
    <el-table-column prop="created_at" sortable label="Date" width="200">
      <template #default="scope">
        {{ $filters.datetime(scope.row.created_at) }}
      </template>
    </el-table-column>
    <el-table-column
      prop="name"
      sortable
      property="name"
      label="Name"
      :min-width="25"
    >
      <template #default="scope">
        <router-link
          :to="{ name: 'Project Detail', params: { projectId: scope.row.id } }"
        >
          {{ scope.row.name }}
        </router-link>
      </template>
    </el-table-column>
    <el-table-column
      prop="summary"
      sortable
      property="summary"
      label="Summary"
      :min-width="35"
      show-overflow-tooltip
    />
    <el-table-column
      sortable
      property="author"
      label="Author"
      show-overflow-tooltip
      :min-width="25"
    />
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
import { mapActions, useStore } from "vuex";
import { computed, onMounted, ref } from "vue";

export default {
  name: "ProjectsTable",
  setup() {
    const store = useStore();
    onMounted(() => store.dispatch("projects/GetProjectList"));

    return {
      pageSize: computed(() => store.state.projects.pageSize),
      currentPage: computed(() => store.state.projects.currentPage),
      total: computed(() => store.state.projects.total),
      setCurrentPage: (page) => store.dispatch("projects/ChangePage", page),
      setPageSize: (size) => store.dispatch("projects/ChangePageSize", size),
      projects: computed(() => store.state.projects.projects),
      clearSelectedProject: () => store.commit("projects/clearSelectedProject"),
    };
  },
  data() {
    return {
      form: {
        name: null,
        description: null,
        tags: [],
        summary: null,
      },
    };
  },
  methods: {
    ...mapActions("projects", ["SetSelectedProject"]),
    clearProjectForm() {
      this.form.name = null;
      this.form.description = null;
      this.form.summary = null;
      this.form.tags = [];
    },
    async handleSelectionChange(val) {
      if (val.length) {
        await this.SetSelectedProject(val[0].id);
      } else {
        await this.clearSelectedProject();
      }
    },
  },
};
</script>

<style>
.header_table-tool {
  padding: 5px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  border-bottom: 1px solid var(--violet-white-color);
}
.header_image-tool {
  margin-right: 10px;
}
.no_item-tool {
  margin-top: 50px;
}
.float-right {
  float: right;
}

.el-dialog__header {
  padding: 0px 20px;
}

.el-form-item__label {
  text-align: left;
}

.model-header {
  display: flex;
  border-bottom: 1px solid #e5e5e5;
}
</style>

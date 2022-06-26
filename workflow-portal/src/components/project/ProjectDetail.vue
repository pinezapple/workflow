<template>
  <div v-if="projectItems">
    <el-table
      :data="projectItems"
      :default-sort="{ prop: 'class', order: 'ascending' }"
      style="width: 100%"
      @selection-change="handleSelectRows"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column
        prop="name"
        sortable
        property="name"
        label="Name"
        width="400"
      >
        <template #default="{ row }">
          <div
            :draggable="['file', 'folder'].includes(row.class)"
            @dragstart="handleDragStart($event, row)"
            @drop="handleDrop($event, row)"
            @dragover="handleDragOverlap($event, row)"
            @dragenter="handleDragOverlap($event, row)"
          >
            <div v-if="row.class == 'workflow'">
              <em class="el-icon-s-operation project-item-icon" />
              <router-link
                :to="{
                  name: 'Workflow Edit',
                  params: {
                    projectId: row.project_id,
                    workflowId: row.id,
                  },
                }"
              >
                <span>{{ row.name }}</span>
              </router-link>
            </div>
            <div v-else-if="row.class == 'tool'">
              <em class="el-icon-s-tools project-item-icon" />
              <span>{{ row.name }}</span>
            </div>

            <div v-else-if="row.class == 'folder'">
              <router-link
                :to="{
                  name: 'Project Detail',
                  params: {
                    id: row.project_id,
                    path: row.path.slice(1, row.path.length),
                  },
                }"
              >
                <em class="el-icon-folder project-item-icon" />
                <span>{{ row.name }}</span>
              </router-link>
            </div>
            <div v-else-if="row.class == 'file'">
              <em class="el-icon-document project-item-icon" />
              <span>{{ row.name }}</span>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        prop="class"
        sortable
        property="class"
        label="Class"
        width="100"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <span v-if="row.class">
            {{ row.class[0].toUpperCase() + row.class.slice(1) }}
          </span>
          <span v-else>---</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="class"
        sortable
        property="class"
        label="Size"
        width="100"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <span v-if="row.size">
            {{ $filters.fileSize(row.size) }}
          </span>
          <span v-else>---</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="description"
        sortable
        property="description"
        label="Description"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <span v-if="row.description">
            {{ row.description }}
          </span>
          <span v-else>---</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="created_at"
        sortable
        label="Created Date"
        width="200"
      >
        <template #default="{ row }">
          <span v-if="row.created_at">
            {{ $filters.datetime(row.created_at) }}
          </span>
          <span v-else>---</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="author"
        sortable
        property="author"
        label="Author"
        show-overflow-tooltip
      >
        <template #default="{ row }">
          <span v-if="row.author">
            {{ row.author }}
          </span>
          <span v-else>---</span>
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
  </div>
</template>

<script>
import { computed } from "vue";
import { useStore } from "vuex";
export default {
  name: "ProjectDetail",
  emits: ["updateProjectPath"],
  setup() {
    const store = useStore();

    return {
      total: computed(() => store.getters["project/total"]([])),
      projectItems: computed(() =>
        store.getters["project/projectPageItems"]([])
      ),
      pageSize: computed(() => store.state.project.pageSize),
      setPageSize: (size) => store.dispatch("project/ChangePageSize", size),
      setCurrentPage: (page) =>
        store.dispatch("project/ChangeCurrentPage", page),
      getProjectItem: (id) => store.getters["project/getProjectItem"](id),
      UpdateProjectPath: (source, target) =>
        store.dispatch("project/UpdateProjectPath", {
          source: source,
          target: target,
        }),
      setSelectedItems: (items) =>
        store.commit("project/setSelectedItems", items),
    };
  },
  methods: {
    handleSelectRows(val) {
      const ids = val.map((item) => item.id);
      this.setSelectedItems(ids);
    },
    handleDragStart(event, item) {
      if (item.class === "file" || item.class === "folder") {
        event.dataTransfer.dropEffect = "move";
        event.dataTransfer.effectAllowed = "move";
        event.dataTransfer.setData("text/plain", item.id);
      }
    },
    handleDragOverlap(event) {
      event.preventDefault();
    },
    async handleDrop(event, item) {
      const sourceId = event.dataTransfer.getData("text/plain");
      const targetId = item.id;
      if (item.class === "folder") {
        const source = this.getProjectItem(sourceId);
        const target = this.getProjectItem(targetId);
        console.log("Source ", source, sourceId);
        console.log("Target ", target, targetId);
        await this.UpdateProjectPath(source, target);
      }
      event.preventDefault();
    },
  },
};
</script>

<style scoped>
.project-item-icon {
  margin: 10px 10px;
}
button {
  background-color: white;
  border: none;
}
</style>

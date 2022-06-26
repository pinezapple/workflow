<template>
  <div class="fileselect__path">
    <el-breadcrumb separator-class="el-icon-arrow-right" class="select__path">
      <el-breadcrumb-item
        v-for="item in paths"
        :key="item"
        class="select__item"
      >
        <i
          v-if="item.class === 'project'"
          class="fas fa-briefcase select_item project-item-icon"
        />
        <i
          v-else-if="item.class === 'folder'"
          class="el-icon-folder project-item-icon"
        />
        {{ item.name }}
      </el-breadcrumb-item>
    </el-breadcrumb>
  </div>
  <el-row>
    <el-col :span="6">
      <el-tree
        :data="treeFolders"
        :props="folderProps"
        default-expand-all
        highlight-current
        node-key="id"
        class="tree__folders"
      >
        <template #default="{ node, data }">
          <div class="folder__node" @click="handleClickFolder(data)">
            <div class="folder__content">
              <i class="el-icon-folder" />
              <div class="folder__text">
                <span>{{ data.name }}</span>
              </div>
            </div>
          </div>
        </template>
      </el-tree>
    </el-col>
    <el-col :span="18">
      <div v-if="projectItems">
        <el-table
          ref="multipleTable"
          :max-height="500"
          :data="projectItems"
          :default-sort="{ prop: 'class', order: 'ascending' }"
          style="width: 100%; cursor: pointer"
          @selection-change="handleSelectRows"
          @row-click="handleRowClick"
        >
          <el-table-column
            prop="name"
            sortable
            property="name"
            label="Name"
            width="300"
          >
            <template #default="scope">
              <div>
                <div v-if="scope.row.class == 'workflow'">
                  <em class="el-icon-s-operation project-item-icon" />
                  <span>{{ scope.row.name }}</span>
                </div>
                <div v-else-if="scope.row.class == 'tool'">
                  <em class="el-icon-s-tools project-item-icon" />
                  <span>{{ scope.row.name }}</span>
                </div>

                <div v-else-if="scope.row.class == 'folder'">
                  <div>
                    <em class="el-icon-folder project-item-icon" />
                    <span>{{ scope.row.name }}</span>
                  </div>
                </div>
                <div v-else-if="scope.row.class == 'file'">
                  <em class="el-icon-document project-item-icon" />
                  <span>{{ scope.row.name }}</span>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column
            prop="class"
            sortable
            property="class"
            label="Type"
            show-overflow-tooltip
          >
            <template #default="scope">
              <span v-if="scope.row.class">
                {{
                  scope.row.class[0].toUpperCase() + scope.row.class.slice(1)
                }}
              </span>
              <span v-else>---</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="class"
            sortable
            property="class"
            label="Size"
            show-overflow-tooltip
          >
            <template #default="scope">
              <span v-if="scope.row.size">
                {{ $filters.fileSize(scope.row.size) }}
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
            <template #default="scope">
              <span v-if="scope.row.created_at">
                {{ $filters.datetime(scope.row.created_at) }}
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
    </el-col>
  </el-row>
</template>

<script>
import { computed, onMounted, ref, toRefs } from "vue";
import { useStore } from "vuex";

export default {
  name: "FileSelector",
  props: {
    projectName: {
      type: String,
      default: "",
    },
  },
  emits: [
    "changePageSize",
    "changePage",
    "changeSelectRows",
    "updateProjectPath",
    "changeFolder",
    "selectFile",
  ],
  setup(props) {
    const store = useStore();
    const { projectName } = toRefs(props);
    const folderProps = { children: "children", label: "name" };
    const paths = ref([]);

    onMounted(() => {
      paths.value.push({
        name: "All Projects",
        path: "..",
      });

      paths.value.push({
        name: projectName,
        path: "/",
        class: "project",
      });
    });

    return {
      total: computed(() => store.getters["project/total"](["workflow"])),
      projectItems: computed(() =>
        store.getters["project/projectPageItems"](["workflow"])
      ),
      pageSize: computed(() => store.state.project.pageSize),
      projectFolders: computed(() => store.state.project.folders),
      setPageSize: (size) => store.dispatch("project/ChangePageSize", size),
      setCurrentPage: (page) =>
        store.dispatch("project/ChangeCurrentPage", page),
      folderProps,
      treeFolders: computed(() => store.getters["project/treeFolders"]),
      paths,
    };
  },
  methods: {
    handleRowClick(row) {
      if (row.class === "folder") {
        this.handleClickFolder(row);
      } else this.$emit("selectFile", row);
    },
    handleSelectRows(val) {
      const ids = val.map((item) => item.id);
      this.$emit("changeSelectRows", ids);
    },
    handleClickFolder(folder) {
      const path = folder.path;
      const pathSplited = path.split("/");
      const paths = this.paths.slice(0, 2);
      pathSplited.forEach((pathStep) => {
        console.log("Path step: ", pathStep);
        const folder = this.projectFolders.find(
          (folder) => folder.name === pathStep
        );
        if (folder) {
          paths.push({
            name: folder.name,
            path: folder.path,
            class: "folder",
          });
        }
      });
      this.paths = paths;

      this.$emit("changeFolder", {
        projectId: folder.project_id,
        currentFolder: folder.path,
      });
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

.tree__folders {
  margin-top: 10px;
}

.folder__node {
  padding: 10px 0;
  cursor: pointer;

  .folder__content {
    display: flex;
    padding: 2px 0 2px 2px;
    font-size: 18px;

    .folder__text {
      margin-left: 10px;
      font-size: 14px;
    }
  }
}
</style>

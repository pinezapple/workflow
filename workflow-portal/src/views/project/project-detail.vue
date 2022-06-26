<template>
  <title-bar>
    <template #left>
      <div class="header-container">
        <div class="header-title-container only-title">
          <h1 class="header-title">
            <span>Projects</span>
          </h1>
        </div>
      </div>
      <el-menu default-active="1" class="nav-menu">
        <el-menu-item index="">
          <router-link to="/projects" class="link">
            <span>SETTINGS</span>
          </router-link>
        </el-menu-item>
        <el-menu-item index="">
          <router-link to="/biosample" class="link">
            <span>MANAGE</span>
          </router-link>
        </el-menu-item>
        <el-menu-item index="">
          <router-link to="/biosample" class="link">
            <span>MONITOR</span>
          </router-link>
        </el-menu-item>
        <el-menu-item index="">
          <router-link to="/biosample" class="link">
            <span>VISUALIZE</span>
          </router-link>
        </el-menu-item>
      </el-menu>
    </template>
  </title-bar>
  <div class="main-section">
    <div class="menu-bar-container">
      <div class="menu-bar">
        <div class="breadcrumb" style="cursor: pointer">
          <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item
              v-for="item in navPaths"
              :key="item"
              :to="item.to"
            >
              <span class="menu-bar__title">
                <i
                  v-if="item.class === 'project'"
                  class="fas fa-briefcase project-item-icon"
                />
                <i
                  v-if="item.class === 'folder'"
                  class="el-icon-folder project-item-icon"
                />

                {{ item.name }}
              </span>
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="button-wrapper">
          <div v-if="selectedItems.length == 0">
            <el-button
              type="primary"
              class="action-button"
              @click="handleNewFolder"
            >
              <span class="el-icon-circle-plus" />
              <span>New folder</span>
            </el-button>
            <el-button class="action-button" @click="handleUploadData">
              <span class="el-icon-circle-plus" />
              <span>Upload data</span>
            </el-button>
            <el-button class="action-button" @click="handleNewWorkflow">
              <span class="el-icon-circle-plus" />
              <span>New workflow</span>
            </el-button>
          </div>
          <div
            v-else-if="
              selectedTypes.length == 1 &&
                (selectedTypes.includes('workflow') ||
                  selectedTypes.includes('tool'))
            "
          >
            <el-button class="action-button">
              <span class="el-icon-copy-document" />
              <span>Copy</span>
            </el-button>
            <el-button class="action-button" @click="handleDeleteItems">
              <span class="el-icon-delete" />
              <span>Delete</span>
            </el-button>
            <el-button class="action-button">
              <span class="el-icon-edit" />
              <span>Edit</span>
            </el-button>
            <el-button class="action-button">
              <span class="el-icon-video-play" />
              <span>Run</span>
            </el-button>
          </div>
          <div
            v-else-if="
              selectedTypes.length == 1 &&
                (selectedTypes.includes('folder') ||
                  selectedTypes.includes('file'))
            "
          >
            <el-button class="action-button">
              <span class="el-icon-copy-document" />
              <span>Copy</span>
            </el-button>
            <el-button class="action-button" @click="handleDeleteItems">
              <span class="el-icon-delete" />
              <span>Delete</span>
            </el-button>
            <el-button class="action-button">
              <span class="el-icon-download" />
              <span>Download</span>
            </el-button>
          </div>
          <el-button circle class="toggle-panel-button">
            <span class="el-icon-info" />
          </el-button>
        </div>
      </div>
    </div>
    <div class="project-list">
      <el-row :gutter="10">
        <el-col :xs="24" :sm="24" :md="4" :lg="4" :xl="4">
          <el-tree
            class="tree__folders"
            :data="treeFolders"
            :props="treeProps"
            default-expand-all
            highlight-current
            node-key="id"
            draggable
            :allow-drag="checkAllowDrag"
            :allow-drop="checkAllowDrop"
            @node-click="handleFolderTreeClick"
          >
            <template #default="{ node, data }">
              <div class="folder__node" @drop="handleDropFolder(data, $event)">
                <div class="folder__content">
                  <i v-if="data.class === 'folder'" class="el-icon-folder" />
                  <i v-if="data.class === 'project'" class="fas fa-briefcase" />
                  <div class="folder__text">
                    <span>{{ data.name }}</span>
                  </div>
                </div>
              </div>
            </template>
          </el-tree>
        </el-col>
        <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
          <project-detail />
        </el-col>
        <el-col :xs="24" :sm="24" :md="4" :lg="4" :xl="4">
          <project-item-left :items="selectedItems" />
        </el-col>
      </el-row>
    </div>
  </div>

  <el-dialog v-model="showAddFolder" :show-close="false" width="20%">
    <template #title>
      <div class="model-header">
        <h4>
          <em
            class="el-icon-plus"
            style="margin-right: 10px; font-size: 1rem"
          />
          New folder
        </h4>
      </div>
    </template>
    <template #default>
      <div style="display: flex; flex: auto; align-items: center">
        <span style="margin-right: 20px; min-width: 50px">Name</span>
        <el-input v-model="newFolderName" placeholder="" />
      </div>
    </template>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="showAddFolder = false">Cancel</el-button>
        <el-button type="primary" @click="handleCreateFolder">Create</el-button>
      </span>
    </template>
  </el-dialog>

  <upload
    :show-dialog="showUploadData"
    :project-id="projectId"
    :project-name="projectName"
    :upload-path="currentFolder"
    :remote-files="files"
    resumable-upload-path="https://workflow.com/files/_resumable"
    @close-upload="handleCloseUpload"
  />
</template>

<script>
import ProjectDetail from "@/components/project/ProjectDetail";
import ProjectItemLeft from "@/components/project/ProjectItemLeft";
import TitleBar from "@/components/TitleBar";
import Upload from "@/components/upload/Upload";

import { useStore } from "vuex";
import { useRoute, onBeforeRouteUpdate } from "vue-router";
import { ref, computed, onMounted, watch } from "vue";

export default {
  name: "ProjectView",
  components: { ProjectDetail, ProjectItemLeft, TitleBar, Upload },
  setup() {
    const route = useRoute();
    const store = useStore();
    const treeProps = { children: "children", label: "name" };
    const projectId = ref();
    projectId.value = route.params.projectId;

    const navPaths = ref([]);
    navPaths.value.push({
      name: "All Projects",
      class: "project",
      to: { path: "/projects" },
    });

    const getCurrentFolder = (path) => {
      let currentFolder = "";
      if (Array.isArray(path)) {
        currentFolder = "/" + path.join("/");
      } else if (path) {
        currentFolder = "/" + path;
      } else currentFolder = "/";
      return currentFolder;
    };

    const calNavPaths = (path) => {
      if (!path) path = "";

      if (Array.isArray(path)) {
        if (path.length == 1) {
          path = path[0];
        }
      }

      const paths = path.split("/");
      if (paths.length > 1) {
        path = paths;
      } else path = paths[0];

      let folderPaths = [];
      if (Array.isArray(path)) {
        folderPaths = path.map((item, index) => {
          return {
            name: item,
            to: {
              name: "Project Detail",
              class: "folder",
              params: {
                projectId: projectId.value,
                path: path.slice(0, index + 1),
              },
            },
          };
        });
      } else if (path) {
        folderPaths = [
          {
            name: path,
            class: "folder",
            to: {
              name: "Project Detail",
              params: { projectId: projectId.value, path: path },
            },
          },
        ];
      }

      if (store.state.project.project) {
        const projectPath = {
          name: store.state.project.project.name,
          class: "project",
          to: {
            name: "Project Detail",
            params: {
              projectId: projectId.value,
            },
          },
        };

        navPaths.value.splice(
          1,
          navPaths.value.length - 1,
          projectPath,
          ...folderPaths
        );
      } else
        navPaths.value.splice(1, navPaths.value.length - 1, ...folderPaths);
    };

    onBeforeRouteUpdate(async (to) => {
      const currentFolder = getCurrentFolder(to.params.path);

      store.dispatch("project/ChangeCurrentFolder", currentFolder);
      store.dispatch("project/GetProjectDataFiles", {
        projectId: projectId.value,
        currentFolder: currentFolder,
      });

      calNavPaths(to.params.path);
    });

    onMounted(() => {
      const currentFolder = getCurrentFolder(route.params.path);

      store.dispatch("project/ChangeCurrentFolder", currentFolder);
      store.dispatch("project/GetProjectDataFiles", {
        projectId: projectId.value,
        currentFolder: currentFolder,
      });
      store.dispatch("project/GetProject", projectId.value);
      store.dispatch("project/GetProjectWorkflows", projectId.value);

      calNavPaths(route.params.path);
    });

    watch(
      () => store.state.project.project,
      () => {
        calNavPaths(route.params.path);
      }
    );

    return {
      navPaths,
      treeProps,
      showAddFolder: ref(false),
      showUploadData: ref(false),
      newFolderName: ref(""),
      projectName: computed(() => store.state.project.project?.name),
      projectId,
      currentFolder: computed(() => store.state.project.currentFolder),
      project: computed(() => store.state.project.project),
      files: computed(() => store.state.project.files),
      treeFolders: computed(() => store.getters["project/treeFolders"]),
      selectedItems: computed(() => store.state.project.selectedItems),
      getProjectItem: (id) => store.getters["project/getProjectItem"](id),
      AddProjectFolder: (folder) =>
        store.dispatch("project/AddProjectFolder", folder),
      DeleteProjectFolder: (folder) =>
        store.dispatch("project/DeleteProjectFolder", folder),
      UpdateProjectPath: (source, target) =>
        store.dispatch("project/UpdateProjectPath", {
          source: source,
          target: target,
        }),
      DeleteProjectItems: (items) =>
        store.dispatch("project/DeleteProjectItems", items),
    };
  },
  computed: {
    selectedTypes() {
      const types = this.selectedItems.map((item) => item.class);
      return [...new Set(types)];
    },
  },
  methods: {
    handleNewFolder() {
      this.showUploadData = false;
      this.showAddFolder = true;
    },
    handleUploadData() {
      this.showUploadData = true;
    },
    handleNewWorkflow() {
      this.$router.push({
        name: "Workflow Edit",
        params: { projectId: this.projectId },
      });
    },
    handleCloseUpload() {
      this.showUploadData = false;
    },
    async handleCreateFolder() {
      await this.AddProjectFolder({
        projectId: this.project.id,
        folder: this.newFolderName,
        path: this.currentFolder + this.newFolderName,
      });
      this.newFolderName = "";
      this.showAddFolder = false;
    },
    async handleDeleteFolder() {
      await this.DeleteProjectFolder({
        projectId: this.project.id,
        folderId: this.selectedItems[0].id,
      });
    },
    async handleFolderTreeClick(event) {
      this.$router.push({
        name: "Project Detail",
        params: { path: event.path.slice(1) },
      });
    },
    async handleDropFolder(data, event) {
      const fileId = event.dataTransfer.getData("text/plain");
      console.log("Drop folder from ", fileId, " to ", data);
      const source = this.getProjectItem(fileId);
      if (data.path === "/") {
        await this.UpdateProjectPath(source, { path: "/", class: "folder" });
      } else {
        const target = this.getProjectItem(folderId);
        await this.UpdateProjectPath(source, target);
      }
    },
    checkAllowDrag(node) {
      return false;
    },
    checkAllowDrop(draggingNode, dropNode, type) {
      return true;
    },
    handleDeleteItems() {
      this.DeleteProjectItems(this.selectedItems);
    },
  },
};
</script>

<style lang="scss" scoped>
.padding__items {
  margin: 0px 3px;
}
.main-table {
  &:after {
    background: #eeeeee;
    content: "";
    width: 2px;
    height: 300px;
    display: block;
  }
}
.header-title-container {
  display: flex;
  align-items: center;

  .only-title {
    max-width: 100%;
    padding-right: 15px;
  }

  .header-title {
    font-size: 14;
    text-align: center;
    line-height: 24px;
    white-space: nowrap;
    text-overflow: ellipsis;
    padding: 0 0 0 5px;
    margin: 0;
  }
}

ul {
  margin-block-start: 1em;
  margin-block-end: 1em;
}

.el-menu {
  background-color: #ffffff;
  border-right: none;
}

.nav-menu {
  position: relative;
  margin: 0px;
  display: flex;
  padding: 0px;
  list-style: none;
  align-items: center;

  &:before {
    width: 2px;
    height: 50px;
    content: "";
    display: block;
    position: relative;
    background: #eeeeee;
    margin-left: 5px;
    margin-right: 5px;
  }

  li {
    display: flex;
    box-sizing: border-box;
    align-items: center;
    font-weight: bold;
    border-bottom: 3px solid transparent;
  }

  .link {
    cursor: pointer;
    display: block;
    overflow: hidden;
    max-width: 175px;
    text-align: center;
    font-weight: bold;
    font-size: 14px;
    line-height: 16px;
    white-space: nowrap;
    text-overflow: ellipsis;
    text-transform: uppercase;
  }
}

.menu-bar-container {
  display: flex;
  position: relative;
  width: 100%;

  .menu-bar {
    flex: auto;
    display: flex;
    padding-left: 50px;
    padding: 5px 20px 4px 20px;
    flex-wrap: wrap;
    align-items: center;
    border-bottom: 1px solid #e5e5e5;

    .menu-bar__title {
      font-size: 18px;
    }
  }

  .button-wrapper {
    flex: auto;
    display: flex;
    min-height: 38px;
    align-items: center;
    justify-content: flex-end;
  }
}

.dropdown-flyout {
  border: 2px solid;
  border-radius: 5px;
  border-color: #d5d5d5;
  box-shadow: 0px 2px 12px 0px;
  background-color: #f8f8f8;
  right: 100px;

  .new-folder-widget {
    align-items: center;
    padding: 5px;
    background: #ffffff;
  }

  .input-line {
    display: flex;
    font-size: 16px;
    align-items: center;

    .input-folder-name {
      height: 30px;
      margin-left: 20px;
      margin-right: 20px;
    }
  }
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

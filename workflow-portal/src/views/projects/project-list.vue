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
            <span>ALL</span>
          </router-link>
        </el-menu-item>
        <el-menu-item index="">
          <router-link to="/biosample" class="link">
            <span>RESOURCES</span>
          </router-link>
        </el-menu-item>
      </el-menu>
    </template>
    <template #right>
      <div class="title-bar-right">
        <el-button class="action-button" @click="showNewProject = true">
          <span class="el-icon-circle-plus" />
          <span>New Project</span>
        </el-button>
        <el-button circle class="toggle-panel-button">
          <span class="el-icon-info" />
        </el-button>
      </div>
    </template>
  </title-bar>
  <div class="project-list">
    <el-row :gutter="10">
      <el-col :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
        <projects-table />
      </el-col>
      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
        <project-left />
      </el-col>
    </el-row>
  </div>

  <el-dialog v-model="showNewProject" width="600px">
    <template #title>
      <div class="model-header">
        <h4>
          <em
            class="el-icon-plus"
            style="margin-right: 10px; font-size: 1rem"
          />
          New Project
        </h4>
      </div>
    </template>
    <template #default>
      <el-form :model="form" label-width="150px">
        <el-form-item label="Project Name" style="font-weight: bold">
          <el-input v-model="form.name" placeholder="Untitled Project" />
        </el-form-item>
        <el-divider content-position="left">
          <span style="font-weight: light">More info</span>
        </el-divider>
        <el-form-item label="Tags" style="font-weight: bold">
          <el-input v-model="form.tags" placeholder="" />
        </el-form-item>
        <el-form-item label="Project Summary" style="font-weight: bold">
          <el-input v-model="form.summary" placeholder="" />
        </el-form-item>
        <el-form-item label="Project Description" style="font-weight: bold">
          <el-input
            v-model="form.description"
            placeholder=""
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
    </template>
    <template #footer>
      <span class="el-dialog__footer">
        <el-button @click="showNewProject = false">Cancel</el-button>
        <el-button class="action-button" @click="hanldeNewProject">
          Create Project
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { mapActions } from "vuex";
import ProjectsTable from "@/components/project/ProjectsTable";
import ProjectLeft from "@/components/project/ProjectLeft";
import TitleBar from "@/components/TitleBar.vue";

export default {
  name: "ProjectList",
  components: {
    ProjectsTable,
    ProjectLeft,
    TitleBar,
  },
  data() {
    return {
      showNewProject: false,
      isShowFull: false,
      form: {
        name: null,
        description: null,
        tags: [],
        summary: null,
      },
    };
  },
  methods: {
    ...mapActions("projects", ["CreateProject"]),
    async hanldeNewProject() {
      const project = {
        name: this.form.name,
        description: this.form.description,
        summary: this.form.summary,
      };
      await this.CreateProject(project);
      await this.GetProjectList();
      this.showNewProject = false;
    },
  },
};
</script>

<style lang="scss" scoped>
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

.el-dialog__body {
  padding-bottom: 0px;
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

.action-button {
  margin: 2px 25px 2px 5px;
  color: #ffffff;
  border: 1px solid;
  outline: none;
  padding: 2px 8px;
  font-weight: 700;
  border-color: #00adae;
  border-radius: 3px;
  background-color: #00adae;
}

.toggle-panel-button {
  background-color: #e0f2fa;
}

.title-bar-right {
  margin-right: 40px;
}
</style>

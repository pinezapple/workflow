<template>
  <title-bar>
    <template #left>
      <div class="header-container">
        <div class="header-title-container only-title">
          <h1 class="header-title">
            <span>Workflows</span>
          </h1>
        </div>
      </div>
      <el-menu default-active="1" class="nav-menu">
        <el-menu-item index="">
          <router-link to="/workflow" class="link">
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
    <!-- <template #right>
      <el-button class="action-button">
        <span class="el-icon-circle-plus" />
        <span>New Analyses</span>
      </el-button>
    </template> -->
  </title-bar>
  <div class="workflow-list">
    <el-row :gutter="10">
      <el-col :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
        <workflows-table />
      </el-col>
      <el-col :xs="24" :sm="24" :md="6" :lg="6" :xl="6">
        <workflow-left />
      </el-col>
    </el-row>
  </div>
</template>

<script>
import WorkflowsTable from "@/components/workflow/WorkflowsTable";
import WorkflowLeft from "@/components/workflow/WorkflowLeft";
import TitleBar from "@/components/TitleBar";

export default {
  name: "WorkflowList",
  components: {
    WorkflowsTable,
    WorkflowLeft,
    TitleBar,
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
  padding: 4px 8px;
  font-weight: 700;
  border-color: #00adae;
  border-radius: 3px;
  background-color: #00adae;
}
</style>

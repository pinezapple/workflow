<template>
  <div class="nav-table">
    <div class="header_table-tool">
      <img
        class="header_image-tool"
        src="../../assets/header_logo.png"
        alt="logo header"
      >
      <h6>SELECTED A WORKFLOW</h6>
    </div>
    <div v-if="!workflow" class="no_item-tool">
      <img alt="Dont have items selected" src="../../assets/no_items.png">
      <h6>No Items Selected</h6>
      <p>Seleted items to show details</p>
    </div>
    <div v-else style="text-align: left; padding-top: 20px">
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Name</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ workflow.name }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Workflow ID</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ workflow.id }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Class</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ workflow.class }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Project ID</span>
        </el-col>
        <el-col>
          <router-link
            class="info-block-value"
            :to="{
              name: 'Project Detail',
              params: { projectId: workflow.project_id },
            }"
          >
            {{ workflow.project_id }}
          </router-link>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Project name</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ workflow.project_name }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Created by</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ workflow.owner }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Created</span>
        </el-col>
        <el-col>
          <span class="info-block-value">
            {{ $filters.datetime(workflow.created_at) }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Updated</span>
        </el-col>
        <el-col>
          <span class="info-block-value">
            {{ $filters.datetime(workflow.updated_at) }}</span>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { computed } from "@vue/runtime-core";
import { useStore } from "vuex";
export default {
  name: "WorkflowLeft",
  components: {},
  setup() {
    const store = useStore();

    return {
      workflow: computed(() => store.state.workflow.selectedWorkflow),
    };
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

.info-block {
  padding-bottom: 10px;
}

.info-block-label {
  padding-bottom: 5px;
  font-weight: bold;
  font-size: 14px;
  color: #505050;
}

.info-block-value {
  font-weight: normal;
  font-size: 12px;
}
</style>

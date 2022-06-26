<template>
  <div class="nav-table">
    <div class="header_table-tool">
      <img
        class="header_image-tool"
        src="../../assets/header_logo.png"
        alt="logo header"
      >
      <h6>SELECTED A PROJECT</h6>
    </div>
    <div v-if="!project" class="no_item-tool">
      <img alt="Dont have items selected" src="../../assets/no_items.png">
      <h6>No Items Selected</h6>
      <p>Seleted items to show details</p>
    </div>
    <div v-else style="text-align: left; padding-top: 20px">
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Project ID</span>
        </el-col>
        <el-col>
          <span class="info-block-value">{{ project.id }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Project name</span>
        </el-col>
        <el-col>
          <span class="info-block-value">{{ project.name }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Summary</span>
        </el-col>
        <el-col>
          <span class="info-block-value">{{ project.summary }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Description</span>
        </el-col>
        <el-col>
          <span class="info-block-value">{{ project.description }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Author</span>
        </el-col>
        <el-col>
          <span class="info-block-value">{{ project.author }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Created</span>
        </el-col>
        <el-col>
          <span class="info-block-value">{{
            $filters.datetime(project.created_at)
          }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Modified</span>
        </el-col>
        <el-col>
          <span class="info-block-value">{{
            $filters.datetime(project.updated_at)
          }}</span>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { computed } from "vue";
import { useStore } from "vuex";
export default {
  name: "LeftTable",
  setup() {
    const store = useStore();
    return {
      project: computed(() => store.state.projects.selectedProject),
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

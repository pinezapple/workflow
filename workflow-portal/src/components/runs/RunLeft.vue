<template>
  <div class="nav-table">
    <div class="header_table-tool">
      <img
        class="header_image-tool"
        src="../../assets/header_logo.png"
        alt="logo header"
      >
      <h6>SELECTED AN ANALYSE</h6>
    </div>
    <div v-if="!analyse" class="no_item-tool">
      <img alt="Dont have items selected" src="../../assets/no_items.png">
      <h6>No Items Selected</h6>
      <p>Seleted an item to show detail</p>
    </div>
    <div v-else style="text-align: left; padding-top: 20px">
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">ID</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ analyse.id }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Name</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ analyse.name }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">State</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ analyse.state }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Workflow ID</span>
        </el-col>
        <el-col>
          <span class="info-block-value">
            {{ analyse.request.workflow_url }}</span>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { computed } from "vue";
import { useStore } from "vuex";
export default {
  name: "RunLeft",
  setup() {
    const store = useStore();

    return {
      analyse: computed(() => store.state.analyses.selectedAnalyse),
    };
  },
};
</script>

<style>
.button_float-right {
  float: right;
}

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

<template>
  <div class="nav-table">
    <div v-if="items.length == 0" class="no_item-tool">
      <img alt="Dont have items selected" src="../../assets/no_items.png">
      <h6>No Items Selected</h6>
      <p>Seleted items to show details</p>
    </div>
    <div v-else-if="items.length > 1">
      <div class="header_table-tool">
        <div style="display: flex; padding: 10px; align-items: center">
          <span style="margin-right: 6px">
            <em
              class="el-icon-copy-document"
              style="font-size: 2em; margin-right: 5px"
            />
          </span>
          <div
            style="
              display: flex;
              overflow: hidden;
              flex-direction: column;
              align-items: start;
            "
          >
            <div
              style="
                font-size: 18px;
                line-height: 22px;
                vertical-align: top;
                white-space: nowrap;
              "
            >
              <span>{{ items.length }} items selected</span>
            </div>
            <div
              style="
                font-size: 12px;
                line-height: 18px;
                vertical-align: top;
                white-space: nowrap;
              "
            >
              Multiple items
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else style="text-align: left">
      <div class="header_table-tool">
        <div style="display: flex; padding: 10px; align-items: center">
          <span style="margin-right: 6px">
            <em
              class="el-icon-copy-document"
              style="font-size: 2em; margin-right: 5px"
            />
          </span>
          <div
            style="
              display: flex;
              overflow: hidden;
              flex-direction: column;
              align-items: start;
            "
          >
            <div
              style="
                font-size: 18px;
                line-height: 22px;
                vertical-align: top;
                white-space: nowrap;
              "
            >
              <span>{{ items[0].name }}</span>
            </div>
            <div
              style="
                font-size: 12px;
                line-height: 18px;
                vertical-align: top;
                white-space: nowrap;
              "
            >
              {{ items[0].class }}
            </div>
          </div>
        </div>
      </div>
      <el-row justify="left" class="info-block" style="padding-top: 10px">
        <el-col>
          <span class="info-block-label">Name</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ items[0].name }}</span>
        </el-col>
      </el-row>
      <el-row
        v-if="items[0].class == 'workflow' || items[0].class == 'tool'"
        justify="left"
        class="info-block"
      >
        <el-col>
          <span class="info-block-label">Workflow ID</span>
        </el-col>
        <el-col>
          <router-link
            class="info-block-value"
            :to="{
              name: 'Workflow Edit',
              params: {
                workflowId: items[0].id,
                projectId: items[0].project_id,
              },
            }"
          >
            {{ items[0].id }}
          </router-link>
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
              params: { id: items[0].project_id },
            }"
          >
            {{ items[0].project_id }}
          </router-link>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Project name</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ items[0].project_name }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Path</span>
        </el-col>
        <el-col>
          <span class="info-block-value">
            <span v-if="items[0].class === 'file'" style="margin-right: 3px">
              {{ items[0].project_path }}
            </span>
            <span v-else style="margin-right: 3px">
              {{ items[0].path }}
            </span>
          </span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Class</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ items[0].class }}</span>
        </el-col>
      </el-row>
      <el-row v-if="items[0].class == 'file'" justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">File size</span>
        </el-col>
        <el-col>
          <span class="info-block-value">
            {{ $filters.fileSize(items[0].size) }}
          </span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Created by</span>
        </el-col>
        <el-col>
          <span class="info-block-value"> {{ items[0].author }}</span>
        </el-col>
      </el-row>
      <el-row justify="left" class="info-block">
        <el-col>
          <span class="info-block-label">Created</span>
        </el-col>
        <el-col>
          <span class="info-block-value">
            {{ $filters.datetime(items[0].created_at) }}</span>
        </el-col>
      </el-row>
      <el-row
        v-if="
          items[0].class == 'workflow' ||
            items[0].class == 'tool' ||
            items[0].class == 'folder'
        "
        justify="left"
        class="info-block"
      >
        <el-col>
          <span class="info-block-label">Updated</span>
        </el-col>
        <el-col>
          <span class="info-block-value">
            {{ $filters.datetime(items[0].updated_at) }}</span>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script>
import { computed } from "vue";
import { useStore } from "vuex";
export default {
  name: "ProjectItemLeft",
  setup() {
    const store = useStore();

    return {
      items: computed(() => store.state.project.selectedItems),
    };
  },
};
</script>

<style lang="scss" scoped>
.header_table-tool {
  padding: 5px;
  display: flex;
  flex-direction: row;
  align-items: center;
  flex: auto;
  border-bottom: 1px solid var(--violet-white-color);
  border-top: 1px solid var(--violet-white-color);
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

<template>
  <div class="tool">
    <div class="tool__title">
      <div style="font-size: 1.5rem" @click="addEditorTool(tool)">
        <em class="fas fa-plus-circle tool__icon" />
      </div>
      <div class="tool__name">
        <span>{{ tool.name }} </span>
        <em class="el-icon-info tool__icon-info" />
      </div>
      <div class="tool__versions">
        <el-dropdown>
          <span class="el-dropdown-link">
            Versions <i class="el-icon-arrow-down el-icon--right" />
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-for="version in versions" :key="version.id">
                <span>{{ version.semver }}</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div class="tool__category">
        <span>Category: {{ tool.category }}</span>
      </div>
    </div>
    <div class="tool__inout">
      <ul class="tool__input">
        <li v-for="input in inputs" :key="input.id" class="tool__inout-item">
          <em class="el-icon-document tool__icon" />
          <span class="tool__inout-item-label">{{ input.extension }}</span>
        </li>
      </ul>
      <span v-if="inputs.length" class="tool__inout-separator">
        <i class="fas fa-chevron-right" />
      </span>

      <ul class="tool__output">
        <li v-for="output in outputs" :key="output.id" class="tool__inout-item">
          <em class="el-icon-document tool__icon" />
          <span class="tool__inout-item-label">{{ output.extension }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
import { toRefs, ref } from "vue";
import { mapMutations } from "vuex";

export default {
  name: "ToolSmall",
  props: {
    tool: {
      type: Object,
      required: true,
    },
    showTooltip: {
      type: Boolean,
      default: true,
    },
  },
  setup(props) {
    const { tool } = toRefs(props);
    const versions = ref([]);
    versions.value = tool.value.versions?.map(
      (version) => ({
        id: version.id,
        semver: version.semver,
      }),
      []
    );
    const inputs = ref([]);
    inputs.value = tool.value.versions?.[0]?.inputs?.filter(
      (input) => input.extension !== null && input.extension !== undefined
    );
    if (!inputs.value) inputs.value = [];

    const outputs = ref([]);
    outputs.value = tool.value.versions?.[0]?.outputs;
    if (!outputs.value) outputs.value = [];

    return {
      versions,
      inputs,
      outputs,
    };
  },
  methods: {
    ...mapMutations("editor", ["addEditorTool"]),
  },
};
</script>

<style lang="scss" scoped>
.tool {
  margin: 5px;
  border-bottom: 1px solid #e9e9e9;
  margin-top: 12px;

  .tool__title {
    display: flex;
    align-items: center;

    & > * {
      margin-right: 15px;
    }

    .tool__name {
      font-size: 16px;
      font-weight: bold;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;

      .tool__icon {
        margin-right: 10px;
      }

      .tool__icon-info {
        margin-left: 10px;
      }
    }
    .tool__versions {
      font-weight: normal;
      font-size: 10px;
    }

    .tool__category {
      font-weight: normal;
      font-size: 12px;
    }
  }

  .tool__inout {
    display: flex;
    align-items: center;
    justify-content: flex-end;

    .tool__input,
    .tool__output {
      list-style-type: circle;
      display: flex;
      align-items: center;
      padding: 0px 5px;

      li {
        list-style: none;
        background-color: #f4f4f6;
        border: 1px solid #e5e5e5;
        border-radius: 20px;
        padding: 9px 12px;
        font-size: 9px;
        margin-right: 2px;
        letter-spacing: 0.5px;
        text-transform: uppercase;
      }
    }

    .tool__inout-separator {
      font-size: 16px;
      padding: 0 8px;
      text-align: center;
      vertical-align: middle;
    }
  }
}
</style>

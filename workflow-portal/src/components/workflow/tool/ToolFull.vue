<template>
  <div class="tool">
    <div class="tool__actions">
      <el-dropdown @command="handleToolCommand">
        <span class="el-dropdown-link">
          Actions
          <i class="el-icon-arrow-down el-icon--right" />
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <!-- <el-dropdown-item command="moreinfo">
              <i class="fas fa-info" />
              <span style="margin-left: 10px">More info</span>
            </el-dropdown-item> -->
            <el-dropdown-item command="remove">
              <i class="fas fa-times" />
              <span style="margin-left: 10px">Remove</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <el-row style="align-items: center">
      <el-col :md="9">
        <ul class="tool__inputs">
          <li v-for="input in inputs" :key="input.id" class="inputs__item">
            <div
              :class="[
                hoverObject &&
                  hoverObject.type === 'input' &&
                  hoverObject.id === input.id
                  ? 'hover-item'
                  : '',
                hasLink(input.id) ? 'input__main--link' : 'input__main--nolink',
                highlight('input', input.id) ? 'input__main--hover' : '',
                'input__main',
              ]"
              draggable="true"
              @mouseover="handleMouseOverInput($event, input)"
              @mouseleave="handleMouseLeave($event, input)"
              @drop="handleDrop($event, input)"
              @dragenter="handleDragEnter($event)"
              @dragover="handleDragOver($event)"
            >
              <div v-if="hasLink(input.id)" class="input__info">
                <div class="input__extension">
                  <i class="far fa-file" />
                  <span> via</span>
                </div>
                <div class="input__link">
                  <span>{{ linkName(input.id) }}</span>
                </div>
                <div
                  class="input__name input__name--maxwidth"
                  @click="handleInputClicked(input)"
                >
                  <span>{{ input.name }}</span>
                </div>
                <div
                  class="input__link--delete"
                  @click="handleRemoveLink(tool.tool_index, input.id)"
                >
                  <i class="fas fa-times-circle" />
                </div>
                <span class="input__divider" />
                <div class="input__edit">
                  <i class="fas fa-pencil-alt" />
                </div>
              </div>
              <div v-else-if="hasInputValue(input.id)" class="input__info">
                <div class="input__extension">
                  <i class="far fa-file" />
                </div>
                <div class="input__link">
                  <span>{{ inputName(input.id) }}</span>
                </div>
                <div
                  class="input__name input__name--maxwidth"
                  @click="handleInputClicked(input)"
                >
                  <span>{{ input.name }}</span>
                </div>
                <div
                  class="input__link--delete"
                  @click="handleRemoveLink(tool.tool_index, input.id)"
                >
                  <i class="fas fa-times-circle" />
                </div>
                <span class="input__divider" />
                <div class="input__edit">
                  <i class="fas fa-pencil-alt" />
                </div>
              </div>
              <div v-else class="input__info">
                <div class="input__extension">
                  <i class="far fa-file" />
                  <span>{{ input.extension }}</span>
                </div>
                <div class="input__name" @click="handleInputClicked(input)">
                  <span>{{ input.name }}</span>
                </div>
                <span class="input__divider" />
                <div class="input__edit">
                  <i class="fas fa-pencil-alt" />
                </div>
              </div>
            </div>
          </li>
        </ul>
      </el-col>
      <el-col :md="6">
        <div class="tool__main">
          <i class="fas fa-chevron-right tool__arrow" />
          <div class="tool__header el-col-20">
            <div class="tool__name">
              <span>{{ tool.name }}</span>
            </div>
            <div class="tool__settings">
              <i class="fas fa-cog" />
            </div>
          </div>
          <i class="fas fa-chevron-right tool__arrow" />
        </div>
      </el-col>
      <el-col :md="9">
        <ul class="tool__outputs">
          <li v-for="output in outputs" :key="output.id" class="outputs__item">
            <div
              class="output__main"
              :class="[
                hoverObject &&
                  hoverObject.type === 'output' &&
                  hoverObject.id === output.id
                  ? 'hover-item'
                  : '',
                highlight('output', output.id) ? 'output__main--hover' : '',
              ]"
              draggable="true"
              drop
              @dragstart="handleDragStart($event, output)"
              @mouseover="handleMouseOverOutput($event, output)"
              @mouseleave="handleMouseLeave($event, output)"
            >
              <div class="output__extension">
                <i class="far fa-file" />
                <span>{{ output.extension }}</span>
              </div>
              <div class="output__name">
                <span>{{ output.name }}</span>
              </div>
              <div class="output__divider" />
              <div class="output__edit">
                <i class="fas fa-pencil-alt" />
              </div>
            </div>
          </li>
        </ul>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { toRefs, ref } from "vue";
import { mapMutations } from "vuex";

export default {
  name: "ToolFull",
  props: {
    tool: {
      type: Object,
      required: true,
    },
    links: {
      type: Array,
      required: false,
      default() {
        return [];
      },
    },
    highlightLinks: {
      type: Array,
      required: false,
      default() {
        return [];
      },
    },
    inputValues: {
      type: Array,
      required: false,
      default() {
        return [];
      },
    },
    showTooltip: {
      type: Boolean,
      default: true,
    },
  },
  emits: ["setLinkInout", "removeLinkInout", "hoverLink", "removeEditorTool"],
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

    const selectedVersion = ref({});
    selectedVersion.value = tool.value.versions?.[0];

    const inputs = ref([]);
    inputs.value = selectedVersion.value?.inputs.filter(
      (input) =>
        (input.extension !== null && input.extension !== undefined) ||
        input.type === "Directory"
    );
    if (!inputs.value) inputs.value = [];

    const outputs = ref([]);
    outputs.value = selectedVersion.value?.outputs;
    if (!outputs.value) outputs.value = [];

    const hoverObject = ref();
    return {
      versions,
      selectedVersion,
      inputs,
      outputs,
      hoverObject,
    };
  },
  methods: {
    ...mapMutations("editor", ["showSelectableData"]),
    hasLink(input_id) {
      const link = this.links.find((l) => l.to_input.id === input_id);
      return link ? true : false;
    },
    linkName(input_id) {
      const link = this.links.find((l) => l.to_input.id === input_id);
      return link.from_output.name;
    },
    hasInputValue(input_id) {
      const inputValue = this.inputValues.find((i) => i.input.id === input_id);
      return inputValue ? true : false;
    },
    inputName(input_id) {
      const inputValue = this.inputValues.find((i) => i.input.id === input_id);
      return inputValue.selectFile.name;
    },
    highlight(iotype, id) {
      if (iotype === "input") {
        const link = this.highlightLinks.find(
          (l) => l.to_input.id === id && l.to_tool.id === this.tool.id
        );
        return link ? true : false;
      } else if (iotype === "output") {
        const link = this.highlightLinks.find(
          (l) => l.from_output.id === id && l.from_tool.id === this.tool.id
        );
        return link ? true : false;
      }

      return false;
    },
    handleDragStart(event, output) {
      event.dataTransfer.dropEffect = "copy";
      event.dataTransfer.effectAllowed = "copy";
      event.dataTransfer.setData("text/plain-tool_index", this.tool.tool_index);
      event.dataTransfer.setData("text/plain-output_id", output.id);
    },
    handleDrop(event, input) {
      const tool_index = parseInt(
        event.dataTransfer.getData("text/plain-tool_index")
      );
      const output_id = parseInt(
        event.dataTransfer.getData("text/plain-output_id")
      );

      this.$emit("setLinkInout", {
        source: { tool_index: tool_index, output_id: output_id },
        target: { tool_index: this.tool.tool_index, input_id: input.id },
      });
    },
    handleDragEnter(event) {
      event.dataTransfer.dropEffect = "copy";
      event.preventDefault();
    },
    handleDragOver(event) {
      event.dataTransfer.dropEffect = "copy";
      event.preventDefault();
    },
    handleMouseOverInput(event, input) {
      this.hoverObject = {
        tool_id: this.tool.id,
        type: "input",
        id: input.id,
      };
      this.$emit("hoverLink", { iotype: "input", id: input.id });
    },
    handleMouseOverOutput(event, output) {
      this.hoverObject = {
        tool_id: this.tool.id,
        type: "output",
        id: output.id,
      };
      this.$emit("hoverLink", { iotype: "output", id: output.id });
    },
    handleMouseLeave(event, inout) {
      this.hoverObject = null;
      this.$emit("hoverLink", { iotype: "", id: inout.id });
    },
    handleRemoveLink(tool_index, input_id) {
      console.log("remove tool link ", tool_index, input_id);
      this.$emit("removeLinkInout", {
        tool_index: tool_index,
        input_id: input_id,
      });
    },
    handleToolCommand(command) {
      if (command == "remove") {
        this.$emit("removeEditorTool", this.tool.tool_index);
      }
    },
    handleInputClicked(input) {
      this.showSelectableData({ tool: this.tool, input: input });
    },
  },
};
</script>

<style lang="scss" scoped>
.hover-item {
  background-color: #f89406;
}

.tool {
  box-sizing: border-box;

  .tool__inputs {
    padding: 0px;
    margin: 0px 5px;
    display: flex;
    flex-wrap: wrap;
    justify-content: flex-end;
    list-style-type: none;

    .inputs__item {
      border: 1px dashed #9dcee0;
      border-radius: 4px;
      overflow: hidden;
      text-align: left;
      margin: 5px;
      padding: 3px;

      .input__main--nolink {
        box-shadow: inset 0 2px 3px rgb(0 0 0 / 30%);
      }
      .input__main--link {
        background-color: #e6e6e6;
      }
      .input__main--hover {
        background-color: #5e60e6;
      }

      .input__main {
        display: flex;
        align-items: center;
        flex-wrap: nowrap;
        border: 1px solid #f89406;
        border-radius: 4px;

        .input__info {
          display: flex;
          align-items: center;

          & > * {
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
          }

          .input__extension {
            margin: 5px;
            font-size: 0.8rem;
            opacity: 0.6;
          }
          .input__link {
            background-color: #323536;
            border-radius: 4px;
            color: #fff;
            margin: 2px 0px;
            padding: 3px;
            max-width: 150px;
          }
          .input__name {
            border: 1px solid;
            border-radius: 4px;
            color: #fff;
            background-color: rgb(31, 119, 180);
            margin: 5px 0px;
            padding: 3px;
          }
          .input__name--maxwidth {
            max-width: 150px;
          }
          .input__divider {
            content: "";
            width: 1px;
            height: 20px;
            display: inline-block;
            background-color: #a7a9aa;
          }
          .input__edit {
            margin: 5px;
          }
          .input__link--delete {
            margin: 5px;
          }
        }
      }
    }
  }

  .tool__main {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .tool__header {
      display: inline-block;
      justify-content: space-between;
      margin: 20px;
      height: 80px;
      background-color: #5b919b;
      border-radius: 5px;

      .tool__name {
        display: inline-block;
        line-height: 80px;
        max-width: 70%;
        color: #fafafa;
        text-align: center;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
      }

      .tool__settings {
        float: right;
        position: relative;
        box-sizing: content-box;
        border-left: 1px solid #9dcee0;
        margin: 0 10px;
        width: 30px;
        font-size: 1.5rem;
        line-height: 80px;
      }
    }
  }

  .tool__outputs {
    padding: 0px;
    margin: 0px 15px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    justify-content: flex-start;
    list-style-type: none;

    .outputs__item {
      border: 1px dashed #9dcee0;
      border-radius: 4px;
      overflow: hidden;
      text-align: left;
      margin: 5px;
      padding: 3px;

      .output__main--hover {
        background-color: #5e60e6;
      }

      .output__main {
        display: flex;
        align-items: center;
        flex-wrap: nowrap;
        border: 1px solid #f89406;
        border-radius: 4px;

        .output__extension {
          margin: 5px;
          text-overflow: ellipsis;
          font-size: 0.8rem;
          opacity: 0.6;
        }

        .output__name {
          text-overflow: ellipsis;
          border: 1px solid;
          border-radius: 4px;
          color: #fff;
          background-color: rgb(31, 119, 180);
          margin: 5px;
          padding: 3px;
        }

        .output__divider {
          width: 1px;
          background-color: #a7a9aa;
          height: 20px;
        }

        .output__edit {
          margin: 5px;
        }
      }
    }
  }

  .tool__actions {
    float: right;
    border: solid 1px #9dcee0;
    border-radius: 4px;
    padding: 5px;

    .command__label {
      margin-left: 10px;
    }
  }
}
</style>

<template>
  <div class="editor">
    <el-row class="editor__header">
      <el-col :span="9">
        <span>Inputs</span>
      </el-col>
      <el-col :span="6">
        <div class="editor__header--app">
          <i class="fas fa-chevron-right editor__header--arrow" />
          <span>App </span>
          <i class="fas fa-chevron-right editor__header--arrow" />
        </div>
      </el-col>
      <el-col :span="9">
        Outputs
      </el-col>
    </el-row>
    <tool-full
      v-for="(tool, index) in tools"
      :key="tool.id"
      :tool="tool"
      :links="
        toolLinks.filter(
          (link) => link.to_tool.id === tool.id && link.to_tool_index === index
        )
      "
      :highlight-links="
        highlightLinks.filter(
          (link) =>
            (link.to_tool.id === tool.id &&
              hoverOnType === 'output' &&
              link.to_tool_index === index) ||
            (link.from_tool.id === tool.id &&
              hoverOnType === 'input' &&
              link.from_tool_index === index)
        )
      "
      :input-values="
        workflowInputs.filter(
          (workflowInput) => workflowInput.tool.tool_index === index
        )
      "
      class="editor__tool"
      @setLinkInout="handleSetupLinkInout"
      @removeLinkInout="handleRemoveLinkInout"
      @hoverLink="handleHoverLink"
      @removeEditorTool="handleRemoveEditorTool"
    />
  </div>

  <el-dialog :model-value="showSelectableData" width="50%" :show-close="false">
    <template #title>
      <div class="dialog__title">
        <i class="fas fa-briefcase dialog__icon" />
        <div>
          <div>Select data for reads input</div>
          <div class="dialog__subtext">
            BWA-MEM FASTQ Read Mapper
          </div>
        </div>
      </div>
    </template>
    <project-selector
      v-if="showProjects"
      @select-project="handleSelectProject"
    />
    <file-selector
      v-else
      :project-name="fileProjectName"
      @change-folder="handleChangeFolder"
      @select-file="handleSelectFile"
    />
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCloseSelectData">Cancel</el-button>
        <el-button type="primary">Confirm</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { ref, computed } from "vue";
import { mapMutations, useStore } from "vuex";
import ToolFull from "./ToolFull";
import FileSelector from "./FileSelector";
import ProjectSelector from "./ProjectSelector";
import { cloneDeep } from "lodash";

export default {
  name: "WorkflowEditor",
  components: { ToolFull, FileSelector, ProjectSelector },
  props: {
    tools: {
      type: Array,
      required: true,
    },
  },
  setup() {
    const store = useStore();

    const highlightLinks = ref([]);
    const hoverOnType = ref("");
    const showProjects = ref(true);
    const fileProjectName = ref("");

    return {
      showSelectableData: computed(() => store.state.editor.showSelectData),
      selectedInput: computed(() => store.state.editor.selectedInput),
      selectedTool: computed(() => store.state.editor.selectedTool),
      toolLinks: computed(() => store.state.editor.toolLinks),
      workflowInputs: computed(() => store.state.editor.workflowInputs),
      setCurrentFolder: (folder) =>
        store.dispatch("project/ChangeCurrentFolder", folder),
      GetProjectDataFiles: (projectId, currentFolder) =>
        store.dispatch("project/GetProjectDataFiles", {
          projectId: projectId,
          currentFolder: currentFolder,
        }),
      GetProject: (projectId) =>
        store.dispatch("project/GetProject", projectId),
      showProjects,
      fileProjectName,
      highlightLinks,
      hoverOnType,
    };
  },
  methods: {
    ...mapMutations("editor", [
      "removeEditorTool",
      "addToolLink",
      "removeToolLink",
      "hideSelectableData",
      "addWorkflowInput",
      "removeWorkflowInput",
    ]),
    handleSetupLinkInout({ source, target }) {
      const from_tool = this.tools.find(
        (tool) => tool.tool_index === source.tool_index
      );
      const from_version = from_tool.selected_version;
      const from_output = from_version.outputs.find(
        (output) => output.id === source.output_id
      );

      const to_tool = this.tools.find(
        (tool) => tool.tool_index === target.tool_index
      );
      const to_version = to_tool.selected_version;
      const to_input = to_version.inputs.find(
        (input) => input.id === target.input_id
      );

      this.addToolLink({
        from_tool_index: source.tool_index,
        from_tool: from_tool,
        from_version: from_version,
        from_output: from_output,
        to_tool_index: target.tool_index,
        to_tool: to_tool,
        to_version: to_version,
        to_input: to_input,
      });
    },
    handleRemoveLinkInout({ tool_index, input_id }) {
      console.log("editor remove link ", tool_index, input_id);
      this.removeToolLink({ tool_index: tool_index, input_id: input_id });
      this.removeWorkflowInput({ tool_index: tool_index, input_id: input_id });
    },
    handleHoverLink({ iotype, id }) {
      if (iotype === "input") {
        this.hoverOnType = "input";
        this.highlightLinks = this.toolLinks.filter(
          (link) => link.to_input.id === id
        );
      } else if (iotype === "output") {
        this.hoverOnType = "output";
        this.highlightLinks = this.toolLinks.filter(
          (link) => link.from_output.id === id
        );
      } else this.hoverOnType = null;
    },
    handleRemoveEditorTool(tool_index) {
      this.highlightLinks = this.highlightLinks.filter(
        (link) =>
          link.from_tool.tool_index !== tool_index ||
          link.to_tool.tool_index !== tool_index
      );

      this.removeEditorTool(tool_index);
    },
    handleSelectProject({ projectId, projectName }) {
      this.fileProjectName = projectName;
      this.setCurrentFolder("/");
      this.GetProject(projectId);
      this.GetProjectDataFiles(projectId, "/");
      this.showProjects = false;
    },
    handleChangeFolder({ projectId, currentFolder }) {
      console.log("Change folder: ", currentFolder);
      this.setCurrentFolder(currentFolder);
      this.GetProjectDataFiles(projectId, currentFolder);
      this.showProjects = false;
    },
    handleCloseSelectData() {
      this.showProjects = true;
      this.hideSelectableData();
    },
    handleSelectFile(selectFile) {
      const workflowInput = {
        tool: this.selectedTool,
        input: this.selectedInput,
        selectFile: cloneDeep(selectFile),
      };
      this.addWorkflowInput(workflowInput);
      this.showProjects = true;
      this.hideSelectableData();
    },
  },
};
</script>

<style lang="scss" scoped>
.editor {
  color: #323336;

  .editor__header {
    border-top: solid 1px #e0e5e6;
    border-bottom: solid 1px #e0e5e6;
    border-radius: 5px;
    line-height: 2.5;
    align-items: center;

    .editor__header--arrow {
      position: relative;
    }

    .editor__header--app {
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
  }

  .editor__tool {
    border: 1px solid #e0e5e6;
    border-radius: 10px;
    margin: 20px 5px;
    padding: 0px;
  }
}

.dialog__title {
  display: flex;
  text-transform: uppercase;
  align-items: center;
  margin: 15px 10px;
  font-size: 22px;
  font-weight: bold;
  padding-bottom: 10px;
  border-bottom: 2px solid #e5e5e5;

  .dialog__icon {
    font-size: 20px;
    margin-right: 10px;
  }

  .dialog__subtext {
    font-size: 14px;
    text-transform: initial;
    font-weight: 300;
    line-height: 16px;
    color: #a7a9aa;
    float: left;
  }
}

.select__path {
  font-weight: 400;
  font-size: 18px;
  text-transform: capitalize;
}
</style>

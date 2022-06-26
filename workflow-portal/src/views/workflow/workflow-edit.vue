<template>
  <div>
    <div class="warp__titleEditAndAdd-box">
      <el-row>
        <el-col :span="8">
          <div v-if="editWorkflowName">
            <el-input
              v-model="workflowName"
              class="editor__name"
              @change="editWorkflowName = false"
            />
          </div>
          <h3 v-else>
            {{ workflowName }}
            <div style="display: inline-block" @click="editWorkflowName = true">
              <i class="el-icon-edit project-item-icon" />
            </div>
          </h3>
        </el-col>
        <el-col :span="8" class="start__analyses--box">
          <el-button type="success" @click="showSampleNameDialog">
            Start Analyze
          </el-button>
          <el-button type="success" @click="saveWorkflow">
            Save workflow
          </el-button>
        </el-col>
      </el-row>
    </div>
    <div class="workflow-list">
      <el-row :gutter="10">
        <el-col :lg="16">
          <workflow-editor :tools="editor_tools" />
        </el-col>
        <el-col :lg="8">
          <div style="border: 1px solid #e6f6f6">
            <list-tool :tools="tools" />
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
  <el-dialog
    :model-value="displaySampleNameDialog"
    title="Sample name"
    width="30%"
    :show-close="false"
  >
    <template #title>
      <div class="model-header">
        <h4>
          <em
            class="el-icon-info"
            style="margin-right: 10px; font-size: 1rem"
          />
          Sample name
        </h4>
      </div>
    </template>
    <template #default>
      <el-form label-width="150px">
        <el-form-item label="Sample name" style="font-weight: bold">
          <el-input v-model="sampleName" />
        </el-form-item>
      </el-form>
    </template>
    <template #footer>
      <span class="el-dialog__footer">
        <el-button @click="handleCancelSampleName">Cancel</el-button>
        <el-button class="action-button" @click="startAnalyses">
          Confirm
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>
<script>
import { mapActions, mapMutations, useStore } from "vuex";
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import ListTool from "../../components/workflow/tool/ListTool";
import WorkflowEditor from "../../components/workflow/tool/WorkflowEditor";
import { generateCwl } from "@/utils/cwl";
const { stringify } = require("yaml");

export default {
  name: "WorkflowEdit",
  components: {
    ListTool,
    WorkflowEditor,
  },
  setup: function () {
    const store = useStore();
    const route = useRoute();
    const sampleName = ref("");
    const projectId = route.params.projectId;
    const workflowId = route.params.workflowId;

    console.log("project id:", projectId);
    console.log("workflow id: ", workflowId ? true : false);

    onMounted(() => {
      store.dispatch("tool/GetAllTools");

      if (workflowId) {
        store.dispatch("editor/GetWorkflow", workflowId);
      } else {
        store.commit(
          "editor/setWorkflowName",
          "Workflow-" + new Date().toISOString()
        );
      }
    });

    const editWorkflowName = ref(false);
    return {
      editWorkflowName,
      projectId,
      sampleName,
      displaySampleNameDialog: computed(
        () => store.state.editor.displaySampleNameDialog
      ),
      currentWorkflow: computed(() => store.state.editor.workflow),
      tools: computed(() => store.state.tool.tools),
      editor_tools: computed(() => store.state.editor.tools),
      toolLinks: computed(() => store.state.editor.toolLinks),
      workflowInputs: computed(() => store.state.editor.workflowInputs),
    };
  },
  computed: {
    workflowName: {
      get() {
        return this.$store.state.editor.workflow.name;
      },
      set(name) {
        this.$store.commit("editor/setWorkflowName", name);
      },
    },
  },
  methods: {
    ...mapMutations("editor", ["showSampleNameDialog", "hideSampleNameDialog"]),
    ...mapActions("editor", ["SaveWorkflow"]),
    ...mapActions("analyses", ["RunAnalyse"]),
    saveWorkflow() {
      const { cwlWorkflow, cwlSteps } = generateCwl(
        this.editor_tools,
        this.toolLinks
      );
      const steps = [];
      for (let [name, content] of cwlSteps.entries()) {
        console.log("append step: ", name);
        steps.push({ name: name, content: stringify(content) });
      }

      this.SaveWorkflow({
        projectId: this.projectId,
        content: stringify(cwlWorkflow),
        steps: steps,
      });
    },
    handleItemLeave(index) {
      for (var m = 0; m < this.columns[0].cards.length; m++) {
        for (var k = 0; k < this.columns[0].cards[m].input.length; k++) {
          this.columns[0].cards[m].input[k].isHover = false;
        }
      }
    },
    handleItemHover(item) {
      this.columns[0].cards = item;
    },
    changeWorkflowName() {
      console.log("edit workflow name");
      this.editWorkflowName = true;
    },
    handleCancelSampleName() {
      this.hideSampleNameDialog();
    },
    startAnalyses() {
      console.log("start analyses");
      this.hideSampleNameDialog();

      const workflowParams = {};
      for (const workflowInput of this.workflowInputs) {
        workflowParams[
          [
            workflowInput.tool.name,
            workflowInput.tool.tool_index,
            workflowInput.input.abbrev,
          ].join("_")
        ] = workflowInput.selectFile.path;
      }
      workflowParams["sample_name"] = this.sampleName;

      const payload = {
        workflow_url: this.currentWorkflow.id,
        workflow_params: workflowParams,
        tags: {
          runner: "Docker",
        },
      };

      console.log("Worklfow run payload: ", payload);
      this.RunAnalyse(payload);
    },
  },
};
</script>

<style scoped>
.warp__titleEditAndAdd-box {
  background: var(--violet-white-color);
  margin-bottom: 10px;
}
.editor__name {
  max-width: 300px;
}
.start__analyses--box {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-end;
}
.column_warpper--box {
  width: 100%;
}
.warp__process-box {
  display: flex;
}
.app__process--box {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-evenly;
}
</style>

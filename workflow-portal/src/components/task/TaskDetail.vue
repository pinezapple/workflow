<template>
  <gstc :items="[task]" :show-index="false" />
  <div class="task__output">
    <div class="output_button">
      <el-button type="primary" @click="showOutput = !showOutput">
        <span v-if="!showOutput">Show Outputs</span>
        <span v-else>Hide Outputs</span>
      </el-button>
    </div>
    <div v-if="showOutput" class="output__list">
      <ul>
        <li v-for="item in task?.outputs" :key="item.id">
          {{ item.name }}
          <a v-bind:href="`${item.url}`">Download</a>
        </li>
      </ul>
    </div>
  </div>
</template>
<script>
import gstc from "@/components/timeline/GSTC";
import { useStore } from "vuex";
import { computed, ref } from "vue";

export default {
  name: "TaskDetail",
  components: { gstc },
  setup() {
    const store = useStore();
    const showOutput = ref(false);

    return {
      showOutput,
      task: computed(() => store.state.run.selectedTask),
    };
  },
};
</script>

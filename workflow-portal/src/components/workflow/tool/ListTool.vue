<template>
  <tool-small v-for="tool in tools" :key="tool.id" :tool="tool" />
</template>

<script>
import ToolSmall from "./ToolSmall.vue";
export default {
  name: "ListTool",
  components: { ToolSmall },
  props: {
    tools: {
      type: Array,
      required: true,
    },
  },
  setup() {},
};
</script>

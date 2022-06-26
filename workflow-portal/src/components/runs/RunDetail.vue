<template>
  <gstc :items="tasks" />
</template>

<script>
import { computed } from "vue";
import { useStore } from "vuex";
import gstc from "@/components/timeline/GSTC";

export default {
  name: "RunDetail",
  components: { gstc },

  setup() {
    const store = useStore();

    return {
      tasks: computed(() => store.state.run.tasks),
    };
  },
};
</script>

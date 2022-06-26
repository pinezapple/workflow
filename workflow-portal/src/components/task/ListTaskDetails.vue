<template>
  <div v-if="task" class="warp__task-box">
    <div>
      <span class="text-bold">Create Time:</span> {{ task.creation_time }}
    </div>
    <div><span class="text-bold">Name:</span> {{ task.name }}</div>
    <div><span class="text-bold">State:</span> {{ task.state }}</div>
    <div v-if="task?.outputs?.length" class="margin-top">
      <span class="text-bold">Output:</span>

      <div v-for="items in task.outputs" :key="items.uuid" :items="items">
        <div>
          <span class="text-bold">description: {{ items.description }}</span>
        </div>
        <div>
          <span class="text-bold">name: {{ items.name }}</span>
        </div>
        <div>
          <span class="text-bold">path: {{ items.path }}</span>
        </div>
        <div>
          <span class="text-bold">type: {{ items.type }}</span>
        </div>
        <div>
          <span class="text-bold">url: {{ items.url }}</span>
        </div>
      </div>
    </div>
    <div v-if="task?.inputs?.length" class="margin-top">
      <span class="text-bold">Output:</span>

      <div v-for="items in task.inputs" :key="items.uuid" :items="items">
        <div>
          <span class="text-bold">description: {{ items.description }}</span>
        </div>
        <div>
          <span class="text-bold">name: {{ items.name }}</span>
        </div>
        <div>
          <span class="text-bold">path: {{ items.path }}</span>
        </div>
        <div>
          <span class="text-bold">type: {{ items.type }}</span>
        </div>
        <div>
          <span class="text-bold">url: {{ items.url }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: "ListTaskDetails",
  data() {
    return {
      id: this.$route.params.id,
    };
  },
  computed: {
    ...mapGetters("run", {
      task: "stateTask",
    }),
  },
  created: function () {
    this.getListTaskDetails(this.id);
  },
  methods: {
    ...mapActions("run", ["getListTaskDetails"]),
  },
};
</script>

<style scoped>
.warp__task-box {
  float: left;
  text-align: left;
  padding: 10px;
}
.text-bold {
  font-weight: 600;
  font-size: 16px;
}

.margin-top {
  margin-top: 20px;
}
</style>

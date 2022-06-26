<template>
  <div />
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: "LogDetails",
  data() {
    return {
      id: this.$route.params.id,
    };
  },
  computed: {
    ...mapGetters("logs", {
      log: "getLog",
    }),
  },
  created: function () {
    this.getLogDetails(this.id);
  },
  methods: {
    ...mapActions("logs", ["getLogDetails"]),
  },
};
</script>

<style scoped></style>

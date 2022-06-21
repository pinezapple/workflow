Í
<template>
  <el-row class="title_warpper">
    <el-col :span="4">
      <h3>{{ getRouteName }}</h3>
    </el-col>
  </el-row>
</template>

<script>
export default {
  name: "TitleBar",
  components: {},
  computed: {
    getRouteName() {
      return this.$route.name;
    },
  },
};
</script>

<style>
.title_warpper {
  height: 70px;
}
</style>

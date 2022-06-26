<template>
  <el-form ref="form" :model="form" label-width="120px">
    <el-form-item label="Project Name">
      <el-input v-model="form.project_name" placeholder="Project name" />
    </el-form-item>
    <el-form-item label="Description">
      <el-input
        v-model="form.description"
        type="textarea"
        :autosize="{ minRows: 4 }"
        placeholder="Description  "
      />
    </el-form-item>
  </el-form>
</template>

<script>
export default {
  name: "ProjectForm",
  data() {
    return {
      form: {
        project_name: "",
        description: "",
      },
    };
  },
};
</script>

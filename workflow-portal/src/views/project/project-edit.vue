<template>
  <div class="project-list">
    <el-row :gutter="10" justify="center" type="flex">
      <el-col :xs="24" :sm="24" :md="18" :lg="18" :xl="18">
        <project-form />
      </el-col>
    </el-row>
  </div>
</template>

<script>
import ProjectForm from "@/components/project/ProjectForm";

export default {
  name: "ProjectList",
  components: {
    ProjectForm,
  },
};
</script>

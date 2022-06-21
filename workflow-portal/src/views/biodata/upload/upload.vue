<template>
  <div class="upload">
    <el-row :gutter="10">
      <upload />
    </el-row>
  </div>
</template>

<script>
import Upload from "@/components/upload/Upload";
export default {
  name: "UploadData",
  components: {
    Upload,
  },
  setup() {},
};
</script>

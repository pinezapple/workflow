<template>
  <div>
    <el-upload
      class="upload-demo"
      drag
      action="https://jsonplaceholder.typicode.com/posts/"
      :on-preview="handlePreview"
      :on-remove="handleRemove"
      :file-list="fileList"
      multiple
    >
      <em class="el-icon-upload" />
      <div class="el-upload__text">
        Drop file here or <em>click to upload</em>
      </div>
      <div class="el-upload__tip" :v-slot="tip">
        jpg/png files with a size less than 500kb
      </div>
    </el-upload>
  </div>
</template>

<script>
export default {
  name: "Input",
  props: {
    input: {
      type: Object,
      default() {
        return {};
      },
    },
  },
};
</script>

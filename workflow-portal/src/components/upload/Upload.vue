<template>
  <div>
    <el-dialog
      v-model="showUpload"
      width="40%"
      top="10vh"
      :show-close="false"
      :before-close="cancelUpload"
    >
      <template #title>
        <h3 class="upload-title">
          <em
            class="el-icon-upload"
            style="margin-left: 15px; margin-right: 15px"
          />Upload Data
        </h3>
        <div class="content-divider" />
      </template>
      <template #default>
        <div>
          <div>
            <p style="display: flex; margin-left: 10px">
              Upload data from your computer to project:
              <span style="font-weight: bold; margin-left: 5px">{{
                projectName
              }}</span>
            </p>
            <div v-if="fileList.length">
              <el-table :data="fileList" stripe :show-header="false">
                <el-table-column
                  v-if="withSampleName"
                  label="Sample name"
                  prop="sample_name"
                  :min-width="20"
                  border
                >
                  <template #default="scope">
                    <el-input
                      v-model="scope.row.sample_name"
                      placeholder="Sample name"
                      style="text-align: right"
                    />
                  </template>
                </el-table-column>
                <el-table-column :min-width="35">
                  <template #default="scope">
                    <span>
                      <em class="el-icon-document" />
                      <span style="margin-left: 10px">{{
                        scope.row.name
                      }}</span>
                    </span>
                  </template>
                </el-table-column>
                <el-table-column label="Size" :min-width="15" align="left">
                  <template #default="scope">
                    <div>Size</div>
                    <span>{{ fileSize(scope.row.size) }}</span>
                  </template>
                </el-table-column>
                <el-table-column label="Progress" :min-width="25" align="left">
                  <template #default="scope">
                    <el-progress
                      v-if="status === 'UPLOADING'"
                      :text-inside="true"
                      :stroke-width="20"
                      :percentage="scope.row.percentage"
                    />
                    <div v-else />
                  </template>
                </el-table-column>
                <el-table-column
                  v-if="withSampleName"
                  label="Action Save"
                  :min-width="10"
                  align="right"
                >
                  <template #default="scope">
                    <el-button
                      size="small"
                      type="primary"
                      :disabled="status === 'UPLOADING' ? true : false"
                      @click="handleSaveSampleName(scope.row)"
                    >
                      <em class="el-icon-check" />
                    </el-button>
                  </template>
                </el-table-column>
                <el-table-column
                  label="Action Delete"
                  :min-width="10"
                  align="right"
                >
                  <template #default="scope">
                    <el-button
                      size="small"
                      type="primary"
                      class="action-button"
                      :disabled="status === 'UPLOADING' ? true : false"
                      @click="handleRemove(scope.row)"
                    >
                      <em class="el-icon-delete" />
                    </el-button>
                  </template>
                </el-table-column>
                <template #empty>
                  <div style="height: 0px" />
                </template>
              </el-table>
            </div>
          </div>
          <div style="margin-top: 10px">
            <el-upload
              ref="upload"
              drag
              multiple
              list-type="text"
              action="resumableUploadPath"
              :auto-upload="false"
              :on-remove="handleRemove"
              :on-change="handleChangeFiles"
              :disabled="status === 'UPLOADING' ? true : false"
            >
              <template #default>
                <div v-if="errorMsg !== ''" class="el-upload__text">
                  <em class="el-icon-upload" />
                  <p>
                    {{ errorMsg }}
                  </p>
                </div>
                <div v-else-if="status !== 'UPLOADING'" class="el-upload__text">
                  <em class="el-icon-upload" />
                  <p>
                    Drop file here to upload <i>or</i> <em>Select Files</em>
                  </p>
                  <div class="el-upload__tip">
                    Note: Uploading folders is not supported. Select individual
                    files for upload.
                  </div>
                </div>
                <div v-else class="el-upload__text">
                  <em class="el-icon-upload" />
                  <p>
                    Uploading {{ numberFiles }}
                    <span v-if="numberFiles == 1">file</span>
                    <span v-else>files</span> to the server.
                  </p>
                  <div class="el-upload__tip">
                    Note: Uploading folders is not supported. Select individual
                    files for upload.
                  </div>
                </div>
              </template>
              <template #tip>
                <div class="el-upload__tip upload__notice">
                  <p>
                    <em class="el-icon-info" style="padding-right: 5px" />
                    Do not include personally identifiable information in
                    filenames or uploaded genetic data.
                  </p>
                </div>
              </template>
              <template #file>
                <div />
              </template>
            </el-upload>
          </div>
          <el-footer style="margin-top: 10px">
            <div class="content-divider" />
            <el-row class="footer-content">
              <el-col :span="16">
                <div style="display: flex; justify-content: start">
                  <span v-if="numberFiles">
                    {{ numberFiles }} <span v-if="numberFiles == 1">file</span>
                    <span v-else>files</span> selected for upload ({{
                      fileSize(totalFileSize)
                    }})
                  </span>
                </div>
              </el-col>
              <el-col :span="8">
                <div style="display: flex; justify-content: flex-end">
                  <el-button
                    type="primary"
                    size="small"
                    class="action-button"
                    @click="cancelUpload"
                  >
                    Cancel
                  </el-button>
                  <el-button
                    class="action-button"
                    type="primary"
                    size="small"
                    :disabled="
                      (fileList.length == 0 ? true : false) || errorMsg !== ''
                    "
                    @click="submitUpload"
                  >
                    Upload file
                  </el-button>
                </div>
              </el-col>
            </el-row>
          </el-footer>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import Resumable from "resumablejs";
import { ref, toRefs } from "@vue/reactivity";
import { computed, watch } from "vue";

export default {
  name: "UploadData",
  props: {
    resumableUploadPath: {
      type: String,
      default: "http://localhost:10006/files/_resumable",
    },
    uploadPath: {
      type: String,
      default: "/",
    },
    projectName: {
      type: String,
      default: "Analysis WGS",
    },
    projectId: {
      type: String,
      default: "",
    },
    withSampleName: {
      type: Boolean,
      default: false,
    },
    showDialog: {
      type: Boolean,
      default: false,
    },
    remoteFiles: {
      type: Array,
      default: () => [],
    },
  },
  emits: ["closeUpload"],
  setup(props) {
    const { showDialog } = toRefs(props);
    const showUpload = ref(false);
    const updateShowUpload = async () => {
      showUpload.value = showDialog.value;
    };
    watch(showDialog, updateShowUpload);

    const status = ref("SELECTING");
    const updateStatus = async (newStatus) => {
      status.value = newStatus;
    };

    const fileList = ref([]);
    fileList.value = [];
    const addUploadFile = async (file) => fileList.value.push(file);
    const removeUploadFile = async (file) => {
      fileList.value = fileList.value.filter((obj) => obj.name != file.name);
    };
    const numberFiles = computed(() => fileList.value.length);
    const totalFileSize = computed(() =>
      fileList.value.reduce((acc, current) => acc + current.size, 0)
    );

    const errorMsg = ref("");
    const resumableObj = ref();

    return {
      errorMsg,
      resumableObj,
      fileList,
      addUploadFile,
      removeUploadFile,
      numberFiles,
      totalFileSize,
      status,
      updateStatus,
      showUpload,
      updateShowUpload,
    };
  },
  methods: {
    handleChangeFiles(file) {
      if (file.status == "ready") {
        for (const obj of this.fileList) {
          if (obj.name == file.name) return;
        }
        this.addUploadFile(Object.assign({ sample_name: "" }, file));
      }
      this.fileExisted();
    },
    handleRemove(file) {
      this.removeUploadFile(file);
      this.fileExisted();
    },
    handleSaveSampleName(file) {
      for (let obj of this.fileList) {
        if (obj.name == file.name) {
          obj.sample_name = file.sample_name;
        }
      }
    },
    submitUpload() {
      if (this.fileList.length == 0) return;

      this.errorMsg = "";
      this.status = "UPLOADING";
      this.resumableObj = Resumable({
        target: this.resumableUploadPath,
        chunkSize: 20 * 1024 * 1024,
        maxChunkRetries: 1,
        simultaneousUploads: 1,
        chunkRetryInterval: 1000,
        testChunks: false,
        query: { projectPath: this.uploadPath, projectID: this.projectId },
        headers: {
          Authorization:
            "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImZlbmNlX2tleV9rZXlzIn0.eyJwdXIiOiJhY2Nlc3MiLCJhdWQiOlsib3BlbmlkIiwidXNlciIsImNyZWRlbnRpYWxzIiwiZGF0YSIsImFkbWluIiwiZ29vZ2xlX2NyZWRlbnRpYWxzIiwiZ29vZ2xlX3NlcnZpY2VfYWNjb3VudCIsImdvb2dsZV9saW5rIiwiZ2E0Z2hfcGFzc3BvcnRfdjEiXSwic3ViIjoiNDEiLCJpc3MiOiJodHRwczovL2dlbm9tZS52aW5iaWdkYXRhLm9yZy91c2VyIiwiaWF0IjoxNjI2NjYyMjMzLCJleHAiOjE2MjY2OTgyMzMsImp0aSI6ImQ5Y2RhN2JkLWI0MTAtNGEyMC1hZGQ2LTA5YzRlMGY1NzE4OCIsInNjb3BlIjpbIm9wZW5pZCIsInVzZXIiLCJjcmVkZW50aWFscyIsImRhdGEiLCJhZG1pbiIsImdvb2dsZV9jcmVkZW50aWFscyIsImdvb2dsZV9zZXJ2aWNlX2FjY291bnQiLCJnb29nbGVfbGluayIsImdhNGdoX3Bhc3Nwb3J0X3YxIl0sImNvbnRleHQiOnsidXNlciI6eyJuYW1lIjoiZG9uZ29jdHVhbi4wMTAxQGdtYWlsLmNvbSIsImlzX2FkbWluIjp0cnVlLCJnb29nbGUiOnsicHJveHlfZ3JvdXAiOm51bGx9LCJwcm9qZWN0cyI6e319fSwiYXpwIjoiIn0.IfY0LKj4zvV7XJCg_qkht09S8w3mNCNs2CC_r4fAUoLXGMZKNrsEHjIaUo_bJKjHetsuswljSpHYifEiZRKk0Ec47x2StwyfFrulytADmrcvFvj2PRIw2bJicmCdwx7sVUCoQxuudMuqqjcYEDtksMmmkaffi7i84oePV7WVKf2kNwOgPbD5FyIiuzA6_VwYEiFova-rd36UCHOnAeR3HChhJuwwb2TG03XP5fWQtD5JasyqY_ex1jkE4QVYRWI7wiC1hPWserLsWBJdJCQFjejJSroZWdCTfgcSvRnlD3eR3Vr2UiLFzjlumMIdf1itJIuVTKF-fGTqZxmeFM8Oqg",
        },
      });

      this.resumableObj.on("fileAdded", function () {
        this.upload();
      });

      if (this.fileExisted()) return;

      this.resumableObj.on("error", () => {
        this.cancel();
      });

      for (const obj of this.fileList) {
        const uploadFile = obj;
        this.resumableObj.addFile(uploadFile.raw);
      }
    },
    cancelUpload() {
      if (this.resumableObj) this.resumableObj.cancel();
      this.errorMsg = "";
      this.$emit("closeUpload");
    },
    fileSize(bytes, decimals = 2) {
      if (bytes == 0) return "0 Bytes";
      var k = 1024,
        dm = decimals || 2,
        sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"],
        i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
    },
    is_blank(str) {
      return !str || str.length === 0 || str.trim() == "";
    },
    fileExisted() {
      for (const obj of this.fileList) {
        for (const remoteFile of this.remoteFiles) {
          if (obj.name == remoteFile.name) {
            this.errorMsg = "File " + remoteFile.name + " existed";
            return true;
          }
        }
      }

      this.errorMsg = "";
      return false;
    },
  },
};
</script>

<style lang="scss" scoped>
.el-dialog {
  .upload-title {
    display: flex;
    padding: 10px 10px 0px 10px;
  }
  .content-divider {
    content: "";
    width: 100%;
    height: 2px;
    background-color: #eeeeee;
  }

  .footer-content {
    padding-top: 15px;
  }
}

.upload__notice {
  border: 1px solid #ccc;
  background-color: var(--gray-white-color);
  display: flex;
  padding-left: 10px;
}
</style>

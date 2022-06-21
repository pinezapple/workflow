<template>
  <div style="width: `{tableSize}`">
    <div>
      <h4>Upload Data</h4>
      <em icon="el-icon-upload" />
    </div>
    <el-divider style="margin: 12px" />
    <el-main :style="tableStyle">
      <div>
        <p style="display: flex; margin-left: 10px">
          Upload data from your computer to project:
          <span style="font-weight: bold">{{ project_name }}</span>
        </p>
        <div v-if="fileList.length">
          <el-table
            :data="fileList"
            stripe
            :show-header="false"
            :style="tableStyle"
          >
            <el-table-column
              label="Sample name"
              prop="sample_name"
              :min-width="20"
              v-if="enable_sample_name"
            >
              <template #default="scope">
                <el-input
                  v-model="scope.row.sample_name"
                  placeholder="Sample name"
                  style="text-align: right"
                ></el-input>
              </template>
            </el-table-column>
            <el-table-column :min-width="35">
              <template #default="scope">
                <span>
                  <em class="el-icon-document" />
                  <span style="margin-left: 10px">{{ scope.row.name }}</span>
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
                ></el-progress>
                <div v-else></div>
              </template>
            </el-table-column>
            <el-table-column
              label="Action Save"
              :min-width="10"
              align="right"
              v-if="enable_sample_name"
            >
              <template #default="scope">
                <el-button
                  @click="handleSaveSampleName(scope.row)"
                  size="small"
                  type="primary"
                  :disabled="status === 'UPLOADING' ? true : false"
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
                  @click="handleRemove(scope.row)"
                  size="small"
                  type="primary"
                  :disabled="status === 'UPLOADING' ? true : false"
                >
                  <em class="el-icon-delete" />
                </el-button>
              </template>
            </el-table-column>
            <template #empty>
              <div style="height: 0px"></div>
            </template>
          </el-table>
        </div>
      </div>
      <div>
        <el-upload
          drag
          multiple
          ref="upload"
          list-type="text"
          action="resumableUploadPath"
          :auto-upload="false"
          :on-remove="handleRemove"
          :on-progress="handleProgress"
          :on-change="handleChangeFiles"
          :on-error="handleError"
          :before-upload="handleBeforeUpload"
          :disabled="status === 'UPLOADING' ? true : false"
        >
          <template #default>
            <div class="el-upload__text" v-if="status !== 'UPLOADING'">
              <em class="el-icon-upload" />
              <p>Drop file here to upload <i>or</i> <em>Select Files</em></p>
              <div class="el-upload__tip">
                Note: Uploading folders is not supported. Select individual
                files for upload.
              </div>
            </div>
            <div class="el-upload__text" v-else>
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
            <div
              class="el-upload__tip"
              style="border: 1px solid #ccc; border-radius: 10px"
            >
              <p>
                Do not include personally identifiable information in filenames
                or uploaded genetic data.
              </p>
            </div>
          </template>
          <template #file>
            <div></div>
          </template>
        </el-upload>
      </div>
      <el-footer style="margin-top: 10px">
        <el-row align="middle">
          <el-col :span="16">
            <span v-if="numberFiles">
              {{ numberFiles }} <span v-if="numberFiles == 1">file</span>
              <span v-else>files</span> selected for upload ({{
                fileSize(totalFileSize)
              }})
            </span>
          </el-col>
          <el-col :span="8">
            <div style="display: flex">
              <el-button type="primary" size="small" @click="cancelUpload"
                >Cancel</el-button
              >
              <el-button
                type="primary"
                @click="submitUpload"
                size="small"
                :disabled="fileList.length == 0 ? true : false"
                >Upload file</el-button
              >
            </div>
          </el-col>
        </el-row>
      </el-footer>
    </el-main>
  </div>
</template>

<script>
import Resumable from "resumablejs";

export default {
  name: "UploadData",
  props: {
    tableSize: {
      type: Number,
      default: 800,
    },
    uploadSize: {
      type: Number,
      default: 750,
    },
    resumableUploadPath: {
      type: String,
      default: "https://workflow.com/valkyrie/files/_resumable",
    },
    uploadPath: {
      type: String,
    },
    project_name: {
      type: String,
      default: "Analysis WGS",
    },
    enable_sample_name: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return { fileList: [], status: "SELECTING" };
  },
  computed: {
    tableStyle() {
      return { width: this.tableSize + "px", "margin-bottom": "10px" };
    },
    numberFiles() {
      return this.fileList.length;
    },
    totalFileSize() {
      return this.fileList.reduce((acc, current) => acc + current.size, 0);
    },
  },
  methods: {
    handleChangeFiles(file, fileList) {
      if (file.status == "ready") {
        for (const obj of this.fileList) {
          if (obj.name == file.name) return;
        }
        this.fileList.push(Object.assign({ sample_name: "" }, file));
      }
    },
    handleRemove(file, fileList) {
      this.fileList = this.fileList.filter((obj) => obj.name != file.name);
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
      for (const obj of this.fileList) {
        console.log("Submit object: ", obj);
        if (this.is_blank(obj.sample_name)) {
          console.log("File " + obj.name + " is missing file name");
        }
      }
      this.status = "UPLOADING";

      var r = new Resumable({
        target: this.resumableUploadPath,
        chunkSize: 20 * 1024 * 1024,
        maxChunkRetries: 3,
        simultaneousUploads: 1,
        chunkRetryInterval: 1000,
        testChunks: false,
        headers: {
          Authorization:
            "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImZlbmNlX2tleV9rZXlzIn0.eyJwdXIiOiJhY2Nlc3MiLCJhdWQiOlsib3BlbmlkIiwidXNlciIsImNyZWRlbnRpYWxzIiwiZGF0YSIsImFkbWluIiwiZ29vZ2xlX2NyZWRlbnRpYWxzIiwiZ29vZ2xlX3NlcnZpY2VfYWNjb3VudCIsImdvb2dsZV9saW5rIiwiZ2E0Z2hfcGFzc3BvcnRfdjEiXSwic3ViIjoiNDEiLCJpc3MiOiJodHRwczovL2dlbm9tZS52aW5iaWdkYXRhLm9yZy91c2VyIiwiaWF0IjoxNjIzNzU0NTUzLCJleHAiOjE2MjM3OTA1NTMsImp0aSI6IjhjODJhMWZlLWQ1NDktNDU5My05NjE4LThmM2U5NDJkN2FmYSIsInNjb3BlIjpbIm9wZW5pZCIsInVzZXIiLCJjcmVkZW50aWFscyIsImRhdGEiLCJhZG1pbiIsImdvb2dsZV9jcmVkZW50aWFscyIsImdvb2dsZV9zZXJ2aWNlX2FjY291bnQiLCJnb29nbGVfbGluayIsImdhNGdoX3Bhc3Nwb3J0X3YxIl0sImNvbnRleHQiOnsidXNlciI6eyJuYW1lIjoiZG9uZ29jdHVhbi4wMTAxQGdtYWlsLmNvbSIsImlzX2FkbWluIjp0cnVlLCJnb29nbGUiOnsicHJveHlfZ3JvdXAiOm51bGx9LCJwcm9qZWN0cyI6e319fSwiYXpwIjoiIn0.u76QjEJQkCh1cBpf4pGu89OhJNXsmKwg-ocUllTY3RmJJEj4J9PktvcAivkLKMJ2NysXdkzUN9E4Yko3oBZnBZXzZnft3WWwrbVn2_aHf53bHi6VZODU0l_vyvupZOZdUInHstRKvD2W9snDLK_cgIDjFyWZwD4tIP7ADAfFb56zC5eucmao7cN685ZrDyjOgcjT6NniCRQUoJCMNIGYM0U9slnKqG8PBboEi16Ourofvu7KOC-sem9Pe36Dj1DaLmGD3NDLfg-FWagu0biQpPj7rsMwH_2-fJ3Uh3aAMEpipuBT3N4C-ksIo-Uz9QrmJfDIawWhWMiE_4IlugJkhQ",
        },
      });

      r.on("fileAdded", function (file) {
        r.upload();
      });

      r.on("fileProgress", (file) => {
        for (let uploadFile of this.fileList) {
          if (uploadFile.name === file.fileName) {
            uploadFile.percentage = Math.round(file.progress() * 100);
          }
        }
      });

      r.on("fileError", function (file, message) {
        console.log("File error: ", file);
        console.log("File error message: ", message);
      });

      r.on("complete", function () {});

      r.on("error", () => {
        r.abort();
      });

      for (const obj of this.fileList) {
        const uploadFile = obj;
        r.addFile(uploadFile.raw);
      }
    },
    cancelUpload() {},
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
  },
};
</script>

<style lang="scss" scoped></style>

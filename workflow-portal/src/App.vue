<template>
  <el-container class="body">
    <el-header>
      <NavBar />
    </el-header>
    <el-main class="main">
      <router-view />
    </el-main>
    <el-footer>
      <div class="footer" />
    </el-footer>
  </el-container>
</template>

<script>
import { mapGetters, mapActions, mapState } from "vuex";

import NavBar from "./components/NavBar";

export default {
  name: "App",
  components: {
    NavBar,
  },
};
</script>

<style>
@import "./assets/styles/variables.css";

#app {
  font-family: Avenir, Helvetica, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
.el-header {
  padding: 0px;
  height: 100% !important;
}

.el-dialog__header {
  padding: 5px 20px 5px;
}

.el-dialog__body {
  padding-top: 0px;
  padding-bottom: 5px;
}

.el-footer {
  padding: 0px 0px;
}

.el-upload--text {
  width: 100%;
}

.el-upload-dragger {
  width: 100%;
}

.el-main {
  padding: 10px;
}

.el-table__row .el-input .el-input__inner {
  border-style: none;
}
.hover-row .el-input .el-input__inner {
  border-style: solid;
}
a {
  text-decoration: none;
}
.el-table__row .el-input .el-input__inner {
  border-style: none;
}
.hover-row .el-input .el-input__inner {
  border-style: solid;
}
.body {
  display: flex;
  flex-direction: column;
}
.main {
  min-height: calc(100vh - 120px);
}
.footer {
  margin-top: auto;
  height: 50px;
  width: 100%;
  background-color: var(--violet-dark-color);
  color: white;
  text-align: center;
  z-index: 9999;
}

.action-button {
  margin: 2px 5px 2px 5px;
  color: #ffffff;
  border: 1px solid;
  outline: none;
  min-width: 80px;
  font-weight: 700;
  border-color: #00adae;
  border-radius: 3px;
  background-color: #00adae;
}

.project-item-icon {
  margin-right: 5px;
  margin-left: 2px;
}
</style>

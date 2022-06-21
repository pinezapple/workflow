<template>
  <el-container>
    <el-header>
      <NavBar></NavBar>
    </el-header>
    <el-main>
      <div class="title_bar-box">
        <TitleBar></TitleBar>
      </div>
      <router-view></router-view>
    </el-main>
    <el-footer>Footer</el-footer>
  </el-container>
</template>

<script>
import NavBar from "./components/NavBar";
import TitleBar from "./components/TitleBar";

export default {
  name: "App",
  components: {
    NavBar,
    TitleBar,
  },
};
</script>

<style>
@import "./assets/styles/variables.css";

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
.el-header {
  padding: 0px;
  height: 100% !important;
}
.el-main {
  padding: 0px;
}
.title_bar-box {
  background-color: var(--violet-white-color);
  margin-bottom: 10px;
}

.el-table__row .el-input .el-input__inner {
  border-style: none;
}
.hover-row .el-input .el-input__inner {
  border-style: solid;
}
</style>

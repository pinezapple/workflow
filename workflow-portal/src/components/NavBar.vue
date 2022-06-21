<template>
  <div class="navbar_warpper">
    <el-row>
      <el-col :span="2" class="navbar__logo-box d-flex flex-column">
        <img
          src="../assets/dna1.png"
          alt="logo genomics"
          class="navbar__logo"
        />
      </el-col>
      <el-col :span="2">
        <h4 class="text_logo-white">
          Gen<span class="text_logo-orange">omics</span>
        </h4>
      </el-col>
      <el-col :span="20" class="navbar__menu-box">
        <el-col :span="18" :offset="6">
          <ul class="el-menu el-menu--horizontal navbar__menu">
            <router-link class="el-menu-item" to="home">Home</router-link>
            <router-link class="el-menu-item" to="projects"
              >Projects</router-link
            >
            <router-link class="el-menu-item" to="workflow"
              >Workflow</router-link
            >
            <router-link class="el-menu-item" to="analyses"
              >Analyses</router-link
            >
            <router-link class="el-menu-item" to="biosample"
              >Biosample</router-link
            >
            <router-link class="el-menu-item" to="demodata"
              >Demodata</router-link
            >
          </ul>
        </el-col>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  },
  methods: {
    handleNavDropdownCommand: function (command) {
      if (command == "logout") {
        this.$confirm("确认退出吗?", "提示", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }).then(
          function () {
            this.$message("你的退出逻辑");
          }.bind(this)
        );
      }
    },
  },
};
</script>
<style>
.text_logo-white {
  color: var(--white-color);
}
.text_logo-orange {
  color: var(--orange-color);
}
.navbar_warpper {
  background-color: var(--violet-dark-color);
}
.navbar__logo-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.navbar__logo {
  display: block;
  margin-left: auto;
  margin-right: auto;
}
.navbar__menu-box {
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.2);
}

.navbar__dropdown-menu a {
  color: inherit;
  text-decoration: none;
}
.navbar__title {
  margin-top: 0;
  margin-bottom: 0;
  margin-left: 12px;
  height: 60px;
  line-height: 60px;
  font-size: 1.5em;
  font-weight: bold;
}
.navbar__title a {
  color: var(--white-color);
  text-decoration: none;
}
.el-menu--horizontal > .el-menu-item {
  color: var(--white-color);
}
.el-menu {
  background-color: var(--violet-dark-color);
}
.el-menu.el-menu--horizontal {
  border: none;
}
</style>

<template>
  <div class="header">
    <div class="header-main">
      <div class="header-left-side">
        <slot name="left" />
      </div>
      <div class="header-right-side">
        <slot name="right" />
      </div>
    </div>
  </div>
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
.header {
  border-bottom: 2px solid;
  border-color: #00adae;
}
.header-main {
  height: 60px;
  display: flex;
  padding: 0 20px;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
}

.header-left-side {
  height: 100%;
  display: flex;
  align-items: center;
}

.header-right-side {
  display: flex;
  align-items: center;
}
.warp__title-bar {
  background-color: var(--violet-white-color);
  margin-bottom: 10px;
}
</style>

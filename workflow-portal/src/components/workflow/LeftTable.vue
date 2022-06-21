<template>
  <div class="nav-table">
    <div v-if="!item">
      <img alt="Dont have items selected" src="../../assets/no_items.png" />
      <h6>No Items Selected</h6>
      <p>Seleted items to show details</p>
    </div>
    <div v-else>
      <p>
        Id : <span>{{ item.uuid }}</span>
      </p>
      <p>
        Name : <span>{{ item.name }}</span>
      </p>
      <p>
        Description : <span>{{ item.description }}</span>
      </p>
      <p>
        Created : <span>{{ item.created_time }}</span>
      </p>
      <p v-html="item.content"></p>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  name: "LeftTable",
  components: {},
  computed: {
    ...mapGetters({
      item: "stateItemsSelected",
    }),
  },
};
</script>

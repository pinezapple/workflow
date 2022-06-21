<template>
  <div class="nav-table">
    <div v-if="!item">
      <img alt="Dont have items selected" src="../../assets/no_items.png" />
      <h6>No Items Selected</h6>
      <p>Seleted items to show details</p>
    </div>
    <div v-else>
      <p>
        UIID : <span>{{ item.uuid }}</span>
      </p>
      <p>
        User : <span>{{ item.user }}</span>
      </p>
      <p>
        State : <span>{{ item.state }}</span>
      </p>
      <p>
        Start Time : <span>{{ item.start_time }}</span>
      </p>
      <p>
        End Time : <span>{{ item.end_time }}</span>
      </p>
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
      item: "stateItemsAnalysesSelected",
    }),
  },
};
</script>

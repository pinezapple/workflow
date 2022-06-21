<template>
  <div v-if="lists">
    <el-table
      ref="multipleTable"
      :data="lists?.tableData"
      :default-sort="{ prop: 'created_time', order: 'descending' }"
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55"> </el-table-column>
      <el-table-column
        prop="user"
        sortable
        property="user"
        label="User"
        width="120"
      >
      </el-table-column>
      <el-table-column
        prop="state"
        sortable
        property="state"
        label="State"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column
        prop="start_time"
        sortable
        property="start_time"
        label="Start Time"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column
        prop="end_time"
        sortable
        property="end_time"
        label="End Time"
        show-overflow-tooltip
      >
      </el-table-column>
    </el-table>
    <el-pagination
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :page-sizes="[10, 20, 50, 100]"
      :page-size="page_size"
      layout="total, sizes, prev, pager, next"
      :total="total"
    >
    </el-pagination>
  </div>
  <div style="margin-top: 20px">
    <el-button @click="toggleSelection([tableData[1], tableData[2]])"
      >Toggle selection status of second and third rows</el-button
    >
    <el-button @click="toggleSelection()">Clear selection</el-button>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  user: "TableAnalyses",
  components: {},

  data() {
    return {
      isShowFull: false,
    };
  },
  created: function () {
    this.GetAnalyses();
  },
  computed: {
    ...mapGetters({
      lists: "stateAnalysesList",
      item: "stateItemsAnalysesSelected",
      page: "getAnalysesPage",
      page_size: "getAnalysesPageSize",
      total: "getAnalysesTotal",
    }),
  },

  methods: {
    ...mapActions([
      "GetItemsSelectedAnalyses",
      "GetAnalyses",
      "ChangePageSizeAnalyses",
      "ChangePageAnalyses",
    ]),
    async handleSizeChange(val) {
      await this.ChangePageSizeAnalyses(val);
    },
    async handleCurrentChange(val) {
      await this.ChangePageAnalyses(val);
    },

    toggleSelection(rows) {
      if (rows) {
        rows.forEach((row) => {
          this.$refs.multipleTable.toggleRowSelection(row);
        });
      } else {
        this.$refs.multipleTable.clearSelection();
      }
    },
    async handleSelectionChange(val) {
      //this is for ui update
      this.multipleSelection = val;
      if (val.length) {
        await this.GetItemsSelectedAnalyses(val[0].uuid);
      } else {
        await this.GetItemsSelectedAnalyses();
      }
    },

    changeShow() {
      this.isShowFull = !this.isShowFull;
    },
  },
};
</script>

<style>
.icon-more {
  margin-top: 10px;
}
</style>

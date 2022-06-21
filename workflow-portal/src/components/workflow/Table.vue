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
      <el-table-column prop="created_time" sortable label="Date" width="120">
        <template #default="scope">{{ scope.row.created_time }}</template>
      </el-table-column>
      <el-table-column
        prop="name"
        sortable
        property="name"
        label="Name"
        width="120"
      >
      </el-table-column>
      <el-table-column
        prop="description"
        sortable
        property="description"
        label="Description"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column
        prop="tags.author"
        sortable
        property="tags.author"
        label="Author"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column
        prop="tags.version"
        sortable
        property="tags.version"
        label="Version"
        show-overflow-tooltip
      >
      </el-table-column>
      <el-table-column label="Steps">
        <template v-slot="props">
          <div v-for="(f, i) in props.row.steps" :key="i">
            <div v-if="!isShowFull">
              <div v-if="i <= 5">
                <el-link
                  v-on:click.stop.prevent="doThat"
                  icon="el-icon-download"
                  type="primary"
                >
                  {{ f.name }}
                </el-link>
              </div>
            </div>
            <div v-else>
              <el-link
                v-on:click.stop.prevent="doThat"
                icon="el-icon-download"
                type="primary"
              >
                {{ f.name }}
              </el-link>
            </div>
          </div>
          <div class="icon-more">
            <button>
              <el-link
                v-on:click.stop.prevent="changeShow()"
                icon="el-icon-more-outline"
              ></el-link>
            </button>
          </div>
        </template>
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
  name: "TableWorkflow",
  components: {},

  data() {
    return {
      isShowFull: false,
    };
  },
  created: function () {
    this.GetWorkflowList();
  },
  computed: {
    ...mapGetters({
      lists: "stateWorkflowList",
      item: "stateItemsSelected",
      page: "getPage",
      page_size: "getPageSize",
      total: "getTotal",
    }),
  },

  methods: {
    ...mapActions([
      "GetItemsSelected",
      "GetWorkflowList",
      "ChangePageSize",
      "ChangePage",
    ]),
    async handleSizeChange(val) {
      await this.ChangePageSize(val);
    },
    async handleCurrentChange(val) {
      await this.ChangePage(val);
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
        await this.GetItemsSelected(val[0].uuid);
      } else {
        await this.GetItemsSelected();
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

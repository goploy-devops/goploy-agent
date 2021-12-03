<template>
  <div class="app-container">
    <el-table :data="list" style="width: 100%">
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="expression" label="Expression" />
      <el-table-column prop="command" label="Command" />
      <el-table-column prop="singleMode" label="Single mode">
        <template #default="scope">
          <span v-if="scope.row.singleMode === 0">no</span>
          <span v-else>yes</span>
        </template>
      </el-table-column>
      <el-table-column prop="logLevel" label="Log level">
        <template #default="scope">
          <span v-if="scope.row.logLevel === 0">none</span>
          <span v-else-if="scope.row.logLevel === 1">stdout</span>
          <span v-else-if="scope.row.logLevel === 2">stdout+stderr</span>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="Description" />
      <el-table-column prop="creator" label="Creator" />
      <el-table-column prop="editor" label="Editor" />
      <el-table-column prop="insertTime" label="Insert time" width="135" />
      <el-table-column prop="updateTime" label="Update time" width="135" />
      <el-table-column prop="operation" width="180" align="center">
        <template #default="scope">
          <el-button
            type="primary"
            icon="el-icon-tickets"
            @click="handleLogs(scope.row)"
          >
            Logs
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="dialogVisible" title="Log">
      <el-row type="flex">
        <div v-loading="logLoading" class="cron-left">
          <el-radio-group
            v-model="selectedLogIndex"
            @change="handleSelectedLog"
          >
            <el-row v-for="(item, index) in logs" :key="index">
              <el-row style="margin: 5px 0">
                <el-radio class="cron-list" :label="index" border>
                  <span class="cron-time">{{ item.reportTime }}</span>
                  <i
                    v-if="item.execCode === 0"
                    class="el-icon-check icon-success"
                    style="float: right"
                  />
                  <i
                    v-else
                    class="el-icon-close icon-fail"
                    style="float: right"
                  />
                </el-radio>
              </el-row>
            </el-row>
          </el-radio-group>
          <el-pagination
            v-model:current-page="pagination.page"
            :total="pagination.total"
            :page-size="pagination.rows"
            style="text-align: right"
            layout="prev, next"
            @current-change="handlePageChange"
          />
        </div>
        <el-row
          class="cron-message"
          style="
            width: 100%;
            flex: 1;
            align-content: flex-start;
            white-space: pre-wrap;
          "
        >
          <div style="width: 100%">{{ message }}</div>
        </el-row>
      </el-row>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { CronList, CronLogs, CronData } from '@/api'
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'Index',
  data() {
    return {
      dialogVisible: false,
      list: [] as CronList['datagram'],
      logs: [] as CronLogs['datagram'],
      pagination: {
        total: 9999,
        page: 1,
        rows: 11,
      },
      selectedItem: {} as CronData['datagram'],
      selectedLogIndex: -1,
      logLoading: false,
      message: '',
    }
  },
  created() {
    this.getCronList()
  },
  methods: {
    enterToBR(detail: string) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    },
    handleLogs(data: CronData['datagram']) {
      this.dialogVisible = true
      this.selectedItem = data
      this.pagination.page = 1
      this.pagination.total = 9999
      this.getCronLogs(data.id)
    },
    handleSelectedLog(index: number) {
      this.message = this.logs[index].message
    },
    handlePageChange(page = 1) {
      this.pagination.page = page
      this.getCronLogs(this.selectedItem.id)
    },
    getCronList() {
      new CronList().request().then((response) => {
        this.list = response.data
      })
    },
    getCronLogs(id: number) {
      this.logLoading = true
      new CronLogs({ id }, this.pagination)
        .request()
        .then((response) => {
          this.logs = response.data
          if (this.logs.length < this.pagination.rows) {
            this.pagination.total = this.logs.length
          }
        })
        .finally(() => {
          this.logLoading = false
        })
    },
  },
})
</script>

<style lang="scss" scoped>
@import '@/styles/mixin.scss';

.cron {
  &-left {
    width: 180px;
  }
  &-list {
    margin-right: 5px;
    padding-right: 8px;
    width: 180px;
    line-height: 12px;
  }
  &-time {
    display: inline-block;
    text-align: center;
  }
  &-message {
    padding: 5px 0 0 15px;
    height: 470px;
    overflow-y: auto;
    @include scrollBar();
  }
}

.icon-success {
  color: #67c23a;
  font-size: 14px;
  font-weight: 900;
}

.icon-fail {
  color: #f56c6c;
  font-size: 14px;
  font-weight: 900;
}
</style>

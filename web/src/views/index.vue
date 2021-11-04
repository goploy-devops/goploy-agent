<template>
  <div class="app-container">
    <el-row :gutter="10">
      <el-col :xs="24" :sm="24" :md="8">
        <el-card>
          <template #header>
            <i class="el-icon-info"></i>
            <span style="margin-left: 4px">General Info</span>
            <el-button
              type="text"
              style="margin-left: 10px; font-size: 16px"
              icon="el-icon-refresh"
              @click="getGeneralInfo()"
            ></el-button>
          </template>
          <el-descriptions :column="1" size="medium">
            <el-descriptions-item label="Kernel Version">
              {{ generalInfo.kernelVersion }}
            </el-descriptions-item>
            <el-descriptions-item label="OS">
              {{ generalInfo.os }}
            </el-descriptions-item>
            <el-descriptions-item label="Cores">
              {{ generalInfo.cores }}
            </el-descriptions-item>
            <el-descriptions-item label="Hostname">
              {{ generalInfo.hostname }}
            </el-descriptions-item>
            <el-descriptions-item label="Uptime">
              {{ generalInfo.uptime }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8">
        <el-card>
          <template #header>
            <i class="el-icon-odometer"></i>
            <span style="margin-left: 4px">Load average</span>
            <el-button
              type="text"
              style="margin-left: 10px; font-size: 16px"
              icon="el-icon-refresh"
              @click="getLoadavgInfo()"
            ></el-button>
          </template>
          <el-row style="min-height: 160px">
            <el-col
              :span="8"
              style="
                text-align: center;
                font-size: 24px;
                border-right: solid 1px #e6e6e6;
              "
            >
              <div class="card-blue-title">1 min</div>
              <div class="card-value">
                {{ ((loadavgInfo.avg / loadavgInfo.cores) * 100).toFixed() }}
                <span class="card-unit">%</span>
              </div>
              <div class="card-value">{{ loadavgInfo.avg }}</div>
            </el-col>
            <el-col
              :span="8"
              style="
                text-align: center;
                font-size: 24px;
                border-right: solid 1px #e6e6e6;
              "
            >
              <div class="card-blue-title">5 min</div>
              <div class="card-value">
                {{ ((loadavgInfo.avg5 / loadavgInfo.cores) * 100).toFixed() }}
                <span class="card-unit">%</span>
              </div>
              <div class="card-value">{{ loadavgInfo.avg5 }}</div>
            </el-col>
            <el-col :span="8" style="text-align: center; font-size: 24px">
              <div>
                <div class="card-blue-title">15 min</div>
                <div class="card-value">
                  {{
                    ((loadavgInfo.avg15 / loadavgInfo.cores) * 100).toFixed()
                  }}
                  <span class="card-unit">%</span>
                </div>
                <div class="card-value">{{ loadavgInfo.avg15 }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8">
        <el-card>
          <template #header>
            <i class="el-icon-help"></i>
            <span style="margin-left: 4px">RAM</span>
            <el-button
              type="text"
              style="margin-left: 10px; font-size: 16px"
              icon="el-icon-refresh"
              @click="getRAMInfo()"
            ></el-button>
          </template>
          <el-row style="min-height: 160px">
            <el-col
              :span="8"
              style="
                text-align: center;
                font-size: 24px;
                border-right: solid 1px #e6e6e6;
              "
            >
              <div class="card-blue-title">Total</div>
              <div class="card-value">
                {{ humanSize(ramInfo.total) }}
              </div>
            </el-col>
            <el-col
              :span="8"
              style="
                text-align: center;
                font-size: 24px;
                border-right: solid 1px #e6e6e6;
              "
            >
              <div class="card-blue-title">Used</div>
              <div class="card-value">
                {{ humanSize(ramInfo.total - ramInfo.free) }}
              </div>
              <div class="card-value">
                {{ ((1 - ramInfo.free / ramInfo.total) * 100).toFixed() }}
                <span class="card-unit">%</span>
              </div>
            </el-col>
            <el-col :span="8" style="text-align: center; font-size: 24px">
              <div>
                <div class="card-blue-title">Free</div>
                <div class="card-value">
                  {{ humanSize(ramInfo.free) }}
                </div>
                <div class="card-value">
                  {{ ((ramInfo.free / ramInfo.total) * 100).toFixed() }}
                  <span class="card-unit">%</span>
                </div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <i class="el-icon-cpu"></i>
            <span style="margin-left: 4px">CPU</span>
            <el-button
              type="text"
              style="margin-left: 10px; font-size: 16px"
              icon="el-icon-refresh"
              @click="getCPUInfo()"
            ></el-button>
          </template>
          <el-table :data="cpuList" style="width: 100%">
            <el-table-column prop="0" label="name" />
            <el-table-column prop="1" label="user" />
            <el-table-column prop="2" label="nice" />
            <el-table-column prop="3" label="system" />
            <el-table-column prop="4" label="idle" />
            <el-table-column prop="5" label="iowait" />
            <el-table-column prop="6" label="irq" />
            <el-table-column prop="7" label="softirq" />
            <el-table-column prop="8" label="steal" />
            <el-table-column prop="9" label="guest" />
            <el-table-column prop="10" label="guest_nice" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <i class="el-icon-postcard"></i>
            <span style="margin-left: 4px">Disk Usage</span>
            <el-button
              type="text"
              style="margin-left: 10px; font-size: 16px"
              icon="el-icon-refresh"
              @click="getDiskUsageInfo()"
            ></el-button>
          </template>
          <el-table :data="diskUsageList" style="width: 100%">
            <el-table-column
              prop="filesystem"
              label="Filesystem"
              min-width="150"
            />
            <el-table-column prop="size" label="Size" />
            <el-table-column prop="used" label="Used" />
            <el-table-column prop="avail" label="Avail" />
            <el-table-column prop="usedPcent" label="Use%" />
            <el-table-column prop="mountedOn" label="Mounted on" />
            <el-table-column prop="inodes" label="Inodes" />
            <el-table-column prop="iUsed" label="IUsed" />
            <el-table-column prop="iFree" label="IFree" />
            <el-table-column prop="iUsedPcent" label="IUse%" />
            <el-table-column prop="type" label="Type" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <i class="el-icon-bank-card"></i>
            <span style="margin-left: 4px">Disk IO Stats</span>
            <el-button
              type="text"
              style="margin-left: 10px; font-size: 16px"
              icon="el-icon-refresh"
              @click="getDiskIOStats()"
            ></el-button>
          </template>
          <el-table :data="diskIOStats.list" style="width: 100%">
            <el-table-column
              v-for="(name, index) in diskIOStats.header"
              :key="index"
              :prop="'' + index"
              :label="name"
            />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <i class="el-icon-box"></i>
            <span style="margin-left: 4px">Net</span>
            <el-button
              type="text"
              style="margin-left: 10px; font-size: 16px"
              icon="el-icon-refresh"
              @click="getNetInfo()"
            ></el-button>
          </template>
          <el-table :data="netList" style="width: 100%">
            <el-table-column prop="0" label="Inter" />
            <el-table-column label="Receive">
              <el-table-column prop="1" label="bytes" />
              <el-table-column prop="2" label="packets" />
              <el-table-column prop="3" label="errs" />
              <el-table-column prop="4" label="drop" />
              <el-table-column prop="5" label="fifo" />
              <el-table-column prop="6" label="frame" />
              <el-table-column prop="7" label="compressed" />
              <el-table-column prop="8" label="multicast" />
            </el-table-column>
            <el-table-column label="Transmit">
              <el-table-column prop="9" label="bytes" />
              <el-table-column prop="10" label="packets" />
              <el-table-column prop="11" label="errs" />
              <el-table-column prop="12" label="drop" />
              <el-table-column prop="13" label="fifo" />
              <el-table-column prop="14" label="colls" />
              <el-table-column prop="15" label="carrier" />
              <el-table-column prop="16" label="compressed" />
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts">
import { General, Loadavg, RAM, CPU, Net, DiskUsage, DiskIOStat } from '@/api'
import { humanSize } from '@/utils'
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'Index',
  data() {
    return {
      generalInfo: {} as General['datagram'],
      loadavgInfo: {} as Loadavg['datagram'],
      ramInfo: {} as RAM['datagram'],
      cpuList: [],
      netList: [],
      diskUsageList: '',
      diskIOStats: { header: [], list: [] },
    }
  },
  created() {
    this.getGeneralInfo()
    this.getLoadavgInfo()
    this.getRAMInfo()
    this.getCPUInfo()
    this.getNetInfo()
    this.getDiskUsageList()
    this.getDiskIOStats()
  },
  methods: {
    humanSize,
    enterToBR(detail: string) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    },
    getGeneralInfo() {
      new General().request().then((response) => {
        this.generalInfo = response.data
      })
    },
    getLoadavgInfo() {
      new Loadavg().request().then((response) => {
        this.loadavgInfo = response.data
      })
    },
    getRAMInfo() {
      new RAM().request().then((response) => {
        this.ramInfo = response.data
      })
    },
    getCPUInfo() {
      new CPU().request().then((response) => {
        this.cpuList = response.data
      })
    },
    getNetInfo() {
      new Net().request().then((response) => {
        this.netList = response.data
      })
    },
    getDiskUsageList() {
      new DiskUsage().request().then((response) => {
        this.diskUsageList = response.data
      })
    },
    getDiskIOStats() {
      new DiskIOStat().request().then((response) => {
        this.diskIOStats = response.data
      })
    },
  },
})
</script>

<style lang="scss" scoped>
.app-container {
  .card-blue-title {
    margin-top: 16px;
    color: #409eff;
    padding: 6px 0;
  }
  .card-value {
    font-weight: 700;
    padding: 6px 0;
  }
  .card-unit {
    font-weight: 500;
    font-size: 14px;
  }
  .cmd-output {
    white-space: pre;
    font-size: 15px;
    line-height: 24px;
    font-family: monospace;
    letter-spacing: 0.5px;
    overflow-x: auto;
    color: #606266;
  }
}
</style>

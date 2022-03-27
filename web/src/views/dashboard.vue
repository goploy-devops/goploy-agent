<template>
  <div class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-date-picker
          v-model="datetimeRange"
          :shortcuts="shortcuts"
          type="datetimerange"
          range-separator="-"
          start-placeholder="Start date"
          end-placeholder="End date"
          @change="datetimeChange"
        />
        <el-button icon="el-icon-refresh" @click="handleRefresh" />
      </el-row>
    </el-row>
    <el-row class="chart-container" :gutter="10">
      <el-col
        v-for="(_, name) in chartNameMap"
        :key="name"
        :xs="24"
        :sm="24"
        :lg="12"
        style="margin-bottom: 10px"
      >
        <div
          :id="name"
          :ref="name"
          style="height: 288px; border: solid 1px #e6e6e6; padding: 10px 0"
        ></div>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts">
import { Chart } from '@/api'
import { deepClone, parseTime } from '@/utils'
import * as echarts from 'echarts'
import { ElDatePicker } from 'element-plus'
import dayjs, { Dayjs } from 'dayjs'
import { defineComponent } from 'vue'
const chartBaseOption = {
  title: {
    text: '',
    textStyle: {
      fontSize: 14,
    },
    padding: [5, 20],
  },
  tooltip: {
    trigger: 'axis',
  },
  xAxis: {
    type: 'time',
  },
  yAxis: {
    type: 'value',
  },
  legend: {
    bottom: 0,
    data: [],
  },
  series: [],
}
export default defineComponent({
  name: 'Index',
  data() {
    return {
      shortcuts: [
        {
          text: 'Last hour',
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
        {
          text: 'Last 6 hours',
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 6)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
        {
          text: 'Last 6 day',
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
        {
          text: 'Last 6 week',
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
      ],
      chartNameMap: <
        Record<string, { type: number; title: string; subtitle: string }>
      >{
        cpuChart: {
          type: 1,
          title: 'CPU Usage',
          subtitle: '(%)',
        },
        ramChart: {
          type: 2,
          title: 'RAM Usage',
          subtitle: '(%)',
        },
        loadavgChart: {
          type: 3,
          title: 'Loadavg',
          subtitle: '',
        },
        tcpChart: {
          type: 4,
          title: 'TCP',
          subtitle: '(count)',
        },
        pubNetChart: {
          type: 5,
          title: 'Pub Band width',
          subtitle: '(bit/s)',
        },
        loNetChart: {
          type: 6,
          title: 'Lo Band width',
          subtitle: '(bit/s)',
        },
        diskUsageChart: {
          type: 7,
          title: 'Disk Usage',
          subtitle: '(%)',
        },
        diskIOChart: {
          type: 8,
          title: 'Disk IO',
          subtitle: '(count/s)',
        },
      },
      datetimeRange: <Dayjs[]>[],
    }
  },
  mounted() {
    this.handleRefresh()
  },
  methods: {
    handleRefresh() {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000)
      this.datetimeRange = [dayjs(start), dayjs(end)]
      this.datetimeChange(
        this.datetimeRange.map((datetime) => datetime.toDate())
      )
    },
    datetimeChange(values: Date[]) {
      for (const key in this.chartNameMap) {
        this.chart(key, values)
      }
    },
    chart(chartName: string, values: Date[]) {
      new Chart({
        type: this.chartNameMap[chartName].type,
        datetimeRange: values.map((value) => parseTime(value)).join(','),
      })
        .request()
        .then((response) => {
          echarts.dispose(<HTMLDivElement>document.getElementById(chartName))
          let chart = echarts.init(
            <HTMLDivElement>document.getElementById(chartName)
          )
          let chartOption = deepClone(chartBaseOption)
          chartOption.title.text = this.chartNameMap[chartName].title
          chartOption.title.subtext = this.chartNameMap[chartName].subtitle
          for (const key in response.data.map) {
            chartOption.legend.data.push(key)
            const series = {
              name: key,
              type: 'line',
              symbol: 'none',
              smooth: true,
              data: response.data.map[key].map(
                (item: { reportTime: string; value: string }) => {
                  return [new Date(item.reportTime), item.value]
                }
              ),
            }
            chartOption.series.push(series)
          }
          chart.setOption(chartOption)
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
}

@import '@/styles/mixin.scss';
.chart-container {
  width: 100%;
  max-height: calc(100vh - 160px);
  overflow-y: auto;
  @include scrollBar();
}
</style>

<template>
  <div class="flow-chart">
    <div ref="chartRef" class="chart-container"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import {getFlowStat} from '~/api/dashboard.js'

const props = defineProps({
  data: {
    type: Object,
    default: () => ({ today: [], yesterday: [] })
  }
})

getFlowStat().then(res => {
  // console.log("getFlowStat",res.data.data)
  chartData.value.today = res.data.data.today
  chartData.value.yesterday = res.data.data.yesterday
  if (!chartInstance) return
  chartInstance.setOption({
    series: [
      { data: chartData.value.today },
      { data: chartData.value.yesterday }
    ]
  })
})

const chartData = ref({
  today: [],
  yesterday: []
})
let chartInstance = null
const chartRef = ref(null)

const initChart = () => {
  if (!chartRef.value) return
  chartInstance = echarts.init(chartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['今日流量', '昨日流量']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: Array.from({ length: 24 }, (_, i) => `${i}:00`)
    },
    yAxis: {
      type: 'value',
      name: '请求量',
      axisLabel: {
        formatter: '{value}'
      }
    },
    series: [
      {
        name: '今日流量',
        type: 'line',
        stack: 'Total',
        areaStyle: {},
        emphasis: {
          focus: 'series'
        },
        data: props.data.today
      },
      {
        name: '昨日流量',
        type: 'line',
        stack: 'Total',
        areaStyle: {},
        emphasis: {
          focus: 'series'
        },
        data: props.data.yesterday
      }
    ]
  }

  chartInstance.setOption(option)
}

// 先初始化一个空图表（onMounted）
onMounted(() => {
  initChart() // 初始化空图
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  if (chartInstance) {
    chartInstance.dispose()
  }
  window.removeEventListener('resize', handleResize)
})

const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}
</script>

<style scoped>
.flow-chart {
  width: 100%;
  height: 100%;
}

.chart-container {
  width: 100%;
  height: 100%;
}
</style>
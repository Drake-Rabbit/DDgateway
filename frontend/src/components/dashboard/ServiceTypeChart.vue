<template>
  <div class="service-type-chart">
    <div ref="chartRef" class="chart-container"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import { getServiceStat } from '~/api/dashboard'


// 本地响应式数据
const chartData = ref({
  legend: [],
  data: [] // 格式：[{ name: 'HTTP', value: 2 }, ...]
})

const chartRef = ref(null)
let chartInstance = null

getServiceStat().then(res => {
  const apiResponse = res.data.data // 这是一个对象：{ legend: [...], data: [...] }

  // ✅ 正确提取
  chartData.value.legend = apiResponse.legend || []
  chartData.value.data = apiResponse.data || []
  // 如果图表已初始化，更新它
  if (chartInstance) {
    chartInstance.setOption({
      legend: { data: chartData.value.legend },
      series: [{ data: chartData.value.data }]
    })
  }
})

const initChart = () => {
  if (!chartRef.value) return

  chartInstance = echarts.init(chartRef.value)

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      data: chartData.value.legend  
    },
    series: [
      {
        name: '服务类型',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '20',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: chartData.value.data.map(item => ({
          value: item.value,
          name: item.name
        }))
      }
    ]
  }

  chartInstance.setOption(option)
}

onMounted(() => {
  initChart()
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
.service-type-chart {
  width: 100%;
  height: 100%;
}

.chart-container {
  width: 100%;
  height: 100%;
}
</style>
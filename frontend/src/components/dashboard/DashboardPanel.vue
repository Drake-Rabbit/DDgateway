<template>
  <div class="dashboard-container" >
    <el-card class="card-panel">
      <template #header>
        <div class="card-header">
          <span>仪表盘概览</span>
        </div>
      </template>

      <div class="panel-group">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon">
                  <i class="el-icon-s-data"></i>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ panelData?.serviceNum || 0 }}</div>
                  <div class="stat-label">服务总数</div>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon">
                  <i class="el-icon-apple"></i>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ panelData?.appNum || 0 }}</div>
                  <div class="stat-label">应用总数</div>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon">
                  <i class="el-icon-timer"></i>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ panelData?.todayRequestNum || 0 }}</div>
                  <div class="stat-label">今日请求量</div>
                </div>
              </div>
            </el-card>
          </el-col>

          <el-col :span="6">
            <el-card class="stat-card">
              <div class="stat-content">
                <div class="stat-icon">
                  <i class="el-icon-speedometer"></i>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ panelData?.currentQPS || 0 }}</div>
                  <div class="stat-label">当前QPS</div>
                </div>
              </div>
            </el-card>
          </el-col>

          
        </el-row>
      </div>

    </el-card>

    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>流量统计</span>
            </div>
          </template>
          <FlowChart :data="flowData" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>服务类型分布</span>
            </div>
          </template>
          <ServiceTypeChart :data="serviceStatData" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import FlowChart from './FlowChart.vue'
import ServiceTypeChart from './ServiceTypeChart.vue'
import { getPanelGroupData, getFlowStat, getServiceStat } from '~/api/dashboard'

const panelData = ref(null)
const flowData = ref({ today: [], yesterday: [] })
const serviceStatData = ref([])

const fetchPanelData = async () => {
  try {
    const response = await getPanelGroupData()
    panelData.value = response.data
  } catch (error) {
    ElMessage.error('获取面板数据失败')
  }
}

const fetchFlowStat = async () => {
  try {
    const response = await getFlowStat()
    flowData.value = response.data
  } catch (error) {
    ElMessage.error('获取流量统计失败')
  }
}

const fetchServiceStat = async () => {
  try {
    const response = await getServiceStat()
    serviceStatData.value = response.data
  } catch (error) {
    ElMessage.error('获取服务统计失败')
  }
}

onMounted(() => {
  fetchPanelData()
  fetchFlowStat()
  fetchServiceStat()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.card-panel {
  margin-bottom: 20px;
}

.panel-group {
  margin-top: 20px;
}

.stat-card {
  text-align: center;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 0;
}

.stat-icon {
  font-size: 28px;
  color: #409EFF;
  margin-right: 15px;
}

.stat-info {
  text-align: left;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.chart-row {
  margin-top: 20px;
}

.chart-card {
  height: 400px;
}

.dashboard-container {
  padding: 0px;
  @apply w-full h-full ;
}
</style>
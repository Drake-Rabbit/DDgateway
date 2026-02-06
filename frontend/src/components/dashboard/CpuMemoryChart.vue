<template>
  <div class="cpu-memory-container">
    <el-card span="4" class="stat-card">
      <div class="stat-content">
        <div class="stat-icon">
          <i class="el-icon-monitor"></i>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ cpuUsage }}%</div>
          <div class="stat-label">CPU使用率</div>
        </div>
      </div>
    </el-card>

    <el-card span="4" class="stat-card">
      <div class="stat-content">
        <div class="stat-icon">
          <i class="el-icon-memory"></i>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ memoryUsage }}%</div>
          <div class="stat-label">内存使用率</div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const cpuUsage = ref(0)
const memoryUsage = ref(0)

const fetchSystemStats = async () => {
  try {
    // 模拟获取系统统计数据
    // 实际项目中应该调用后端API获取真实数据
    cpuUsage.value = Math.floor(Math.random() * 80) + 10 // 模拟10-90%的CPU使用率
    memoryUsage.value = Math.floor(Math.random() * 70) + 20 // 模拟20-90%的内存使用率
  } catch (error) {
    console.error('获取系统统计数据失败:', error)
  }
}

onMounted(() => {
  fetchSystemStats()
  // 每5秒更新一次数据
  setInterval(fetchSystemStats, 5000)
})
</script>

<style scoped>
.cpu-memory-container {
  display: flex;
  gap: 20px;
  width: 100%;
}

.stat-card {
  flex: 1;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
  text-align: center;
  padding: 14px 0;
  /* height: 110px; */
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
}

/* CPU使用率 - 紫色 */
.stat-card:nth-child(1) {
  background-color: #ce6ccc;
  color: white;
}

/* 内存使用率 - 绿色 */
.stat-card:nth-child(2) {
  background-color: #67c23a;
  color: white;
}

.stat-icon i {
  font-size: 28px;
  margin-bottom: 8px;
  color: inherit;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}
</style>
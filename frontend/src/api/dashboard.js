import axios from '~/axios.js'

// 获取面板数据
export function getPanelGroupData() {
  return axios.get('/dashboard/panelGruopData')
}

// 获取流量统计
export function getFlowStat() {
  return axios.get('/dashboard/flowstat')
}

// 获取服务统计
export function getServiceStat() {
  return axios.get('/dashboard/serviceStat')
}

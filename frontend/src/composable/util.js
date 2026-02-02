//封装工具类
import { ElMessageBox } from 'element-plus'

//消息弹出确认框
export function showModal(content,type,title) {
  return ElMessageBox.confirm(
      content,
      title, 
  {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: type,
  })
}
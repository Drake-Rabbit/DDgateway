import { reactive } from 'vue'

// 消息状态管理
const state = reactive({
  messages: []
})

let messageId = 0

export const useMessage = () => {
  // 添加消息
  const addMessage = (text, type = 'info') => {
    const id = ++messageId
    state.messages.push({ id, text, type })

    // 3秒后自动移除
    setTimeout(() => {
      removeMessage(id)
    }, 3000)

    return id
  }

  // 移除消息
  const removeMessage = (id) => {
    const index = state.messages.findIndex(m => m.id === id)
    if (index > -1) {
      state.messages.splice(index, 1)
    }
  }

  // 成功消息
  const success = (text) => {
    return addMessage(text, 'success')
  }

  // 错误消息
  const error = (text) => {
    return addMessage(text, 'error')
  }

  // 信息消息
  const info = (text) => {
    return addMessage(text, 'info')
  }

  // 警告消息
  const warning = (text) => {
    return addMessage(text, 'warning')
  }

  return {
    messages: state.messages,
    addMessage,
    removeMessage,
    success,
    error,
    info,
    warning
  }
}

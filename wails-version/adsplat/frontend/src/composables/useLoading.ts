import { ref, readonly } from 'vue'

// 简单的全屏loading状态管理
const isLoading = ref(false)
const loadingText = ref('加载中...')

export const useLoading = () => {
    const show = (text: string = '加载中...') => {
        loadingText.value = text
        isLoading.value = true
    }

    const hide = () => {
        isLoading.value = false
    }

    const updateText = (text: string) => {
        loadingText.value = text
    }

    return {
        isLoading: readonly(isLoading),
        loadingText: readonly(loadingText),
        show,
        hide,
        updateText
    }
}

// 导出只读的响应式状态，供组件使用
export { isLoading, loadingText }
<template>
    <div class="terminal">
        <pre>{{ rawText }}</pre>

        <el-button type="primary" @click="handleParse">粘贴</el-button>
        <el-button type="primary" @click="handleRun" :disabled="isDisabled">执行</el-button>
    </div>
</template>
  
<script setup>
import { ref, defineProps, nextTick } from 'vue'
const props = defineProps({
    serverId: { type: String, required: true },
    projectId: { type: String, required: true },
})

const rawText = ref('请点击粘贴按钮，粘贴内容')

// 执行按钮禁用状态
const isDisabled = ref(true)

// 粘贴
const handleParse = async () => {
    const text = (await navigator.clipboard.readText()) || ''
    if (!text) {
        rawText.value = '剪贴板为空'
        return
    }
    rawText.value = text
    isDisabled.value = false
}

// 执行
const handleRun = async () => {
    console.log('执行', rawText.value)
}
</script>
  
  <style scoped>
.terminal /deep/ pre {
    background-color: #f5f5f5;
    padding: 10px;
    border-radius: 5px;
    white-space: pre-wrap;
    width: auto;
    min-height: 200px;
    max-height: calc(100vh - 200px);
    overflow-y: auto;
}
</style>
  
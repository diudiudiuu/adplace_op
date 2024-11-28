<template>
    <div class="terminal">
        <div>
            <el-upload
                class="upload-demo"
                drag
                :auto-upload="false"
                :before-upload="handleFileRead"
            >
                <el-icon class="el-icon--upload">
                    <upload-filled />
                </el-icon>
                <div class="el-upload__text">
                    拖拽文件到此处或
                    <em>点击上传</em>
                </div>
                <template #tip>
                    <div class="el-upload__tip">txt md等文件</div>
                </template>
            </el-upload>
        </div>
        <div>
            <pre>{{ rawText }}</pre>
            <el-progress :indeterminate="false" />
            <el-button type="primary">启动</el-button>
        </div>
    </div>
</template>
  
  <script setup>
import { ref, defineProps } from 'vue'
import { UploadFilled } from '@element-plus/icons-vue'

const props = defineProps({
    serverId: { type: String, required: true },
    projectId: { type: String, required: true },
})

const rawText = ref(`:::demo
  <div>水电费都是</div>
  <template>
    <el-button type="primary">Primary Button</el-button>
  </template>
  :::`)

const handleFileRead = (file) => {
    const reader = new FileReader()
    reader.onload = (e) => {
        rawText.value = e.target.result
    }
    reader.readAsText(file)
    return false // 阻止默认上传
}
</script>
  
  <style scoped>
.terminal /deep/ pre {
    background-color: #f5f5f5;
    padding: 10px;
    border-radius: 5px;
    white-space: pre-wrap;
}
</style>
  
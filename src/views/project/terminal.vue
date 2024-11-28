<template>
    <div class="terminal">
        <div v-if="!tableData.length">
            <pre>{{ rawText }}</pre>
            <el-button type="primary" @click="handleParse">粘贴</el-button>
            <el-button type="primary" @click="handleAnalyze" :disabled="isDisabled">解析</el-button>
        </div>
        <el-table v-else :data="tableData">
            <el-table-column prop="title" label="命令" />
            <el-table-column prop="content" label="内容" />
            <el-table-column label="状态">
                <template #default="{ row }">
                    <el-tag type="success" v-if="row.status == 0">等待执行</el-tag>
                    <el-tag type="info" v-else-if="row.status == 1">执行中</el-tag>
                    <el-tag type="success" v-else-if="row.status == 2">执行成功</el-tag>
                    <el-tag type="danger" v-else-if="row.status == 3">执行失败</el-tag>
                </template>
            </el-table-column>

            <el-table-column label="操作">
                <template #default="{ row }">
                    <!-- 执行 -->
                    <el-button
                        type="primary"
                        size="small"
                        round
                        @click="handleExecute(row)"
                        :loading="row.status == 1"
                    >执行</el-button>
                    <!-- 删除 -->
                    <el-button type="danger" size="small" round @click="handleDelete(row)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
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

// 解析
const tableData = ref([])
const handleAnalyze = async () => {
    // 按换行符分割并清理空白行
    const lines = rawText.value
        .split('\n')
        .map((line) => line.trim())
        .filter((line) => line)

    // 按 "-- 命令" 分组
    const commands = []
    let currentCommand = null

    lines.forEach((line) => {
        if (line.startsWith('--')) {
            if (currentCommand) {
                commands.push(currentCommand)
            }
            currentCommand = { title: line, content: '', status: 0 }
        } else if (currentCommand) {
            currentCommand.content +=
                (currentCommand.content ? '\n' : '') + line
        }
    })

    // 添加最后一个命令
    if (currentCommand) {
        commands.push(currentCommand)
    }

    tableData.value = commands
}

// 执行
const handleExecute = async (row) => {
    console.log('执行', row)
    row.status = 1
}

// 删除
const handleDelete = async (row) => {
    tableData.value = tableData.value.filter((item) => item !== row)
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
  
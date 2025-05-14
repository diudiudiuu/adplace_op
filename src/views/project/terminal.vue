<template>
    <div class="terminal">
        <div v-if="!tableData.length" class="parse">
            <el-button type="primary" @click="handleParse">粘贴</el-button>
            <el-button type="primary" @click="handleAnalyze" :disabled="isDisabled">解析</el-button>
            <el-divider />
            <pre>{{ rawText }}</pre>
        </div>

        <div v-else class="exec">
            <el-button type="primary" @click="handleCancle">返回</el-button>
            <el-divider />
            <el-progress
                :text-inside="true"
                :stroke-width="18"
                :percentage="progress"
                status="success"
                striped
                striped-flow
            />
            <el-divider />
            <el-table :data="tableData">
                <el-table-column
                    prop="title"
                    label="Title"
                    align="left"
                    width="180px"
                    show-overflow-tooltip
                ></el-table-column>
                <el-table-column prop="content" label="SQL" align="left">
                    <template #default="{ row }">
                        <pre>{{ row.content }}</pre>
                    </template>
                </el-table-column>
                <el-table-column label="状态" width="120px">
                    <template #default="{ row }">
                        <el-tag v-if="row.status == 0">等待执行</el-tag>
                        <el-tag type="info" v-else-if="row.status == 1">执行中</el-tag>
                        <el-tag type="success" v-else-if="row.status == 2">操作完成</el-tag>
                        <el-tag type="danger" v-else-if="row.status == 3">系统异常</el-tag>
                    </template>
                </el-table-column>

                <el-table-column label="操作" width="180px">
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
                        <el-button
                            type="danger"
                            size="small"
                            round
                            :disabled="row.status == 1"
                            @click="handleDelete(row)"
                        >删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
    </div>
</template>
  
<script setup>
import { ref, defineProps, nextTick } from 'vue'
import api from '@/api'

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
const progress = ref(0)

// 计算进度
const progressTimer = setInterval(() => {
    if (tableData.value.length === 0) {
        progress.value = 0
        return
    }

    const total = tableData.value.length
    const done = tableData.value.filter((item) => item.status === 2).length
    progress.value = Math.floor((done / total) * 100)
}, 800)

const handleAnalyze = async () => {
    tableData.value = []
    progress.value = 0

    // 按换行符分割并清理空白行
    const lines = rawText.value
        .split('\n')
        .map((line) => line.trim())
        .filter((line) => line)

    // 按 "-- 命令" 分组
    const commands = []
    let currentCommand = null

    for (const line of lines) {
        if (line.startsWith('--')) {
            if (currentCommand) {
                commands.push(currentCommand)
            }
            currentCommand = { title: line, content: '', status: 0 }
        } else if (currentCommand) {
            currentCommand.content +=
                (currentCommand.content ? '\n' : '') + line
        }
    }

    // 添加最后一个命令
    if (currentCommand) {
        commands.push(currentCommand)
    }

    // 如果有命令，提示没有可执行的命令
    if (commands.length === 0) {
        rawText.value = '没有可执行的命令'
    } else {
        // 去除 commands 中的空白行
        tableData.value = commands.filter((item) => item.content)
    }
}

// 执行
const handleExecute = async (row) => {
    row.status = 1
    api('exec', {
        projectId: props.projectId,
        sql: row.content,
        sqlType: 'insert',
        authorization: localStorage.getItem('authorization'),
    })
        .then((res) => {
            if (res.code === 200) {
                setTimeout(() => {
                    row.status = 2
                    progressTimer()
                }, 800)
            } else {
                row.status = 3
            }
        })
        .catch(() => {
            row.status = 3
        })
}

// 删除
const handleDelete = async (row) => {
    tableData.value = tableData.value.filter((item) => item !== row)
    progressTimer()
}

// 返回
const handleCancle = async () => {
    tableData.value = []
    rawText.value = '请点击粘贴按钮，粘贴内容'
    isDisabled.value = true
}
</script>
  
  <style scoped>
.terminal .parse /deep/ pre {
    background-color: #f5f5f5;
    padding: 10px;
    border-radius: 5px;
    white-space: pre-wrap;
    width: auto;
    min-height: 200px;
    max-height: calc(100vh - 200px);
    overflow-y: auto;
}

.terminal .exec /deep/ pre {
    background-color: #f5f5f5;
    padding: 10px;
    border-radius: 5px;
    white-space: pre-wrap;
}
</style>
  
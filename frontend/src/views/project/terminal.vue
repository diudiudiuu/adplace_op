<template>
    <div class="terminal">
        <div v-if="!tableData.length" class="parse">
            <n-space>
                <n-tooltip>
                    <template #trigger>
                        <n-button type="primary" @click="handleParse">
                            <template #icon>
                                <n-icon><ClipboardOutline /></n-icon>
                            </template>
                            粘贴
                        </n-button>
                    </template>
                    粘贴
                </n-tooltip>
                <n-tooltip>
                    <template #trigger>
                        <n-button type="info" @click="handleAnalyze" :disabled="isDisabled">
                            <template #icon>
                                <n-icon><CodeOutline /></n-icon>
                            </template>
                            解析
                        </n-button>
                    </template>
                    解析
                </n-tooltip>
            </n-space>
            <n-divider />
            <n-card>
                <pre class="terminal-pre">{{ rawText }}</pre>
            </n-card>
        </div>

        <div v-else class="exec">
            <n-tooltip>
                <template #trigger>
                    <n-button type="primary" @click="handleCancle">
                        <template #icon>
                            <n-icon><ArrowBackOutline /></n-icon>
                        </template>
                        返回
                    </n-button>
                </template>
                返回
            </n-tooltip>
            <n-divider />
            <n-progress
                type="line"
                :percentage="progress"
                :height="18"
                :border-radius="4"
                fill-border-radius="4"
                processing
            />
            <n-divider />
            <n-data-table
                :columns="columns"
                :data="tableData"
                :pagination="false"
            />
        </div>
    </div>
</template>
  
<script setup>
import { ref, defineProps, computed, h } from 'vue'
import { NButton, NIcon, NTag, NSpace } from 'naive-ui'
import { 
    ClipboardOutline, 
    CodeOutline, 
    ArrowBackOutline,
    PlayOutline,
    TrashOutline
} from '@vicons/ionicons5'
import api from '@/api'

const props = defineProps({
    serverId: { type: String, required: true },
    projectId: { type: String, required: true },
})

const rawText = ref('请点击粘贴按钮，粘贴内容')

// 执行按钮禁用状态
const isDisabled = ref(true)

// 表格列配置
const columns = computed(() => [
    {
        title: 'Title',
        key: 'title',
        width: 180,
        ellipsis: {
            tooltip: true
        }
    },
    {
        title: 'SQL',
        key: 'content',
        render: (row) => {
            return h('pre', { class: 'sql-content' }, row.content)
        }
    },
    {
        title: '状态',
        key: 'status',
        width: 120,
        render: (row) => {
            const statusMap = {
                0: { type: 'default', text: '等待执行' },
                1: { type: 'info', text: '执行中' },
                2: { type: 'success', text: '操作完成' },
                3: { type: 'error', text: '系统异常' }
            }
            const status = statusMap[row.status] || statusMap[0]
            return h(NTag, { type: status.type }, { default: () => status.text })
        }
    },
    {
        title: '操作',
        key: 'actions',
        width: 180,
        render: (row) => {
            return h(NSpace, { size: 'small' }, {
                default: () => [
                    h(NButton, {
                        type: 'primary',
                        size: 'small',
                        loading: row.status === 1,
                        onClick: () => handleExecute(row)
                    }, {
                        default: () => '执行',
                        icon: () => h(NIcon, null, { default: () => h(PlayOutline) })
                    }),
                    h(NButton, {
                        type: 'error',
                        size: 'small',
                        disabled: row.status === 1,
                        onClick: () => handleDelete(row)
                    }, {
                        default: () => h(NIcon, null, { default: () => h(TrashOutline) })
                    })
                ]
            })
        }
    }
])

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
.terminal-pre {
    padding: 10px;
    border-radius: 8px;
    white-space: pre-wrap;
    width: auto;
    min-height: 200px;
    max-height: calc(100vh - 200px);
    overflow-y: auto;
    margin: 0;
}

.sql-content {
    padding: 8px;
    border-radius: 8px;
    white-space: pre-wrap;
    margin: 0;
    font-size: 12px;
}
</style>
  
<template>
    <div>
        <div class="search-box">
            <n-space>
                <n-tooltip>
                    <template #trigger>
                        <ColorfulButton v-if="!isMaxRecordsReached" type="add" text="添加" size="medium"
                            @click="openForm(false)" />
                        <n-button v-else disabled size="medium" quaternary>
                            <template #icon>
                                <n-icon>
                                    <AddCircleOutline />
                                </n-icon>
                            </template>
                            添加
                        </n-button>
                    </template>
                    {{ isMaxRecordsReached ? maxRecordsMessage : '添加' }}
                </n-tooltip>
                <n-tooltip>
                    <template #trigger>
                        <ColorfulButton type="refresh" text="刷新" size="medium" @click="refreshData" />
                    </template>
                    刷新
                </n-tooltip>
            </n-space>
        </div>

        <n-data-table :columns="columns" :data="tableData" :pagination="pagination" @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange" striped class="special-table" />

        <n-modal v-model:show="isFormVisible" preset="dialog" :title="isEditMode ? '编辑' : '添加'"
            style="width: 600px; max-width: 90vw;">
            <n-form :model="formData" label-placement="left" label-width="180">
                <n-form-item v-for="field in fields" :key="field" :label="field" v-if="shouldShowField(field)" required>
                    <template v-if="fieldsType[field]['type'] === 'enum'">
                        <n-radio-group v-model:value="formData[field]">
                            <n-radio-button v-for="item in fieldsType[field]['value']" :key="item" :value="item">
                                {{ item }}
                            </n-radio-button>
                        </n-radio-group>
                    </template>
                    <template v-else-if="fieldsType[field]['type'] === 'datetime'">
                        <n-date-picker v-model:value="formData[field]" type="datetime" placeholder="选择日期时间"
                            value-format="yyyy-MM-dd HH:mm:ss" format="yyyy-MM-dd HH:mm:ss" style="width: 100%"
                            @update:value="handleChange(field, formData, fieldsType[field])" />
                    </template>
                    <template v-else>
                        <n-input-group v-if="fieldsType[field]['button']">
                            <n-input v-model:value="formData[field]"
                                :disabled="(isPrimaryKey(field) && isEditMode) || fieldsType[field]['disabled']" />
                            <n-button type="primary" @click="handleClick(field, formData, fieldsType[field])">
                                {{ fieldsType[field]['button'].text }}
                            </n-button>
                        </n-input-group>
                        <n-input v-else v-model:value="formData[field]"
                            :disabled="(isPrimaryKey(field) && isEditMode) || fieldsType[field]['disabled']"
                            :placeholder="getFieldPlaceholder(field)" />
                    </template>
                </n-form-item>
            </n-form>
            <template #action>
                <n-space>
                    <n-button @click="isFormVisible = false">
                        <template #icon>
                            <n-icon>
                                <CloseOutline />
                            </n-icon>
                        </template>
                        取消
                    </n-button>
                    <n-button type="primary" @click="isEditMode ? submitForm('update') : submitForm('insert')">
                        <template #icon>
                            <n-icon>
                                <CheckmarkOutline />
                            </n-icon>
                        </template>
                        {{ isEditMode ? '更新' : '提交' }}
                    </n-button>
                </n-space>
            </template>
        </n-modal>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, defineProps, computed, h } from 'vue'
import { useMessage, useDialog, NButton, NIcon, NSpace, NTooltip } from 'naive-ui'
import { AddCircleOutline, RefreshOutline, CreateOutline, TrashOutline, CloseOutline, CheckmarkOutline } from '@vicons/ionicons5'
import ColorfulButton from '@/components/ColorfulButton.vue'
import api from '@/api'
import { encryptAes } from '@/utils'

const props = defineProps({
    model: { type: Object, required: true },
    projectId: { type: String, required: true },
})

const message = useMessage()
const dialog = useDialog()

const tableData = ref([])
const formData = ref({})
const isFormVisible = ref(false)
const isEditMode = ref(false)

const fields = props.model.fields
const fieldsType = props.model.fieldsType
const primaryKey = props.model.primaryKey

// 检查是否达到最大记录数限制
const isMaxRecordsReached = computed(() => {
    if (props.model.maxRecords === null || props.model.maxRecords === undefined) {
        return false; // 不限制
    }
    return tableData.value.length >= props.model.maxRecords;
})

// 获取限制提示信息
const maxRecordsMessage = computed(() => {
    return props.model.maxRecordsMessage || '已达到最大记录数限制';
})

// 分页配置
const pagination = ref({
    page: 1,
    pageSize: 10,
    showSizePicker: true,
    pageSizes: [10, 20, 50],
    showQuickJumper: true,
    itemCount: 0,
    prefix: ({ itemCount }) => `共 ${itemCount} 条`
})

// 表格列配置
const columns = computed(() => {
    const actionColumn = {
        title: '操作',
        key: 'actions',
        width: 120,
        render: (row) => {
            return h(NSpace, { size: 'small' }, {
                default: () => [
                    h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(ColorfulButton, {
                            type: 'edit',
                            size: 'small',
                            iconOnly: true,
                            onClick: () => openForm(true, row)
                        }),
                        default: () => '编辑'
                    }),
                    h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(ColorfulButton, {
                            type: 'delete',
                            size: 'small',
                            iconOnly: true,
                            onClick: () => confirmDelete(row)
                        }),
                        default: () => '删除'
                    })
                ]
            })
        }
    }

    const fieldColumns = fields.map((field: string) => ({
        title: field,
        key: field,
        ellipsis: {
            tooltip: true
        }
    }))

    return [actionColumn, ...fieldColumns]
})

const fetchData = async () => {
    try {
        const res = await api('exec', {
            projectId: props.projectId,
            sql: props.model.selects(),
            sqlType: 'selects',
            authorization: localStorage.getItem('authorization'),
        })
        if (res.code === 200) {
            tableData.value = res.data.result || []
            // 更新分页信息
            pagination.value.itemCount = tableData.value.length
        } else {
            message.error('获取数据失败')
            tableData.value = []
            pagination.value.itemCount = 0
        }
    } catch (error) {
        console.error('Failed to fetch data:', error)
        message.error('获取数据失败')
        tableData.value = []
        pagination.value.itemCount = 0
    }
}
onMounted(fetchData)

const refreshData = () => {
    pagination.value.page = 1 // 重置到第一页
    fetchData()
}

// 处理分页变化
const handlePageChange = (page: number) => {
    pagination.value.page = page
}

// 处理每页大小变化
const handlePageSizeChange = (pageSize: number) => {
    pagination.value.pageSize = pageSize
    pagination.value.page = 1 // 重置到第一页
}

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const isPrimaryKey = (field: any) => field === primaryKey
// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const shouldShowField = (field: any) => isEditMode.value || field !== primaryKey

const openForm = async (editMode: boolean, row = {}) => {
    // 如果是添加模式，检查是否达到最大记录数限制
    if (!editMode && isMaxRecordsReached.value) {
        message.warning(maxRecordsMessage.value)
        return
    }

    isEditMode.value = editMode
    formData.value = { ...row }
    if (!editMode) {
        // 如果是套餐管理且是添加模式，先获取项目信息设置默认端口
        if (props.model.constructor.name === 'Client') {
            await setDefaultPortsFromProject()
        }
        
        for (const field in fieldsType) {
            const type = fieldsType[field].type
            if (type === 'int' || type === 'string') {
                // 如果是端口字段且已经设置了默认值，跳过
                if ((field === 'api_port' || field === 'front_port') && formData.value[field]) {
                    continue
                }
                formData.value[field] = fieldsType[field].value
            }
            if (type === 'enum') {
                formData.value[field] =
                    fieldsType[field].default || fieldsType[field].value[0]
            }
            if (type === 'datetime') {
                // 默认值为当前时间
                formData.value[field] = new Date(
                    new Date().getTime() -
                    new Date().getTimezoneOffset() * 60000
                )
                    .toISOString()
                    .slice(0, 19)
                    .replace('T', ' ')
            }
        }
    }
    isFormVisible.value = true
}

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const confirmDelete = (row: any) => {
    dialog.warning({
        title: '确认删除',
        content: '此操作将永久删除该项, 是否继续?',
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: () => deleteEntry(row)
    })
}

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const deleteEntry = async (row: any) => {
    props.model.formData = row
    const res = await api('exec', {
        projectId: props.projectId,
        sql: props.model.delete(),
        sqlType: 'delete',
        authorization: localStorage.getItem('authorization'),
    })
    if (res.code === 200) {
        message.success('删除成功')
        fetchData()
    } else {
        message.error('删除失败')
    }
}

const submitForm = async (action: string) => {
    isFormVisible.value = false
    props.model.formData = formData.value
    const res = await api('exec', {
        projectId: props.projectId,
        sql: action === 'insert' ? props.model.insert() : props.model.update(),
        sqlType: action,
        authorization: localStorage.getItem('authorization'),
    })
    if (res.code !== 200) {
        message.error(action === 'insert' ? '添加失败' : '编辑失败')
        isFormVisible.value = true
        return
    }
    message.success(action === 'insert' ? '添加成功' : '编辑成功')
    fetchData()
}

const handleClick = async (field: string, formData: any, fieldsType: any) => {
    if (fieldsType['button']['action'] === 'generateLicenseKey') {
        const license_key = encryptAes(formData.expire_time)
        formData.license_key = license_key
    }
}

const handleChange = async (
    field: string,
    formData: any,
    fieldsType: any
) => {
    if (fieldsType['change'] && fieldsType['change'].length) {
        for (const action of fieldsType['change']) {
            if (action === 'generateLicenseKey') {
                const license_key = encryptAes(formData[field])
                formData.license_key = license_key
            }
        }
    }
}

// 从项目信息设置默认端口
const setDefaultPortsFromProject = async () => {
    try {
        const projectInfo = await api('project_info', {
            projectId: props.projectId,
        })
        
        // 设置客户ID为项目ID
        formData.value.client_id = props.projectId
        
        if (projectInfo && projectInfo.api_port && projectInfo.front_port) {
            formData.value.api_port = projectInfo.api_port
            formData.value.front_port = projectInfo.front_port
            console.log('Set default values from project:', {
                client_id: props.projectId,
                api_port: projectInfo.api_port,
                front_port: projectInfo.front_port
            })
        } else {
            // 如果项目没有端口信息，使用默认值
            formData.value.api_port = '9000'
            formData.value.front_port = '3000'
            console.log('Project has no port info, using default values:', {
                client_id: props.projectId,
                api_port: '9000',
                front_port: '3000'
            })
        }
    } catch (error) {
        console.error('Failed to get project info for default ports:', error)
        // 出错时使用默认值
        formData.value.client_id = props.projectId
        formData.value.api_port = '9000'
        formData.value.front_port = '3000'
        message.warning('获取项目端口信息失败，使用默认值')
    }
}

// 获取字段占位符
const getFieldPlaceholder = (field: string) => {
    if (props.model.constructor.name === 'Client') {
        if (field === 'client_id') {
            return isEditMode.value ? '请输入客户ID' : '默认使用项目ID'
        }
        if (field === 'api_port') {
            return isEditMode.value ? '请输入API端口' : '默认使用项目API端口'
        }
        if (field === 'front_port') {
            return isEditMode.value ? '请输入前端端口' : '默认使用项目前端端口'
        }
    }
    return `请输入${field}`
}

</script>

<style scoped>
.search-box {
    margin-bottom: 20px;
}

/* 特殊表格的24px按钮样式 */
:deep(.special-table .special-table-btn) {
    width: 24px !important;
    height: 24px !important;
    min-width: 24px !important;
    max-width: 24px !important;
    padding: 0 !important;
    border-radius: 4px !important;
    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;
}

:deep(.special-table .special-table-btn .n-button__content) {
    padding: 0 !important;
    margin: 0 !important;
    width: 24px !important;
    height: 24px !important;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
}

:deep(.special-table .special-table-btn .n-icon) {
    font-size: 16px !important;
    width: 16px !important;
    height: 16px !important;
    margin: 0 !important;
}

/* 弹框样式优化 */
:deep(.n-modal .n-dialog) {
    width: 600px !important;
    max-width: 90vw !important;
    min-width: 500px !important;
}

:deep(.n-modal .n-dialog .n-dialog__content) {
    padding: 20px 24px !important;
}

:deep(.n-modal .n-form) {
    width: 100%;
}

:deep(.n-modal .n-form-item) {
    margin-bottom: 16px !important;
}

:deep(.n-modal .n-form-item-label) {
    width: 140px !important;
    min-width: 140px !important;
    text-align: right !important;
    padding-right: 12px !important;
}

:deep(.n-modal .n-form-item-blank) {
    flex: 1 !important;
    min-width: 0 !important;
}

/* 输入框组合样式 */
:deep(.n-modal .n-input-group) {
    width: 100% !important;
}

:deep(.n-modal .n-input-group .n-input) {
    flex: 1 !important;
}

/* 单选按钮组样式 */
:deep(.n-modal .n-radio-group) {
    width: 100% !important;
}

:deep(.n-modal .n-radio-button) {
    margin-right: 8px !important;
}

/* 日期选择器样式 */
:deep(.n-modal .n-date-picker) {
    width: 100% !important;
}

/* 响应式弹框 */
@media (max-width: 768px) {
    :deep(.n-modal .n-dialog) {
        width: 95vw !important;
        min-width: auto !important;
        margin: 10px !important;
    }

    :deep(.n-modal .n-form-item-label) {
        width: 100px !important;
        min-width: 100px !important;
        font-size: 12px !important;
    }

    :deep(.n-modal .n-dialog .n-dialog__content) {
        padding: 16px 20px !important;
    }
}
</style>

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
                <!-- 一键填写按钮 - 只在添加客户套餐时显示 -->
                <n-form-item v-if="!isEditMode && isClientModel()">
                    <template #label>
                        <span style="color: #666;">快速填写</span>
                    </template>
                    <n-button type="info" @click="autoFillClientInfo" size="small">
                        <template #icon>
                            <n-icon>
                                <CheckmarkOutline />
                            </n-icon>
                        </template>
                        一键填写项目信息
                    </n-button>
                    <n-text depth="3" style="margin-left: 12px; font-size: 12px;">
                        自动填写客户ID和端口信息
                    </n-text>
                </n-form-item>
                
                <n-form-item v-for="field in fields" :key="field" :label="field" v-if="shouldShowField(field)" required>
                    <template v-if="fieldsType[field]['type'] === 'enum'">
                        <n-radio-group v-model:value="formData[field]">
                            <n-radio-button v-for="item in fieldsType[field]['value']" :key="item" :value="item">
                                {{ item }}
                            </n-radio-button>
                        </n-radio-group>
                    </template>
                    <template v-else-if="fieldsType[field]['type'] === 'datetime'">
                        <n-date-picker 
                            v-model:value="formData[field]" 
                            type="datetime" 
                            placeholder="选择日期时间"
                            value-format="yyyy-MM-dd HH:mm:ss" 
                            format="yyyy-MM-dd HH:mm:ss" 
                            style="width: 100%"
                            :is-date-disabled="() => false"
                            :is-time-disabled="() => false"
                            @update:value="(value) => handleDateTimeChange(field, value)" 
                        />
                    </template>
                    <template v-else>
                        <n-input-group v-if="fieldsType[field]['button']">
                            <n-input v-model:value="formData[field]"
                                :disabled="(isPrimaryKey(field) && isEditMode) || fieldsType[field]['disabled']" />
                            <n-button type="primary" @click="() => handleClick(field)">
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
import { ref, onMounted, defineProps, computed, h, inject } from 'vue'
import { useMessage, useDialog, NButton, NIcon, NSpace, NTooltip } from 'naive-ui'
import { AddCircleOutline, RefreshOutline, CreateOutline, TrashOutline, CloseOutline, CheckmarkOutline } from '@vicons/ionicons5'
import ColorfulButton from '@/components/ColorfulButton.vue'
import api from '@/api'
import { encryptAes } from '@/utils'
import dataManager from '@/utils/dataManager'

const props = defineProps({
    model: { type: Object, required: true },
    projectId: { type: String, required: true },
})

const message = useMessage()
const dialog = useDialog()

// 注入全局 loading
const globalLoading = inject('globalLoading') as any

// 时间格式化函数
const formatDateTime = (value: any): string => {
    if (!value) return ''
    
    let date: Date
    
    // 处理不同的输入格式
    if (typeof value === 'number') {
        // 时间戳
        date = new Date(value)
    } else if (typeof value === 'string') {
        // 字符串格式
        if (value.includes('T')) {
            // ISO格式
            date = new Date(value)
        } else if (value.includes(' ')) {
            // 已经是 YYYY-MM-DD HH:mm:ss 格式
            return value
        } else {
            // 其他字符串格式
            date = new Date(value)
        }
    } else if (value instanceof Date) {
        // Date对象
        date = value
    } else {
        return String(value)
    }
    
    // 检查日期是否有效
    if (isNaN(date.getTime())) {
        return String(value)
    }
    
    // 格式化为 YYYY-MM-DD HH:mm:ss
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

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
        },
        render: (row) => {
            const value = row[field]
            // 如果是时间字段，格式化显示
            if (fieldsType[field] && fieldsType[field].type === 'datetime' && value) {
                return formatDateTime(value)
            }
            return value
        }
    }))

    return [actionColumn, ...fieldColumns]
})

const fetchData = async () => {
    globalLoading.show('正在加载数据...')
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
    } finally {
        globalLoading.hide()
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

// 判断是否为客户套餐管理
const isClientModel = () => {
    const result = props.model.constructor.name === 'Client' || 
                   props.model.tableName === 'tb_client' ||
                   (props.model.fields && props.model.fields.includes('client_id'))
    
    // 添加调试信息
    console.log('isClientModel check:', {
        constructorName: props.model.constructor.name,
        tableName: props.model.tableName,
        hasClientIdField: props.model.fields && props.model.fields.includes('client_id'),
        result: result
    })
    
    return result
}

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const isPrimaryKey = (field: any) => field === primaryKey
// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const shouldShowField = (field: any) => isEditMode.value || field !== primaryKey

const openForm = async (editMode: boolean, row = {}) => {
    console.log('openForm called:', {
        editMode: editMode,
        isClientModel: isClientModel(),
        modelInfo: {
            constructorName: props.model.constructor.name,
            tableName: props.model.tableName,
            fields: props.model.fields
        }
    })
    
    // 如果是添加模式，检查是否达到最大记录数限制
    if (!editMode && isMaxRecordsReached.value) {
        message.warning(maxRecordsMessage.value)
        return
    }

    isEditMode.value = editMode
    
    // 复制行数据并格式化时间字段
    formData.value = { ...row }
    
    // 如果是编辑模式，确保时间字段格式正确
    if (editMode) {
        for (const field in fieldsType) {
            if (fieldsType[field].type === 'datetime' && formData.value[field]) {
                formData.value[field] = formatDateTime(formData.value[field])
            }
        }
    }
    if (!editMode) {
        // 如果是套餐管理且是添加模式，先获取项目信息设置默认端口
        if (isClientModel()) {
            await setDefaultPortsFromProject()
        }
        
        for (const field in fieldsType) {
            const type = fieldsType[field].type
            if (type === 'int' || type === 'string') {
                // 如果是端口字段且已经设置了默认值，跳过
                if ((field === 'api_port' || field === 'front_port') && formData.value[field]) {
                    continue
                }
                // 如果是client_id字段，设置为项目ID
                if (field === 'client_id' && isClientModel()) {
                    formData.value[field] = props.projectId
                } else {
                    formData.value[field] = fieldsType[field].value
                }
            }
            if (type === 'enum') {
                formData.value[field] =
                    fieldsType[field].default || fieldsType[field].value[0]
            }
            if (type === 'datetime') {
                // 默认值为当前时间，格式化为 YYYY-MM-DD HH:mm:ss
                const now = new Date()
                formData.value[field] = formatDateTime(now)
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
    globalLoading.show('正在删除数据...')
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
    globalLoading.hide()
}

const submitForm = async (action: string) => {
    globalLoading.show(action === 'insert' ? '正在添加数据...' : '正在更新数据...')
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
        globalLoading.hide()
        return
    }
    message.success(action === 'insert' ? '添加成功' : '编辑成功')
    fetchData()
    globalLoading.hide()
}

const handleClick = async (field: string) => {
    console.log('handleClick called for field:', field)
    
    const fieldConfig = fieldsType[field]
    if (fieldConfig && fieldConfig['button'] && fieldConfig['button']['action'] === 'generateLicenseKey') {
        console.log('Generating license key via button click, expire_time:', formData.value.expire_time)
        
        if (formData.value.expire_time) {
            try {
                // 使用统一的时间格式化函数
                const expireTimeStr = formatDateTime(formData.value.expire_time)
                
                console.log('Final expire_time format:', expireTimeStr)
                
                const license_key = encryptAes(expireTimeStr)
                formData.value.license_key = license_key
                
                message.success('License Key生成成功')
            } catch (error) {
                console.error('Error generating license key:', error)
                message.error('生成License Key失败: ' + error.message)
            }
        } else {
            message.warning('请先设置过期时间')
        }
    }
}

// 处理日期时间变化
const handleDateTimeChange = async (field: string, value: any) => {
    console.log('handleDateTimeChange called for field:', field, 'value:', value)
    
    // 确保时间格式统一为 YYYY-MM-DD HH:mm:ss
    if (value) {
        formData.value[field] = formatDateTime(value)
    } else {
        formData.value[field] = ''
    }
    
    // 调用原有的变化处理逻辑
    await handleChange(field)
}

const handleChange = async (field: string) => {
    console.log('handleChange called for field:', field)
    
    const fieldConfig = fieldsType[field]
    if (fieldConfig && fieldConfig['change'] && fieldConfig['change'].length) {
        for (const action of fieldConfig['change']) {
            if (action === 'generateLicenseKey') {
                console.log('Generating license key for field:', field, 'value:', formData.value[field])
                if (formData.value[field]) {
                    try {
                        // 使用统一的时间格式化函数
                        const valueStr = formatDateTime(formData.value[field])
                        
                        console.log('Final field value format:', valueStr)
                        
                        const license_key = encryptAes(valueStr)
                        formData.value.license_key = license_key
                        console.log('Generated license key:', license_key)
                    } catch (error) {
                        console.error('Error generating license key:', error)
                        message.error('生成License Key失败: ' + error.message)
                    }
                }
            }
        }
    }
}

// 一键填写客户信息 - 只使用本地缓存数据
const autoFillClientInfo = () => {
    try {
        // 设置基本默认值
        formData.value.client_id = props.projectId
        formData.value.api_port = '9000'
        formData.value.front_port = '3000'
        
        console.log('一键填写：设置基本默认值', {
            client_id: formData.value.client_id,
            api_port: formData.value.api_port,
            front_port: formData.value.front_port
        })
        
        // 从本地缓存获取项目信息（不进行API调用）
        const cachedServers = dataManager.getCachedServerData()
        let projectInfo = null
        
        // 在缓存的服务器数据中查找项目信息
        for (const server of cachedServers) {
            if (server.project_list && Array.isArray(server.project_list)) {
                const project = server.project_list.find((p: any) => p.project_id === props.projectId)
                if (project) {
                    projectInfo = project
                    break
                }
            }
        }
        
        console.log('从缓存获取的项目信息:', projectInfo)
        
        if (projectInfo) {
            // 使用缓存中的端口信息
            if (projectInfo.api_port && String(projectInfo.api_port).trim() !== '') {
                formData.value.api_port = String(projectInfo.api_port).trim()
                console.log('一键填写：从缓存更新API端口', formData.value.api_port)
            }
            
            if (projectInfo.front_port && String(projectInfo.front_port).trim() !== '') {
                formData.value.front_port = String(projectInfo.front_port).trim()
                console.log('一键填写：从缓存更新前端端口', formData.value.front_port)
            }
            
            message.success('项目信息填写完成！')
        } else {
            console.log('缓存中未找到项目信息，使用默认配置')
            message.info('使用默认端口配置')
        }
        
        console.log('一键填写完成:', {
            client_id: formData.value.client_id,
            api_port: formData.value.api_port,
            front_port: formData.value.front_port
        })
        
    } catch (error) {
        console.error('一键填写失败:', error)
        message.error('填写失败，请手动输入')
    }
}

// 从项目信息设置默认端口
const setDefaultPortsFromProject = async () => {
    try {
        // 优先从缓存获取项目信息
        const projectInfo = await dataManager.getProjectById(props.projectId)
        
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
    if (isClientModel()) {
        if (field === 'client_id') {
            return isEditMode.value ? '请输入客户ID' : '点击上方"一键填写"按钮自动填写'
        }
        if (field === 'api_port') {
            return isEditMode.value ? '请输入API端口' : '点击上方"一键填写"按钮自动填写'
        }
        if (field === 'front_port') {
            return isEditMode.value ? '请输入前端端口' : '点击上方"一键填写"按钮自动填写'
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

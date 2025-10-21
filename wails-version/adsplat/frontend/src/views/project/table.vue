<template>
    <div>
        <div class="search-box">
            <n-space>
                <n-tooltip>
                    <template #trigger>
                        <n-button type="primary" @click="openForm(false)">
                            <template #icon>
                                <n-icon><AddCircleOutline /></n-icon>
                            </template>
                            添加
                        </n-button>
                    </template>
                    添加
                </n-tooltip>
                <n-tooltip>
                    <template #trigger>
                        <n-button type="info" @click="refreshData">
                            <template #icon>
                                <n-icon><RefreshOutline /></n-icon>
                            </template>
                            刷新
                        </n-button>
                    </template>
                    刷新
                </n-tooltip>
            </n-space>
        </div>
        
        <n-data-table
            :columns="columns"
            :data="tableData"
            :pagination="pagination"
            striped
            class="special-table"
        />

        <n-modal v-model:show="isFormVisible" preset="dialog" :title="isEditMode ? '编辑' : '添加'">
            <n-form :model="formData" label-placement="left" label-width="180">
                <n-form-item
                    v-for="field in fields"
                    :key="field"
                    :label="field"
                    v-if="shouldShowField(field)"
                    required
                >
                    <template v-if="fieldsType[field]['type'] === 'enum'">
                        <n-radio-group v-model:value="formData[field]">
                            <n-radio-button
                                v-for="item in fieldsType[field]['value']"
                                :key="item"
                                :value="item"
                            >
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
                            @update:value="handleChange(field, formData, fieldsType[field])"
                        />
                    </template>
                    <template v-else>
                        <n-input-group v-if="fieldsType[field]['button']">
                            <n-input
                                v-model:value="formData[field]"
                                :disabled="(isPrimaryKey(field) && isEditMode) || fieldsType[field]['disabled']"
                            />
                            <n-button
                                type="primary"
                                @click="handleClick(field, formData, fieldsType[field])"
                            >
                                {{ fieldsType[field]['button'].text }}
                            </n-button>
                        </n-input-group>
                        <n-input
                            v-else
                            v-model:value="formData[field]"
                            :disabled="(isPrimaryKey(field) && isEditMode) || fieldsType[field]['disabled']"
                        />
                    </template>
                </n-form-item>
            </n-form>
            <template #action>
                <n-space>
                    <n-button @click="isFormVisible = false">
                        <template #icon>
                            <n-icon><CloseOutline /></n-icon>
                        </template>
                        取消
                    </n-button>
                    <n-button
                        type="primary"
                        @click="isEditMode ? submitForm('update') : submitForm('insert')"
                    >
                        <template #icon>
                            <n-icon><CheckmarkOutline /></n-icon>
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

// 分页配置
const pagination = {
    pageSize: 10,
    showSizePicker: true,
    pageSizes: [10, 20, 50],
    showQuickJumper: true
}

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
                        trigger: () => h(NButton, {
                            size: 'small',
                            type: 'primary',
                            class: 'special-table-btn',
                            onClick: () => openForm(true, row)
                        }, {
                            icon: () => h(NIcon, { size: 16 }, { default: () => h(CreateOutline) })
                        }),
                        default: () => '编辑'
                    }),
                    h(NTooltip, { trigger: 'hover' }, {
                        trigger: () => h(NButton, {
                            size: 'small',
                            type: 'error',
                            class: 'special-table-btn',
                            onClick: () => confirmDelete(row)
                        }, {
                            icon: () => h(NIcon, { size: 16 }, { default: () => h(TrashOutline) })
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
    const res = await api('exec', {
        projectId: props.projectId,
        sql: props.model.selects(),
        sqlType: 'selects',
        authorization: localStorage.getItem('authorization'),
    })
    if (res.code === 200) {
        tableData.value = res.data.result
    } else {
        message.error('获取数据失败')
    }
}
onMounted(fetchData)

const refreshData = () => fetchData()

// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const isPrimaryKey = (field: any) => field === primaryKey
// biome-ignore lint/suspicious/noExplicitAny: <explanation>
const shouldShowField = (field: any) => isEditMode.value || field !== primaryKey

const openForm = (editMode: boolean, row = {}) => {
    isEditMode.value = editMode
    formData.value = { ...row }
    if (!editMode) {
        for (const field in fieldsType) {
            const type = fieldsType[field].type
            if (type === 'int' || type === 'string') {
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
</style>

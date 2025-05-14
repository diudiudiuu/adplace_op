<template>
    <div>
        <div class="search-box">
            <el-button text bg type="primary" @click="openForm(false)">
                <i class="el-icon-lx-roundaddfill"></i> 添加
            </el-button>
            <el-button text bg type="success" @click="refreshData">
                <i class="el-icon-lx-refresh"></i> 刷新
            </el-button>
        </div>
        <el-table
            :data="tableData"
            class="table"
            ref="multipleTable"
            header-cell-class-name="table-header"
            stripe
        >
            <el-table-column fixed label="操作" align="left" min-width="120px">
                <template #default="scope">
                    <el-tooltip content="编辑" placement="top-start" :hide-after="0">
                        <el-button
                            text
                            bg
                            type="primary"
                            size="small"
                            @click="openForm(true, scope.row)"
                        >
                            <i class="el-icon-lx-edit"></i>
                        </el-button>
                    </el-tooltip>
                    <el-tooltip content="删除" placement="top-start" :hide-after="0">
                        <el-button
                            text
                            bg
                            type="danger"
                            size="small"
                            @click="confirmDelete(scope.row)"
                        >
                            <i class="el-icon-lx-deletefill"></i>
                        </el-button>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column
                v-for="field in fields"
                :key="field"
                :prop="field"
                :label="field"
                align="left"
                width="auto"
                show-overflow-tooltip
            />
        </el-table>

        <el-dialog :title="isEditMode ? '编辑' : '添加'" v-model="isFormVisible" width="65%">
            <el-form :model="formData" label-width="180px">
                <el-form-item
                    v-for="field in fields"
                    :key="field"
                    :label="field"
                    v-if="shouldShowField(field)"
                    required
                >
                    <template v-if="fieldsType[field]['type'] === 'enum'">
                        <el-radio-group v-model="formData[field]">
                            <el-radio-button
                                v-for="item in fieldsType[field]['value']"
                                :key="item"
                                :label="item"
                                :value="item"
                            />
                        </el-radio-group>
                    </template>
                    <template v-else-if="fieldsType[field]['type'] === 'datetime'">
                        <el-date-picker
                            v-model="formData[field]"
                            type="datetime"
                            placeholder="选择日期时间"
                            value-format="YYYY-MM-DD HH:mm:ss"
                            format="YYYY-MM-DD HH:mm:ss"
                            style="width: 100%"
                            @change="handleChange(field, formData, fieldsType[field])"
                        />
                    </template>
                    <template v-else>
                        <el-input
                            v-model="formData[field]"
                            :disabled="(isPrimaryKey(field) && isEditMode) || fieldsType[field]['disabled']"
                        >
                            <template v-if="fieldsType[field]['button']" #append>
                                <el-button
                                    @click="handleClick(field, formData, fieldsType[field])"
                                >{{ fieldsType[field]['button'].text }}</el-button>
                            </template>
                        </el-input>
                    </template>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="isFormVisible = false">取消</el-button>
                <el-button
                    type="primary"
                    @click="isEditMode ? submitForm('update') : submitForm('insert')"
                >提交</el-button>
            </div>
        </el-dialog>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, defineProps } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'
import { encryptAes } from '@/utils'

const props = defineProps({
    model: { type: Object, required: true },
    projectId: { type: String, required: true },
})

const tableData = ref([])
const formData = ref({})
const isFormVisible = ref(false)
const isEditMode = ref(false)

const fields = props.model.fields
const fieldsType = props.model.fieldsType
const primaryKey = props.model.primaryKey

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
        ElMessage.error('获取数据失败')
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
    ElMessageBox.confirm('此操作将永久删除该项, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
    })
        .then(() => deleteEntry(row))
        .catch(() => {})
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
        ElMessage.success('删除成功')
        fetchData()
    } else {
        ElMessage.error('删除失败')
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
        ElMessage.error(action === 'insert' ? '添加失败' : '编辑失败')
        isFormVisible.value = true
        return
    }
    ElMessage.success(action === 'insert' ? '添加成功' : '编辑成功')
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
.table {
    width: 100%;
}
</style>

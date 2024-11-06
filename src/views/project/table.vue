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
            <el-table-column label="操作" align="left">
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
            />
        </el-table>

        <el-dialog :title="isEditMode ? '编辑条目' : '添加条目'" v-model="isFormVisible" width="50%">
            <el-form :model="formData" label-width="100px">
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
                            style="width: 100%"
                        />
                    </template>
                    <template v-else>
                        <el-input
                            v-model="formData[field]"
                            :disabled="isPrimaryKey(field) && isEditMode"
                        />
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
    })
    if (res.code === 200) {
        ElMessage.success('删除成功')
        fetchData()
    } else {
        ElMessage.error('删除失败')
    }
}

const submitForm = async (action: string) => {
    props.model.formData = formData.value
    const res = await api('exec', {
        projectId: props.projectId,
        sql: action === 'insert' ? props.model.insert() : props.model.update(),
        sqlType: action,
    })
    if (res.code !== 200) {
        ElMessage.error(action === 'insert' ? '添加失败' : '编辑失败')
        return
    }
    ElMessage.success(action === 'insert' ? '添加成功' : '编辑成功')
    isFormVisible.value = false
    fetchData()
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

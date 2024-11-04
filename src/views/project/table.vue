<template>
    <div>
        <div class="search-box">
            <el-button text bg type="primary" @click="openForm(false)">
                <i class="el-icon-lx-roundaddfill"></i> 新增
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
            <el-table-column label="操作" fixed align="left">
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
                align="center"
            />
        </el-table>

        <el-dialog :title="isEditMode ? '编辑条目' : '新增条目'" v-model="isFormVisible" width="50%">
            <el-form :model="formData" label-width="100px">
                <el-form-item
                    v-for="field in fields"
                    :key="field"
                    :label="field"
                    v-if="shouldShowField(field)"
                >
                    <template v-if="fieldsType[field]['type'] === 'enum'">
                        <el-select v-model="formData[field]">
                            <el-option
                                v-for="item in fieldsType[field]['value']"
                                :key="item"
                                :label="item"
                                :value="item"
                            />
                        </el-select>
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
console.log('fields', fields)
console.log('fieldsType', fieldsType)
const primaryKey = props.model.primaryKey

const fetchData = async () => {
    const res = await api('exec', {
        projectId: props.projectId,
        sql: props.model.select_list(),
        sqlType: 'select_list',
    })
    tableData.value = res.data.result
}
onMounted(fetchData)

const refreshData = () => fetchData()

const isPrimaryKey = (field) => field === primaryKey
const shouldShowField = (field) => isEditMode.value || field !== primaryKey

const openForm = (editMode, row = {}) => {
    isEditMode.value = editMode
    formData.value = { ...row }
    props.model.formData = formData.value
    isFormVisible.value = true
}

const confirmDelete = (row) => {
    ElMessageBox.confirm('此操作将永久删除该项, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
    })
        .then(() => deleteEntry(row))
        .catch(() => {})
}

const deleteEntry = async (row) => {
    props.model.formData = row
    await api('exec', {
        projectId: props.projectId,
        sql: props.model.delete(),
        sqlType: 'delete',
    })
    ElMessage.success('删除成功')
    fetchData()
}

const submitForm = async (action) => {
    props.model.formData = formData.value
    await api('exec', {
        projectId: props.projectId,
        sql: action === 'insert' ? props.model.insert() : props.model.update(),
        sqlType: action,
    })
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

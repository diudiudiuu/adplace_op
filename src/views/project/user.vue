<template>
    <div>
        <div class="search-box">
            <el-button text bg type="primary" @click="openAddForm">
                <i class="el-icon-lx-roundaddfill"></i> 新增
            </el-button>
            <el-button text bg type="success" @click="handleRefresh">
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
            <el-table-column label="操作" fixed width="240" align="left">
                <template #default="scope">
                    <el-tooltip content="编辑" placement="top-start" :hide-after="0">
                        <el-button
                            text
                            bg
                            type="primary"
                            size="small"
                            @click="openEditForm(scope.row)"
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
                            @click="handleDelete(scope.row)"
                        >
                            <i class="el-icon-lx-deletefill"></i>
                        </el-button>
                    </el-tooltip>
                </template>
            </el-table-column>
            <el-table-column
                v-for="col in columns"
                :key="col.prop"
                :prop="col.prop"
                :label="col.label"
                align="center"
                :width="colWidth"
            ></el-table-column>
        </el-table>

        <!-- Add/Edit Form Modal -->
        <el-dialog :title="isEditMode ? '编辑条目' : '新增条目'" :visible.sync="isFormVisible" width="50%">
            <el-form :model="formData" label-width="100px">
                <el-form-item v-for="field in columns" :key="field.prop" :label="field.label">
                    <el-input v-model="formData[field.prop]" />
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="isFormVisible = false">取消</el-button>
                <el-button
                    type="primary"
                    @click="isEditMode ? submitEditForm() : submitAddForm()"
                >提交</el-button>
            </div>
        </el-dialog>
    </div>
</template>
  
  <script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { defineProps } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

const props = defineProps({
    projectId: {
        type: String,
        required: true,
    },
})

const tableData = ref([])
const columns = ref([])
const formData = ref({})
const isFormVisible = ref(false)
const isEditMode = ref(false)
const colWidth = 'auto' // Set column width to auto

// Fetch data and set up columns dynamically
const getDataList = () => {
    api('exec', {
        projectId: props.projectId,
        sql: 'select * from tb_auth_user',
        sqlType: 'select_list',
    }).then((res: any) => {
        tableData.value = res.data.result
        if (tableData.value.length > 0) {
            columns.value = Object.keys(tableData.value[0]).map((key) => ({
                prop: key,
                label: key, // Keep labels as they are in JSON response
            }))
        }
    })
}
onMounted(getDataList)

const handleRefresh = () => {
    getDataList()
}

// Delete item function
const handleDelete = (row) => {
    ElMessageBox.confirm('此操作将永久删除该项, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
    })
        .then(() => {
            api('exec', {
                projectId: props.projectId,
                sql: `delete from tb_data where id = ${row.id}`,
                sqlType: 'delete',
            }).then(() => {
                ElMessage.success('删除成功')
                getDataList()
            })
        })
        .catch(() => {})
}

// Open the form in add or edit mode
const openAddForm = () => {
    formData.value = {}
    isEditMode.value = false
    isFormVisible.value = true
}

const openEditForm = (row) => {
    formData.value = { ...row }
    isEditMode.value = true
    isFormVisible.value = true
}

// Submit add form
const submitAddForm = () => {
    // Call API to add a new item (logic to be implemented as needed)
    ElMessage.success('添加成功')
    isFormVisible.value = false
    getDataList()
}

// Submit edit form
const submitEditForm = () => {
    // Call API to edit the item (logic to be implemented as needed)
    ElMessage.success('编辑成功')
    isFormVisible.value = false
    getDataList()
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
  
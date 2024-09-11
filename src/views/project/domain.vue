<template>
    <div>
        <div class="search-box">
            <el-button text bg type="primary">
                <i class="el-icon-lx-roundaddfill"></i> 新增
            </el-button>
            <el-button text bg type="success" @click="handleReflash">
                <i class="el-icon-lx-refresh"></i>刷新
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
                        <el-button text bg type="primary" size="small">
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
            <el-table-column prop="id" label="ID" width="80" align="center"></el-table-column>
            <el-table-column prop="domain" label="域名" align="center"></el-table-column>
            <el-table-column prop="created_at" label="注册时间" align="center"></el-table-column>
            <el-table-column prop="expired_at" label="到期时间" align="center"></el-table-column>
            <el-table-column prop="ping_flag" label="ping状态" align="center">
                <template #default="{ row }">
                    <el-tag v-if="row.ping_flag == 'normal'" type="success">正常</el-tag>
                    <el-tag v-else type="danger">异常</el-tag>&nbsp;
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import { defineProps } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/api'

// 定义接受 projectId 的 props
const props = defineProps({
    projectId: {
        type: String,
        required: true,
    },
})
// 定义域名列表
const tableData = ref([])
const getDomainList = () => {
    api('exec', {
        projectId: props.projectId,
        sql: 'select * from tb_domain',
        sqlType: 'select_list',
    }).then((res: any) => {
        tableData.value = res.data.result
    })
}
getDomainList()

const handleReflash = () => {
    getDomainList()
}

const delDomain = (id: string) => {
    api('exec', {
        projectId: props.projectId,
        sql: `delete from tb_domain where id = ${id}`,
        sqlType: 'delete',
    }).then((res: any) => {
        ElMessage.success('删除成功')
        getDomainList()
    })
}

const handleDelete = (row): void => {
    // 二次确认
    ElMessageBox.confirm('此操作将永久删除该域名, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
    })
        .then(() => {
            delDomain(row.id)
        })
        .catch(() => {})
}
</script>

<style scoped>
.search-box {
    margin-bottom: 20px;
}

.search-input {
    width: 200px;
}

.mr10 {
    margin-right: 10px;
}
</style>
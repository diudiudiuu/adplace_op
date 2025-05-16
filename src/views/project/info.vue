<template>
    <div>
        <el-descriptions :column="1" direction="vertical" border>
            <el-descriptions-item>
                <template #label>操作</template>
                <el-button text bg type="primary" @click="handleEdit">{{ !eidtmode ? '编辑' : '取消' }}</el-button>
                <el-button text bg type="danger" @click="handleDelete">删除</el-button>
                <Dform
                    v-if="eidtmode"
                    mode="edit"
                    :serverId="serverId"
                    :initialForm="projectInfo"
                    @editSuccess="updateHandle"
                />
            </el-descriptions-item>
            <template v-if="!eidtmode">
                <el-descriptions-item>
                    <template #label>客户ID</template>
                    {{ projectInfo.project_id }}
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>客户名称</template>
                    {{ projectInfo.project_name }}
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>客户管理地址</template>
                    {{ projectInfo.project_manage_url }}
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>客户API地址</template>
                    {{ projectInfo.project_api_url }}
                </el-descriptions-item>
            </template>
        </el-descriptions>
    </div>
</template>
<script lang="ts" setup>
import { ref, defineProps } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useSidebarStore } from '@/store/sidebar'
import Dform from './form.vue'
import api from '@/api'

const sidebar = useSidebarStore()
// 定义接受 projectId 的 props
const props = defineProps({
    serverId: {
        type: String,
        required: true,
    },
    projectId: {
        type: String,
        required: true,
    },
})

const projectInfo = ref({})

const eidtmode = ref(false)
// 编辑按钮点击事件
const handleEdit = () => {
    // 切换到编辑模式
    eidtmode.value = !eidtmode.value
}
// 获取项目信息
const getProjectInfo = () => {
    api('project_info', {
        projectId: props.projectId,
        authorization: localStorage.getItem('authorization'),
    }).then((res: any) => {
        projectInfo.value = res
    })
}

getProjectInfo()
const updateHandle = () => {
    eidtmode.value = false
    getProjectInfo()
}

// 删除按钮点击事件
const handleDelete = () => {
    ElMessageBox.confirm('是否删除该客户吗?', '提示', {
        type: 'warning',
    })
        .then(() => {
            api('project_delete', {
                serverId: props.serverId,
                projectId: props.projectId,
                authorization: localStorage.getItem('authorization'),
            }).then((res: any) => {
                ElMessage.success('删除成功')
                sidebar.setboolroute(true)
            })
        })
        .catch(() => {
            ElMessage.error('删除失败')
        })
}
</script>
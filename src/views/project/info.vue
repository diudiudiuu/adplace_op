<template>
    <div>
        <el-descriptions :column="1" direction="vertical" border>
            <el-descriptions-item>
                <template #label>操作</template>
                <el-button text bg type="primary" @click="handleEdit">{{ !eidtmode ? '编辑' : '取消' }}</el-button>
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
                    <template #label>合同日期</template>
                    {{ projectInfo.contract_date }}
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>客户管理地址</template>
                    {{ projectInfo.project_manage_url }}
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>客户API地址</template>
                    {{ projectInfo.project_api_url }}
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>API端口</template>
                    {{ projectInfo.api_port }}
                </el-descriptions-item>
                <el-descriptions-item>
                    <template #label>Front端口</template>
                    {{ projectInfo.front_port }}
                </el-descriptions-item>
            </template>
        </el-descriptions>
    </div>
</template>
<script lang="ts" setup>
import { ref, defineProps } from 'vue'
import Dform from './form.vue'
import api from '@/api'

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
        // biome-ignore lint/suspicious/noExplicitAny: <explanation>
    }).then((res: any) => {
        projectInfo.value = res
    })
}

getProjectInfo()
const updateHandle = () => {
    eidtmode.value = false
    getProjectInfo()
}
</script>
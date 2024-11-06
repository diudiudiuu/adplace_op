<template>
    <div>
        <el-descriptions :column="1" direction="vertical" border>
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
                <template #label>管理端链接</template>
                {{ projectInfo.project_manage_url }}
            </el-descriptions-item>
            <el-descriptions-item>
                <template #label>API链接</template>
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
        </el-descriptions>
    </div>
</template>
<script lang="ts" setup>
import { ref, defineProps } from 'vue'
import api from '@/api'

// 定义接受 projectId 的 props
const props = defineProps({
    projectId: {
        type: String,
        required: true,
    },
})

const projectInfo = ref({})

api('project_info', {
    projectId: props.projectId,
    // biome-ignore lint/suspicious/noExplicitAny: <explanation>
}).then((res: any) => {
    console.log(res)
    projectInfo.value = res
})
</script>
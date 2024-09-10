<template>
    <div>
        <el-descriptions :column="1" direction="vertical" border>
            <el-descriptions-item>
                <template #label>名称</template>
                {{ projectInfo.project_name }}
            </el-descriptions-item>
            <el-descriptions-item>
                <template #label>标识</template>
                {{ projectInfo.project_path }}
            </el-descriptions-item>
            <el-descriptions-item>
                <template #label>管理端链接</template>
                {{ projectInfo.project_web_url }}
            </el-descriptions-item>
            <el-descriptions-item>
                <template #label>API链接</template>
                {{ projectInfo.project_api_url }}
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

const projectInfo = ref({
    project_name: '',
    project_path: '',
    project_web_url: '',
    project_api_url: '',
})

api('project_info', {
    projectId: props.projectId,
}).then((res: any) => {
    projectInfo.value = res
})
</script>
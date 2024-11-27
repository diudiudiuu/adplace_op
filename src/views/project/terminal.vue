<template>
    <div class="terminal">
        <pre>{{ rawText }}</pre>
        <el-progress :indeterminate="true" />
        <el-button type="primary">启动</el-button>
        <el-mention
            v-model="value"
            type="textarea"
            :options="options"
            style="width: 320px"
            placeholder="Please input"
        />
    </div>
</template>
  
<script setup>
import { ref, defineProps } from 'vue'

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

const rawText = ref(`
:::demo
<div>水电费都是</div>
<template>
    <el-button type="primary">Primary Button</el-button>
</template>
:::
`)

const value = ref('')
const options = ref([
    {
        label: 'Fuphoenixes',
        value: 'Fuphoenixes',
    },
    {
        label: 'kooriookami',
        value: 'kooriookami',
    },
    {
        label: 'Jeremy',
        value: 'Jeremy',
    },
    {
        label: 'btea',
        value: 'btea',
    },
])
</script>
  
<style scoped>
.terminal >>> pre {
    background-color: #f5f5f5;
    padding: 10px;
    border-radius: 5px;
}
</style>
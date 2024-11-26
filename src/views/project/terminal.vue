<template>
    <div class="terminal">
        <terminal
            name="my-terminal"
            :show-header="false"
            @exec-cmd="onExecCmd"
            :drag-conf="dragConf"
            theme="dark"
        />
    </div>
</template>
  
  <script setup>
import { ref, defineProps } from 'vue'
import Terminal from 'vue-web-terminal'

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

// 设置 dragConf 配置
const dragConf = ref({
    width: '50%',
    height: '70%',
    zIndex: 100,
    pinned: false,
})

// 定义 onExecCmd 方法
const onExecCmd = (key, command, success, failed) => {
    if (key === 'fail') {
        failed('Something wrong!!!')
    } else {
        const allClass = ['success', 'error', 'system', 'info', 'warning']
        const clazz = allClass[Math.floor(Math.random() * allClass.length)]

        success({
            type: 'normal',
            class: clazz,
            tag: clazz,
            content: `Your command is '${command}'`,
        })
    }
}
</script>
  
<style scoped>
</style>
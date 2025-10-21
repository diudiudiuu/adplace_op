<template>
    <div class="dashboard-container">
        <n-tabs v-model:value="activeTab" type="line" placement="left" size="large" class="dashboard-tabs">
            <n-tab-pane name="info" tab="客户信息">
                <Info v-if="activeTab === 'info'" :serverId="serverId" :projectId="projectId" />
            </n-tab-pane>

            <n-tab-pane name="client" tab="套餐管理">
                <table-component v-if="activeTab === 'client'" :model="clientModel" :projectId="projectId" />
            </n-tab-pane>

            <n-tab-pane name="user" tab="用户管理">
                <table-component v-if="activeTab === 'user'" :model="userModel" :projectId="projectId" />
            </n-tab-pane>

            <n-tab-pane name="user_token" tab="Token管理">
                <table-component v-if="activeTab === 'user_token'" :model="userTokenModel" :projectId="projectId" />
            </n-tab-pane>

            <n-tab-pane name="domain" tab="域名管理">
                <table-component v-if="activeTab === 'domain'" :model="domainModel" :projectId="projectId" />
            </n-tab-pane>

            <n-tab-pane name="user_domain" tab="用户域名">
                <table-component v-if="activeTab === 'user_domain'" :model="userDomainModel" :projectId="projectId" />
            </n-tab-pane>

            <n-tab-pane name="terminal" tab="Terminal">
                <Terminal v-if="activeTab === 'terminal'" :serverId="serverId" :projectId="projectId" />
            </n-tab-pane>
        </n-tabs>
    </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'

import Info from './info.vue'
import tableComponent from './table.vue'
import Terminal from './terminal.vue'
import client from '@/model/client'
import auth_user from '@/model/auth_user'
import user_token from '@/model/auth_user_token'
import domain from '@/model/domain'
import user_domain from '@/model/user_domain'

// 移除 tabPosition，Naive UI 使用 placement 属性
const activeTab = ref('info') // Set initial active tab here

const route = useRoute()
const serverId = ref(route.params.pid)
const projectId = ref(route.params.id)

// 更新路由参数的函数
const updateRouteParams = () => {
    serverId.value = route.params.pid
    projectId.value = route.params.id
    
    console.log('Dashboard: Updated route params:', {
        serverId: serverId.value,
        projectId: projectId.value,
        fullParams: route.params
    })
    
    // 验证参数
    if (!serverId.value || !projectId.value) {
        console.error('Dashboard: Missing route parameters!', {
            serverId: serverId.value,
            projectId: projectId.value,
            fullParams: route.params
        })
    }
}

// 监听路由变化
watch(() => route.params, () => {
    console.log('Dashboard: Route params changed')
    updateRouteParams()
}, { immediate: true })

// 组件挂载时更新参数
onMounted(() => {
    console.log('Dashboard: Component mounted')
    updateRouteParams()
})
const clientModel = new client()
const userModel = new auth_user()
const userTokenModel = new user_token()
const domainModel = new domain()
const userDomainModel = new user_domain()
</script>
<style scoped>
.dashboard-container {
    height: 100%;
    background: transparent;
}

.dashboard-tabs {
    height: 100%;
}

:deep(.n-tabs .n-tabs-nav) {
    padding: 8px;
}

:deep(.n-tabs .n-tabs-tab) {
    margin: 4px 0;
    transition: all 0.3s ease;
    font-weight: 500;
}

:deep(.n-tabs .n-tabs-tab:hover) {
    transform: translateX(2px);
}

:deep(.n-tabs .n-tabs-content) {
    padding: 24px;
    margin-left: 16px;
}

:deep(.n-tabs .n-tabs-pane) {
    padding: 0;
}
</style>
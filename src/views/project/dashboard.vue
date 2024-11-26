<template>
    <div>
        <el-tabs :tab-position="tabPosition" v-model="activeTab">
            <el-tab-pane label="客户信息" name="info">
                <Info v-if="activeTab === 'info'" :serverId="serverId" :projectId="projectId" />
            </el-tab-pane>
            <el-tab-pane label="套餐管理" name="client">
                <table-component
                    v-if="activeTab === 'client'"
                    :model="clientModel"
                    :projectId="projectId"
                />
            </el-tab-pane>
            <el-tab-pane label="用户管理" name="user">
                <table-component
                    v-if="activeTab === 'user'"
                    :model="userModel"
                    :projectId="projectId"
                />
            </el-tab-pane>
            <el-tab-pane label="用户Token管理" name="user_token">
                <table-component
                    v-if="activeTab === 'user_token'"
                    :model="userTokenModel"
                    :projectId="projectId"
                />
            </el-tab-pane>
            <el-tab-pane label="域名管理" name="domain">
                <table-component
                    v-if="activeTab === 'domain'"
                    :model="domainModel"
                    :projectId="projectId"
                />
            </el-tab-pane>
            <el-tab-pane label="用户域名" name="user_domain">
                <table-component
                    v-if="activeTab === 'user_domain'"
                    :model="userDomainModel"
                    :projectId="projectId"
                />
            </el-tab-pane>
            <el-tab-pane label="Terminal" name="terminal">
                <Terminal
                    v-if="activeTab === 'terminal'"
                    :serverId="serverId"
                    :projectId="projectId"
                />
            </el-tab-pane>
        </el-tabs>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'

import Info from './info.vue'
import tableComponent from './table.vue'
import Terminal from './terminal.vue'
import client from '@/model/client'
import auth_user from '@/model/auth_user'
import user_token from '@/model/auth_user_token'
import domain from '@/model/domain'
import user_domain from '@/model/user_domain'

const tabPosition = ref('left')
const activeTab = ref('info') // Set initial active tab here

const route = useRoute()
const serverId = ref(route.params.pid)
const projectId = ref(route.params.id)
const clientModel = new client()
const userModel = new auth_user()
const userTokenModel = new user_token()
const domainModel = new domain()
const userDomainModel = new user_domain()
</script>

<template>
    <div>
        <el-tabs :tab-position="tabPosition" v-model="activeTab">
            <el-tab-pane label="客户信息" name="info">
                <Info v-if="activeTab === 'info'" :projectId="projectId" />
            </el-tab-pane>
            <el-tab-pane label="用户管理" name="user">
                <table-component
                    v-if="activeTab === 'user'"
                    :model="userModel"
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
        </el-tabs>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'

import Info from './info.vue'
import tableComponent from './table.vue'
import auth_user from '@/model/auth_user'
import domain from '@/model/domain'
import user_domain from '@/model/user_domain'

const tabPosition = ref('left')
const activeTab = ref('info') // Set initial active tab here

const route = useRoute()
const projectId = ref(route.params.id)
const userModel = new auth_user()
const domainModel = new domain()
const userDomainModel = new user_domain()

// CREATE TABLE "tb_user_domain" (
//   "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
//   "user_id" INTEGER NOT NULL DEFAULT 0,
//   "domain_id" INTEGER NOT NULL DEFAULT 0,
//   "created_at" datetime,
//   "updated_at" datetime
// );
</script>

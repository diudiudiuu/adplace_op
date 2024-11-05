<template>
    <div>
        <el-tabs :tab-position="tabPosition">
            <el-tab-pane label="项目信息">
                <Info :projectId="projectId" />
            </el-tab-pane>
            <el-tab-pane label="域名管理">
                <Domain :projectId="projectId" />
            </el-tab-pane>
            <el-tab-pane label="用户管理">
                <User :projectId="projectId" />
            </el-tab-pane>
            <el-tab-pane label="SQL在线运行">
                <User :projectId="projectId" />
            </el-tab-pane>
        </el-tabs>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'

// 导入 domain.vue 和 info.vue 中的代码
import Domain from './domain.vue'
import Info from './info.vue'
import User from './user.vue'

const tabPosition = ref('left')

const route = useRoute()
const projectId = ref(route.params.id)
</script>
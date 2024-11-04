<template>
    <table-component :model="userModel" :projectId="projectId" />
</template>

<script lang="ts" setup>
import { defineProps } from 'vue'
import tableComponent from './table.vue'
import auth_user from '@/model/auth_user'

// Props for project ID
defineProps({
    projectId: {
        type: String,
        required: true,
    },
})

const userModel = new auth_user()
</script>

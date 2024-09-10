<template>
    <div class="wrapper">
        <v-header />
        <v-sidebar />
        <div class="content-box" :class="{ 'content-collapse': sidebar.collapse }">
            <div class="content">
                <router-view v-slot="{ Component }" :key="$route.fullPath">
                    <transition name="move" mode="out-in">
                        <keep-alive>
                            <component :is="Component"></component>
                        </keep-alive>
                    </transition>
                </router-view>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { useSidebarStore } from '@/store/sidebar'
import vHeader from '@/components/header.vue'
import vSidebar from '@/components/sidebar.vue'

const sidebar = useSidebarStore()
</script>

<style>
.wrapper {
    height: 100vh;
    overflow: hidden;
}
.content-box {
    position: absolute;
    left: 250px;
    right: 0;
    top: 70px;
    bottom: 0;
    padding-bottom: 30px;
    -webkit-transition: left 0.3s ease-in-out;
    transition: left 0.3s ease-in-out;
    background: #ffffff;
    overflow: hidden;
}

.content {
    width: auto;
    height: 100%;
    padding: 20px;
    overflow-y: scroll;
    box-sizing: border-box;
}

.content::-webkit-scrollbar {
    width: 0;
}

.content-collapse {
    left: 65px;
}
</style>

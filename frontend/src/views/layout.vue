<template>
    <div>
        <n-message-provider>
            <n-layout has-sider style="height: 100vh;">
                <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="200"
                    :collapsed="sidebar.collapse" :show-trigger="false" class="custom-sider">
                    <div class="sidebar-container">


                        <div class="menu-container">
                            <n-menu v-if="showMenu" :value="onRoutes" :collapsed="sidebar.collapse"
                                :collapsed-width="64" :options="menuOptions" @update:value="handleMenuSelect"
                                :indent="20" />

                            <!-- 折叠按钮 -->
                            <div class="collapse-button-container">
                                <n-button quaternary circle size="small" @click="sidebar.handleCollapse"
                                    class="collapse-toggle-btn">
                                    <template #icon>
                                        <n-icon size="18">
                                            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                                <path v-if="!sidebar.collapse" d="M15 18L9 12L15 6" stroke="white"
                                                    stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />
                                                <path v-else d="M9 18L15 12L9 6" stroke="white" stroke-width="3"
                                                    stroke-linecap="round" stroke-linejoin="round" />
                                            </svg>
                                        </n-icon>
                                    </template>
                                </n-button>
                            </div>
                        </div>


                    </div>
                </n-layout-sider>

                <n-layout class="main-layout" :class="{ 'layout-loading': isContentLoading }">
                    <n-layout-content content-style="padding: 8px;" class="main-content"
                        :class="{ 'content-loading-parent': isContentLoading }">
                        <div class="content-wrapper" :class="{ 'content-loading': isContentLoading }">
                            <router-view v-slot="{ Component }" :key="$route.fullPath">
                                <transition name="fade" mode="out-in">
                                    <keep-alive>
                                        <component :is="Component" />
                                    </keep-alive>
                                </transition>
                            </router-view>

                            <!-- 主内容区域 Loading 遮罩 -->
                            <div v-if="isContentLoading" class="content-loading-overlay" @wheel.prevent.stop
                                @touchmove.prevent.stop @scroll.prevent.stop @mousewheel.prevent.stop
                                @DOMMouseScroll.prevent.stop>
                                <div class="loading-spinner">
                                    <img :src="loadingGif" alt="Loading..." class="loading-gif" />
                                    <span class="loading-text">{{ loadingText }}</span>
                                </div>
                            </div>
                        </div>
                    </n-layout-content>
                </n-layout>
            </n-layout>
        </n-message-provider>


    </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, h, onMounted, provide } from 'vue'
import { storeToRefs } from 'pinia'
import { useSidebarStore } from '@/store/sidebar'
import { useRoute, useRouter } from 'vue-router'
import { getMenus, reloadMenus } from '@/components/menu'
import { NIcon, useMessage, useDialog, NSpin } from 'naive-ui'
import { setGlobalInstances } from '@/api'
import ColorfulIcons from '@/components/ColorfulIcons.vue'
import loadingGif from '@/assets/img/loading.gif'

const route = useRoute()
const router = useRouter()
const onRoutes = computed(() => route.path)

const sidebar = useSidebarStore()
const { boolroute } = storeToRefs(sidebar)

// 获取 Naive UI 实例
const message = useMessage()

// Loading 状态管理
let loadingStartTime = 0
let originalBodyOverflow = ''
let originalHtmlOverflow = ''

// 阻止滚动的事件处理函数
const preventScroll = (e: Event) => {
    e.preventDefault()
    e.stopPropagation()
    return false
}

// 阻止键盘滚动的事件处理函数
const preventKeyboardScroll = (e: KeyboardEvent) => {
    // 阻止方向键、空格键、Page Up/Down 等滚动按键
    if ([32, 33, 34, 35, 36, 37, 38, 39, 40].includes(e.keyCode)) {
        e.preventDefault()
        e.stopPropagation()
    }
}

// 全局 Loading 控制函数
const showContentLoading = (text: string = '请稍候...') => {
    loadingText.value = text
    isContentLoading.value = true
    loadingStartTime = Date.now()

    // 保存原始样式
    originalBodyOverflow = document.body.style.overflow
    originalHtmlOverflow = document.documentElement.style.overflow

    // 阻止页面滚动 - 多层防护
    document.body.style.overflow = 'hidden'
    document.documentElement.style.overflow = 'hidden'
    document.body.style.position = 'fixed'
    document.body.style.width = '100%'
    document.body.style.height = '100%'
    document.body.style.top = '0'
    document.body.style.left = '0'

    // 锁定视口，防止任何形式的滚动
    const viewport = document.querySelector('meta[name=viewport]')
    if (viewport) {
        viewport.setAttribute('content', 'width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no')
    }

    // 添加事件监听器阻止滚动
    document.addEventListener('wheel', preventScroll, { passive: false })
    document.addEventListener('touchmove', preventScroll, { passive: false })
    document.addEventListener('scroll', preventScroll, { passive: false })
    document.addEventListener('keydown', preventKeyboardScroll, { passive: false })
}

const hideContentLoading = async () => {
    const elapsed = Date.now() - loadingStartTime
    const minDisplayTime = 500 // 最小显示时间 500ms

    if (elapsed < minDisplayTime) {
        // 如果显示时间不足 500ms，等待剩余时间
        await new Promise(resolve => setTimeout(resolve, minDisplayTime - elapsed))
    }

    isContentLoading.value = false

    // 恢复页面滚动
    document.body.style.overflow = originalBodyOverflow
    document.documentElement.style.overflow = originalHtmlOverflow
    document.body.style.position = ''
    document.body.style.width = ''
    document.body.style.height = ''
    document.body.style.top = ''
    document.body.style.left = ''

    // 恢复视口设置
    const viewport = document.querySelector('meta[name=viewport]')
    if (viewport) {
        viewport.setAttribute('content', 'width=device-width, initial-scale=1.0')
    }

    // 移除事件监听器
    document.removeEventListener('wheel', preventScroll)
    document.removeEventListener('touchmove', preventScroll)
    document.removeEventListener('scroll', preventScroll)
    document.removeEventListener('keydown', preventKeyboardScroll)
}

// 创建全局 Loading 实例
const createGlobalLoading = () => {
    return {
        create: (options: any) => {
            const text = options.description || '请稍候...'
            showContentLoading(text)
            return {
                destroy: () => hideContentLoading()
            }
        }
    }
}

// 设置全局实例
onMounted(() => {
    // 设置message和loading实例
    setGlobalInstances(message, createGlobalLoading())
})

// 暴露给子组件使用的方法
const globalLoading = {
    show: showContentLoading,
    hide: hideContentLoading
}

// 提供给子组件
provide('globalLoading', globalLoading)

const menuData = ref<any[]>([])
const menuOptions = ref<any[]>([])
const showMenu = ref(true)

// 主内容区域 Loading 状态
const isContentLoading = ref(false)
const loadingText = ref('请稍候...')

// 渲染彩色图标
const renderIcon = (iconName: string) => {
    return () => h(ColorfulIcons, { name: iconName, size: 18 })
}

// 转换菜单数据为 Naive UI 格式
const convertMenuData = (data: any[], isCollapsed = false) => {
    return data.map(item => {
        const menuItem: any = {
            label: item.title,
            key: item.index,
            icon: item.icon ? renderIcon(item.icon) : undefined
        }

        // 折叠状态下不显示子菜单
        if (item.children && item.children.length > 0 && !isCollapsed) {
            menuItem.children = convertMenuData(item.children, isCollapsed)
        }

        return menuItem
    })
}

const loadMenuData = async (forceReload = false) => {
    try {
        const res = forceReload ? await reloadMenus() : await getMenus()
        menuData.value = res
        menuOptions.value = convertMenuData(res, sidebar.collapse)
        console.log(`Layout: Menu loaded with ${res.length} items`)
    } catch (error) {
        console.error('Layout: Failed to load menu:', error)
        // 显示默认菜单
        const defaultMenu = [{
            id: 'welcome',
            title: '主页',
            index: '/welcome',
            icon: 'Grid',
        }]
        menuData.value = defaultMenu
        menuOptions.value = convertMenuData(defaultMenu, sidebar.collapse)
    }
}

// 处理菜单选择
const handleMenuSelect = (key: string) => {
    router.push(key)
}



// 初始加载菜单
loadMenuData()

// 监听路由变化，用于重新加载菜单
watch(boolroute, (newVal) => {
    if (newVal) {
        loadMenuData(true) // 强制重新加载
        sidebar.setboolroute(false)
    }
})
</script>

<style scoped>
/* 侧边栏区域样式 */
.custom-sider {
    background: linear-gradient(180deg, #F8FAFC 0%, #F1F5F9 50%, #E2E8F0 100%) !important;
    border-right: 1px solid #E2E8F0 !important;
}

.sidebar-container {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: linear-gradient(180deg, #F8FAFC 0%, #F1F5F9 50%, #E2E8F0 100%);
}


.menu-container {
    flex: 1;
    padding: 16px 0;
    overflow-y: auto;
    position: relative;
}

.collapse-button-container {
    position: absolute;
    right: -12px;
    top: 50%;
    transform: translateY(-50%);
    z-index: 10;
    width: 24px;
    height: 36px;
}

.collapse-toggle-btn {
    width: 24px !important;
    height: 36px !important;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
    border: none !important;
    border-radius: 12px !important;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    position: relative;
    overflow: hidden;
    padding: 0 !important;
    margin: 0 !important;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
}

.collapse-toggle-btn .n-button__content {
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    width: 100% !important;
    height: 100% !important;
    color: white !important;
}

.collapse-toggle-btn .n-icon {
    color: white !important;
}

.collapse-toggle-btn svg {
    color: white !important;
}

.collapse-toggle-btn svg path {
    stroke: white !important;
}

.collapse-toggle-btn::before {
    content: '';
    position: absolute;
    left: -12px;
    top: 0;
    width: 12px;
    height: 100%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 0 12px 12px 0;
}

.collapse-toggle-btn:hover {
    transform: translateX(2px);
    background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 100%) !important;
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.collapse-toggle-btn:hover::before {
    background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 100%);
}

.collapse-toggle-btn:active {
    transform: translateX(1px);
}



/* 主内容区域样式 */
.main-layout {
    background: #FFFFFF !important;
}

.main-content {
    background: #FFFFFF !important;
}

/* Loading 状态下的主布局固定 */
.main-layout.layout-loading {
    overflow: hidden !important;
    height: 100vh !important;
    max-height: 100vh !important;
}

.main-content.content-loading-parent {
    overflow: hidden !important;
    height: 100vh !important;
    max-height: 100vh !important;
}

.content-wrapper {
    background: #FFFFFF;
    border-radius: 6px;
    min-height: calc(100vh - 24px);
    padding: 12px;
    border: 1px solid #F1F5F9;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    font-size: var(--table-font-size, 12px);
    line-height: 1.3;
    position: relative;
    overflow: hidden;
}

/* Loading 状态下的内容区域 */
.content-wrapper.content-loading {
    pointer-events: none;
    user-select: none;
    overflow: hidden !important;
    position: relative;
    /* 强制固定高度，完全锁定内容区域 */
    height: calc(100vh - 24px) !important;
    max-height: calc(100vh - 24px) !important;
    min-height: calc(100vh - 24px) !important;
    /* 阻止所有滚动行为 */
    overscroll-behavior: none;
    scroll-behavior: auto;
    /* 阻止触摸滚动 */
    touch-action: none;
    /* 阻止鼠标滚轮 */
    -ms-overflow-style: none;
    scrollbar-width: none;
    /* 锁定内容，防止任何形式的移动 */
    transform: translateZ(0);
    /* 强制硬件加速，锁定位置 */
    will-change: auto;
    /* 重置 will-change */
}

/* 隐藏滚动条 */
.content-wrapper.content-loading::-webkit-scrollbar {
    display: none;
}

/* 防止 loading 时整个页面滚动 */
.content-wrapper.content-loading * {
    pointer-events: none !important;
    user-select: none !important;
    overflow: hidden !important;
    /* 阻止所有子元素的滚动 */
    overscroll-behavior: none !important;
    touch-action: none !important;
    scroll-behavior: auto !important;
}

/* 特别处理可能的滚动容器 */
.content-wrapper.content-loading .n-scrollbar,
.content-wrapper.content-loading .n-data-table-wrapper,
.content-wrapper.content-loading .n-layout-content {
    overflow: hidden !important;
    overscroll-behavior: none !important;
    touch-action: none !important;
}

/* 主内容区域 Loading 遮罩 */
.content-loading-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.85);
    backdrop-filter: blur(3px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 999999;
    /* 提高层级 */
    border-radius: 6px;
    /* 完全阻止所有交互 */
    pointer-events: all;
    user-select: none;
    touch-action: none;
    /* 防止滚动穿透 - 加强版 */
    overscroll-behavior: none;
    overflow: hidden;
    /* 阻止所有滚动相关的行为 */
    scroll-behavior: auto;
    -webkit-overflow-scrolling: touch;
    /* 确保遮罩能够捕获所有事件 */
    width: 100%;
    height: 100%;
}

.loading-spinner {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 16px;
    /* 移除背景颜色、边框和阴影 */
    background: transparent;
    border: none;
    box-shadow: none;
    /* 确保 loading 区域本身也阻止事件 */
    pointer-events: all;
    user-select: none;
    /* 防止内容被选中或拖拽 */
    -webkit-user-drag: none;
    -webkit-touch-callout: none;
    -webkit-tap-highlight-color: transparent;
}

.loading-gif {
    width: 112px;
    height: 112px;

    /* 确保 GIF 正常播放 */
    image-rendering: auto;
    /* 防止图片被拖拽 */
    -webkit-user-drag: none;
    -moz-user-drag: none;
    -o-user-drag: none;
    user-drag: none;
    /* 防止图片被选中 */
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    /* 防止右键菜单 */
    -webkit-touch-callout: none;
    pointer-events: none;
}

.loading-text {
    margin-top: 12px;
    font-size: 14px;
    color: #666;
    font-weight: 500;
    text-align: center;
    word-wrap: break-word;
    word-break: break-all;
    white-space: pre-wrap;
    max-width: 300px;
    line-height: 1.4;
    padding: 0 16px;
}

.fade-enter-active,
.fade-leave-active {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}

/* 自定义菜单样式 */
:deep(.n-menu) {
    background: transparent !important;
}

:deep(.n-menu .n-menu-item) {
    margin: 2px 6px;
    border-radius: 6px;
    transition: all 0.3s ease;
    color: #64748B !important;
    position: relative;
    overflow: hidden;
}

/* 展开状态的菜单项 */
:deep(.n-menu:not(.n-menu--collapsed) .n-menu-item) {
    margin: 2px 6px;
}

/* 折叠状态的菜单项 */
:deep(.n-menu.n-menu--collapsed .n-menu-item) {
    margin: 3px auto;
    width: 48px;
    height: 48px;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
}

:deep(.n-menu.n-menu--collapsed .n-menu-item-content) {
    padding: 0 !important;
    width: 48px !important;
    height: 48px !important;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    margin: 0 !important;
}

:deep(.n-menu.n-menu--collapsed .n-menu-item-content .n-menu-item-content__icon) {
    margin: 0 !important;
}

:deep(.n-menu .n-menu-item:hover) {
    background: rgba(255, 255, 255, 0.9) !important;
    color: #334155 !important;
    transform: translateX(2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

:deep(.n-menu.n-menu--collapsed .n-menu-item:hover) {
    transform: scale(1.05);
}

:deep(.n-menu .n-menu-item.n-menu-item--selected) {
    background: linear-gradient(135deg, #3B82F6 0%, #1D4ED8 100%) !important;
    color: #FFFFFF !important;
    box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
}

:deep(.n-menu .n-menu-item-content) {
    padding: 8px 12px;
    font-weight: 500;
    font-size: 13px;
    border-radius: 6px;
}

:deep(.n-menu .n-submenu .n-submenu-children .n-menu-item) {
    margin-left: 16px;
    background: transparent !important;
}

:deep(.n-menu .n-submenu .n-submenu-children .n-menu-item:hover) {
    background: rgba(255, 255, 255, 0.7) !important;
}

/* 菜单图标样式 */
:deep(.n-menu .n-menu-item .n-icon) {
    transition: all 0.3s ease;
    color: inherit !important;
    font-size: 16px;
}

:deep(.n-menu.n-menu--collapsed .n-menu-item .n-icon) {
    font-size: 18px;
}

:deep(.n-menu .n-menu-item:hover .n-icon) {
    transform: scale(1.1);
}

/* 隐藏默认的折叠触发器 */
:deep(.n-layout-sider .n-layout-toggle-button) {
    display: none !important;
}

/* 优化内容区域的字体和间距 - 使用CSS变量 */
:deep(.content-wrapper .n-card) {
    font-size: var(--table-font-size, 12px);
    margin-bottom: 12px;
}

:deep(.content-wrapper .n-card-header) {
    padding: 8px 12px;
    font-size: var(--menu-font-size, 14px);
    font-weight: 600;
    min-height: auto;
}

:deep(.content-wrapper .n-card__content) {
    padding: 8px 12px;
    font-size: var(--table-font-size, 12px);
}

:deep(.content-wrapper .n-form-item) {
    margin-bottom: 12px;
}

:deep(.content-wrapper .n-form-item-label) {
    font-size: var(--table-font-size, 12px);
    font-weight: 500;
    padding-bottom: 4px;
}

:deep(.content-wrapper .n-input) {
    font-size: var(--table-font-size, 12px);
    min-height: 28px;
}

:deep(.content-wrapper .n-input .n-input__input-el) {
    padding: 0 8px;
}

:deep(.content-wrapper .n-button) {
    font-size: var(--table-font-size, 12px);
    padding: 0 10px;
    height: 28px;
    min-height: 28px;
}

:deep(.content-wrapper .n-data-table) {
    font-size: var(--table-font-size, 12px);
}

:deep(.content-wrapper .n-data-table th) {
    font-size: var(--table-font-size, 12px);
    font-weight: 600;
    padding: 6px 8px;
    height: 32px;
}

:deep(.content-wrapper .n-data-table td) {
    font-size: var(--table-font-size, 12px);
    padding: 4px 8px;
    height: 28px;
}

:deep(.content-wrapper .n-tabs) {
    font-size: var(--table-font-size, 12px);
}

:deep(.content-wrapper .n-tabs .n-tabs-tab) {
    font-size: var(--menu-font-size, 14px);
    padding: 6px 12px;
}

:deep(.content-wrapper .n-tabs .n-tabs-pane) {
    padding: 8px 0;
}

:deep(.content-wrapper .n-space) {
    gap: 8px !important;
}

:deep(.content-wrapper .n-descriptions) {
    font-size: var(--table-font-size, 12px);
}

:deep(.content-wrapper .n-descriptions-table-wrapper .n-descriptions-table th) {
    font-size: var(--table-font-size, 12px);
    padding: 4px 8px;
}

:deep(.content-wrapper .n-descriptions-table-wrapper .n-descriptions-table td) {
    font-size: var(--table-font-size, 12px);
    padding: 4px 8px;
}

:deep(.content-wrapper .n-tooltip) {
    font-size: var(--table-font-size, 12px);
}

:deep(.content-wrapper h1) {
    font-size: calc(var(--sidebar-font-size, 16px) + 2px);
    margin: 8px 0;
}

:deep(.content-wrapper h2) {
    font-size: var(--sidebar-font-size, 16px);
    margin: 6px 0;
}

:deep(.content-wrapper h3) {
    font-size: var(--menu-font-size, 14px);
    margin: 4px 0;
}

:deep(.content-wrapper p) {
    font-size: var(--table-font-size, 12px);
    margin: 4px 0;
    line-height: 1.3;
}
</style>

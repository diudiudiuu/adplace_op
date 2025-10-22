<template>
    <div>
        <n-message-provider>
            <n-layout has-sider style="height: 100vh;">
                <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="200"
                    :collapsed="sidebar.collapse" :show-trigger="false" class="custom-sider">
                    <div class="sidebar-container">


                        <div class="menu-container">
                            <n-menu v-if="showMenu" :value="onRoutes" :collapsed="sidebar.collapse" :collapsed-width="64"
                                :options="menuOptions" @update:value="handleMenuSelect" :indent="20" />
                            
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

                <n-layout class="main-layout">
                    <n-layout-content content-style="padding: 8px;" class="main-content">
                        <div class="content-wrapper">
                            <router-view v-slot="{ Component }" :key="$route.fullPath">
                                <transition name="fade" mode="out-in">
                                    <keep-alive>
                                        <component :is="Component" />
                                    </keep-alive>
                                </transition>
                            </router-view>
                        </div>
                    </n-layout-content>
                </n-layout>
            </n-layout>
        </n-message-provider>


    </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, h, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useSidebarStore } from '@/store/sidebar'
import { useRoute, useRouter } from 'vue-router'
import { getMenus, reloadMenus } from '@/components/menu'
import { NIcon, useMessage, useDialog } from 'naive-ui'
import { setGlobalInstances, handleUnauthorized } from '@/api'
import ColorfulIcons from '@/components/ColorfulIcons.vue'

const route = useRoute()
const router = useRouter()
const onRoutes = computed(() => route.path)

const sidebar = useSidebarStore()
const { boolroute } = storeToRefs(sidebar)

// 获取 Naive UI 实例
const message = useMessage()
const dialog = useDialog()

// 设置全局实例
onMounted(() => {
    // 创建Naive UI loading实例
    const createLoading = () => {
        return {
            create: (options: any) => {
                const loadingInstance = message.loading(options.description || '加载中...', {
                    duration: 0
                });
                return {
                    destroy: () => loadingInstance.destroy()
                };
            }
        };
    };
    
    setGlobalInstances(message, createLoading())
})

const menuData = ref<any[]>([])
const menuOptions = ref<any[]>([])
const showMenu = ref(true)

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

.content-wrapper {
    background: #FFFFFF;
    border-radius: 6px;
    min-height: calc(100vh - 24px);
    padding: 12px;
    border: 1px solid #F1F5F9;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    font-size: var(--table-font-size, 12px);
    line-height: 1.3;
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

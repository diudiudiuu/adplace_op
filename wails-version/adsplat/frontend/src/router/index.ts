import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { isAuthorized } from '@/utils/auth';
import { getMenus } from '@/components/menu';
import layoutRouter from './layout'

const routes: RouteRecordRaw[] = [
    {
        path: '/lock',
        name: 'lock',
        component: () => import(/* webpackChunkName: "lock" */ '@/views/common/lock.vue'),
    },
    // 始终添加 layout 路由
    layoutRouter,
    {
        path: '/:pathMatch(.*)*',
        redirect: '/lock',
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

// 动态添加菜单路由的函数（现在主要用于日志记录，实际路由已在 layout.ts 中定义）
export const addDynamicRoutes = (menuData: any[]) => {
    console.log('Dynamic routes initialized with menu data:', menuData);
    // 路由现在通过 layout.ts 中的通用模式处理
    // /project/:pid/:id -> dashboard.vue
    // /project_form/:pid -> form.vue
};

router.beforeEach(async (to, _, next) => {
    if (to.path === '/lock') {
        if (isAuthorized()) {
            next('/');
        } else {
            next();
        }
        return;
    }

    if (!isAuthorized()) {
        next('/lock');
        return;
    }

    // 如果访问项目相关页面，确保菜单数据已加载
    if (to.path.startsWith('/project')) {
        try {
            console.log('Router: Loading menus for project page...');
            await getMenus();
            console.log('Router: Menus loaded successfully');
        } catch (error) {
            console.error('Router: Failed to load menus:', error);
        }
    }

    next();
});

export default router;

import layout from '@/views/layout.vue';

const layoutRouter = {
    path: '/',
    name: 'layout',
    redirect: '/welcome',
    component: layout,
    children: [
        {
            path: '/welcome',
            name: 'welcome',
            component: () => import('@/views/common/welcome.vue'),
        },
        {
            path: '/project/:pid/:id',
            name: 'project-dashboard',
            component: () => import('@/views/project/dashboard.vue'),
        },
        {
            path: '/project_form/:pid',
            name: 'project-form',
            component: () => import('@/views/project/form.vue'),
        },
        // 动态路由将通过 addDynamicRoutes 函数添加
    ],
}

export default layoutRouter;
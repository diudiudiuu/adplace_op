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
        {
            path: '/project_info/:serverId/:projectId',
            name: 'project-info',
            component: () => import('@/views/project/info.vue'),
            props: true,
        },
        {
            path: '/server-management',
            name: 'server-management',
            component: () => import('@/views/server/ServerManagement.vue'),
        },
        {
            path: '/server/:serverId',
            name: 'server-info',
            component: () => import('@/views/server/ServerInfo.vue'),
            props: true,
        },
        {
            path: '/page-capture',
            name: 'page-capture',
            component: () => import('@/views/common/PageCapture.vue'),
        },
        // 动态路由将通过 addDynamicRoutes 函数添加
    ],
}

export default layoutRouter;
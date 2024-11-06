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
            component: () => import(/* webpackChunkName: "server/form" */ '@/views/common/welcome.vue'),
        },
        {
            path: '/project/:pid/:id',
            name: 'project/pid/id',
            component: () => import(/* webpackChunkName: "server/list" */ '@/views/project/dashboard.vue'),
        },
        {
            path: '/project_form/:id',
            name: '/project_form/id',
            component: () => import(/* webpackChunkName: "server/list" */ '@/views/project/form.vue'),
        },
        {
            path: '/theme',
            name: 'theme',
            component: () => import(/* webpackChunkName: "theme" */ '@/views/common/theme.vue'),
        },


    ],
}

export default layoutRouter;
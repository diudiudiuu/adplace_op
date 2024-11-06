import layout from '@/views/layout.vue';
const layoutRouter = {
    path: '/',
    name: 'layout',
    component: layout,
    children: [
        {
            path: '/server_form',
            name: 'server_form',
            component: () => import(/* webpackChunkName: "server/form" */ '@/views/server/form.vue'),
        },
        {
            path: '/project/:id',
            name: '/project/id',
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
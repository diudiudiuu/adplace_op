import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import Layout from '@/views/layout.vue';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/server_form',
    },
    {
        path: '/',
        name: 'Layout',
        component: Layout,
        children: [
            {
                path: '/server_form',
                name: 'server_form',
                meta: {
                    auth : true,
                },
                component: () => import(/* webpackChunkName: "server/form" */ '@/views/server/form.vue'),
            },
            {
                path: '/project/:id',
                name: '/project/id',
                meta: {
                    auth : true,
                },
                component: () => import(/* webpackChunkName: "server/list" */ '@/views/project/dashboard.vue'),
            },
            {
                path: '/project_form/:id',
                name: '/project_form/id',
                meta: {
                    auth : true,
                },
                component: () => import(/* webpackChunkName: "server/list" */ '@/views/project/form.vue'),
            },
            {
                path: '/theme',
                name: 'theme',
                meta: {
                    lock : true,
                },
                component: () => import(/* webpackChunkName: "theme" */ '@/views/common/theme.vue'),
            },


        ],
    },
    {
        path: '/lock',
        name: 'lock',
        meta: {
            lock: false,
        },
        component: () => import(/* webpackChunkName: "lock" */ '@/views/common/lock.vue'),
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((_to, _from, next) => {
    NProgress.start();
    // 调用tauri API获取是否锁定
    
    next();
});

router.afterEach(() => {
    NProgress.done();
});

export default router;

import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import Layout from '@/views/layout.vue';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        redirect: '/server-list',
    },
    {
        path: '/',
        name: 'Layout',
        component: Layout,
        children: [
            {
                path: '/server-list',
                name: 'server-list',
                meta: {
                    auth : true,
                },
                component: () => import(/* webpackChunkName: "server/list" */ '@/views/server/list.vue'),
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
    history: createWebHashHistory(),
    routes,
});

router.beforeEach((to, from, next) => {
    NProgress.start();
    // 调用tauri API获取是否锁定
    
    next();
});

router.afterEach(() => {
    NProgress.done();
});

export default router;

import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

const routes: RouteRecordRaw[] = [
    {
        path: '/lock',
        name: 'lock',
        component: () => import(/* webpackChunkName: "lock" */ '@/views/common/lock.vue'),
    },
    {
        path: '/:pathMatch(.*)*',
        redirect: '/lock',
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((_to, _from, next) => {
    NProgress.start();
    next();
});

router.afterEach(() => {
    NProgress.done();
});

export default router;

import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';

import layoutRouter from './layout'
const key = 'GKQSZuLkJI0nPV65'

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

if (localStorage.getItem(key)) {
    //删除
    localStorage.removeItem(key)
    router.addRoute(layoutRouter)
}

router.beforeEach((_to, _from, next) => {
    NProgress.start();
    next();
});

router.afterEach(() => {
    NProgress.done();
});

export default router;

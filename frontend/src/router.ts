import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from './utils';

const router = createRouter({
    history: createWebHistory(),
    scrollBehavior() {
        return { top: 0, behavior: 'smooth' } // aby scroll nezustaval dole na strankach kde se nescrolluje
    },
    routes: [
        {
            path: '/',
            component: () => import('./views/Domu.vue'),
        },
        {
            path: '/o-nas',
            component: () => import('./views/ONas.vue')
        },
        {
            path: '/lekce',
            component: () => import('./views/SeznamLekci.vue')
        },
        {
            path: '/jak-psat',
            component: () => import('./views/Teorie.vue'),
        },
        {
            path: '/prihlaseni',
            component: () => import('./views/Prihlaseni.vue')
        },
        {
            path: '/registrace',
            component: () => import('./views/Registrace.vue')
        },
        {
            path: '/ucet',
            component: () => import('./views/Ucet.vue'),
            meta: { requireAuth: true }
        },
        {
            path: '/lekce/:pismena',
            component: () => import('./views/Lekce.vue'),
        },
        {
            path: '/lekce/:pismena/:id',
            component: () => import('./views/Cviceni.vue'),
            meta: { requireAuth: true }
        },
        {
            path: '/:pathMatch(.*)*',
            component: () => import('./views/404.vue')
        }

    ]
})

router.beforeEach((to, _, next) => { // kdyz potrebuje auth tak => prihlaseni
    if (to.meta.requireAuth) { 
        if (!getToken()) {
            next("/prihlaseni");
        } else {
            next();
        }
    } else {
        next();
    }
});

export default router
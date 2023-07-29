import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    scrollBehavior() {
        return { top: 0, behavior: 'smooth' } // aby scroll nezustaval dole na strankach kde se nescrolluje
    },
    routes: [
        {
            path: '/',
            component: () => import('./views/Domu.vue')
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
            component: () => import('./views/Teorie.vue')
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
            component: () => import('./views/Ucet.vue')
        },
        {
            path: '/lekce/:pismena',
            component: () => import('./views/Lekce.vue'),
        },
        {
            path: '/lekce/:pismena/:id',
            component: () => import('./views/Cviceni.vue')
        },
        {
            path: '/:pathMatch(.*)*',
            component: () => import('./views/404.vue')
        }

    ]
})

export default router

// Composables
import { createRouter, createWebHistory } from 'vue-router'

const routes: any = [
  {
    path: '/',
    component: () => import("@/views/UDefaultView.vue"),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/UHome.vue'),
      },
    ],
  },
  {
    path: '/user',
    component: () => import('@/views/UDefaultView.vue'),
    children: [
      {
        path: ':id',
        name: 'User',
        component: () => import(/* webpackChunkName: "user" */ '@/views/UShow.vue'),
      },
      {
        path: 'edit/:id',
        name: 'Edit',
        component: () => import(/* webpackChunkName: "edit" */ '@/views/UEdit.vue'),
      },
      {
        path: 'new',
        name: 'New',
        component: () => import(/* webpackChunkName: "new" */ '@/views/UNew.vue'),
      },
    ],
  },
  {
    path: '/sng-home',
    name: 'sng-home',
    beforeEnter() {
      const sngHome = import.meta.env.VITE_SNG_HOME
      window.location.replace(sngHome)
    }
  },
  {
    path: '/login',
    name: 'Login',
    beforeEnter() {
      const url = `${import.meta.env.VITE_USER_BACKEND}/sn/user/login`
      window.location.replace(url)
    }
  },
  {
    path: '/logout',
    name: 'Logout',
    beforeEnter() {
      const query = `?redirect=${encodeURIComponent(import.meta.env.VITE_FRONTEND)}`
      const url = `${import.meta.env.VITE_USER_BACKEND}/sn/user/logout${query}`
      window.location.replace(url)
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})
export default router

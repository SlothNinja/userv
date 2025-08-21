// Composables
import { createRouter, createWebHistory } from 'vue-router'

const sngGames = [ 'confucius', 'indonesia', 'all' ]
const sngXGames = [ 'atf', 'confucius' ]

const routes: any = [
  {
    path: '/',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/components/Home.vue'),
      },
    ],
  },
  {
    path: '/user',
    component: () => import('@/layouts/default/Default.vue'),
    children: [
      {
        path: '/user/:id',
        name: 'User',
        component: () => import(/* webpackChunkName: "user" */ '@/components/User/Show.vue'),
      },
      {
        path: '/user/edit/:id',
        name: 'Edit',
        component: () => import(/* webpackChunkName: "edit" */ '@/components/User/Edit.vue'),
      },
      {
        path: '/user/new',
        name: 'New',
        component: () => import(/* webpackChunkName: "new" */ '@/components/User/New.vue'),
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
    path: '/sng-games/:type/:status',
    name: 'sng-games',
    beforeEnter(to: any) {
      if (sngGames.includes(to.params.type)) {
        const sngHome = import.meta.env.VITE_SNG_HOME
        window.location.replace(`${sngHome}#/games/${to.params.status}/${to.params.type}`)
      } else {
        const gotHome = import.meta.env.VITE_GOT_HOME
        const tammanyHome = import.meta.env.VITE_TAMMANY_HOME
        const plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/games/${to.params.status}`)
            break;
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/games/${to.params.status}`)
            break;
          case 'plateau':
            window.location.replace(`${plateauHome}games/${to.params.status}`)
        }
      }
    }
  },
  {
    path: '/sng-ratings/:type',
    name: 'sng-ratings',
    beforeEnter(to: any) {
      if (sngGames.includes(to.params.type)) {
        const sngHome = import.meta.env.VITE_SNG_HOME
        window.location.replace(`${sngHome}ratings/show/${to.params.type}`)
      } else {
        let gotHome = import.meta.env.VUE_APP_GOT_HOME
        let tammanyHome = import.meta.env.VUE_APP_TAMMANY_HOME
        let plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/rank`)
            break
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/rank`)
            break
          case 'plateau':
            window.location.replace(`${plateauHome}join`)
        }
      }
    }
  },
  {
    path: '/sng-ugames/:uid/:status/:type',
    name: 'sng-ugames',
    beforeEnter(to: any) {
      const sngHome = import.meta.env.VITE_SNG_HOME
      window.location.replace(`${sngHome}#/ugames/${to.params.uid}/${to.params.status}/${to.params.type}`)
    }
  },
  {
    path: '/sng-new-game/:type',
    name: 'sng-new-game',
    beforeEnter(to: any) {
      if (sngGames.includes(to.params.type)) {
        let sngHome = import.meta.env.VITE_SNG_HOME
        if (sngXGames.includes(to.params.type)) {
          window.location.replace(`${sngHome}#/invitation/new/${to.params.type}`)
        } else {
          window.location.replace(`${sngHome}${to.params.type}/game/new`)
        }
      } else {
        let gotHome = import.meta.env.VITE_GOT_HOME
        let tammanyHome = import.meta.env.VITE_TAMMANY_HOME
        let plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/invitation/new`)
            break
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/invitation/new`)
            break
          case 'plateau':
            window.location.replace(`${plateauHome}new`)
        }
      }
    }
  },
  {
    path: '/sng-join-game/:type',
    name: 'sng-join-game',
    beforeEnter(to :any) {
      if (sngGames.includes(to.params.type)) {
        const sngHome = import.meta.env.VITE_SNG_HOME
        if (sngXGames.includes(to.params.type)) {
          window.location.replace(`${sngHome}#/invitations/${to.params.type}`)
        } else {
          window.location.replace(`${sngHome}${to.params.type}/games/recruiting`)
        }
      } else {
        let gotHome = import.meta.env.VITE_GOT_HOME
        let tammanyHome = import.meta.env.VITE_TAMMANY_HOME
        let plateauHome = import.meta.env.VITE_PLATEAU_HOME
        switch (to.params.type) {
          case 'got':
            window.location.replace(`${gotHome}#/invitations`)
            break
          case 'tammany2':
            window.location.replace(`${tammanyHome}#/invitations`)
            break
          case 'plateau':
            window.location.replace(`${plateauHome}join`)
        }
      }
    }
  },
  {
    path: '/user/login',
    name: 'Login',
    beforeEnter() {
      const url = `${import.meta.env.VITE_USER_BACKEND}sn/user/login`
      window.location.replace(url)
    }
  },
  {
    path: '/user/logout',
    name: 'Logout',
    beforeEnter() {
      const query = `?redirect=${encodeURIComponent(import.meta.env.VITE_USER_FRONTEND)}`
      const url = `${import.meta.env.VITE_USER_BACKEND}sn/user/logout${query}`
      window.location.replace(url)
    }
  },
  {
    path: '/user/:id/as',
    name: 'AsUser',
    beforeEnter(to: any) {
      const url = `${import.meta.env.VITE_USER_BACKEND}sn/user/${to.params.id}/as`
      window.location.replace(url)
    }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router

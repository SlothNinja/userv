import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/home/Home'
import New from '@/components/user/New'
import Show from '@/components/user/Show'
import Edit from '@/components/user/Edit'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/new',
      name: 'new',
      component: New
    },
    {
      path: '/show/:id',
      name: 'show',
      component: Show
    },
    {
      path: '/edit/:id',
      name: 'edit',
      component: Edit
    },
    {
      path: '/',
      name: 'user-home',
      component: Home
    },
    {
      path: '/logout',
      name: 'logout',
      beforeEnter() {
        window.location.href = '/logout'
      }
    },
    {
      path: '/login',
      name: 'login',
      beforeEnter() {
        window.location.href = '/login'
      }
    },
    {
      path: '/sng-home',
      name: 'sng-home',
      beforeEnter() {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = sngHome
      }
    },
    {
      path: '/sng-ratings/:type',
      name: 'sng-ratings',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = `${sngHome}ratings/show/${to.params.type}`
      }
    },
    {
      path: '/sng-games/:type/:status',
      name: 'sng-games',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = `${sngHome}${to.params.type}/games/${to.params.status}`
      }
    },
    {
      path: '/sng-new-game/:type',
      name: 'sng-new-game',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_SNG_HOME
        window.location.href = `${sngHome}${to.params.type}/game/new`
      }
    },
    {
      path: '/got-ratings',
      name: 'got-ratings',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${gotHome}#/ratings`
      }
    },
    {
      path: '/got-new-game',
      name: 'got-new-game',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${gotHome}#/invitation/new`
      }
    },
    {
      path: '/got-join-game',
      name: 'got-join-game',
      beforeEnter() {
        let gotHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${gotHome}#/invitations`
      }
    },
    {
      path: '/got-games/:status',
      name: 'got-games',
      beforeEnter(to) {
        let sngHome = process.env.VUE_APP_GOT_HOME
        window.location.href = `${sngHome}#/games/${to.params.status}`
      }
    },
  ]
})

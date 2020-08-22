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
      name: 'home',
      component: Home
    }
  ]
})

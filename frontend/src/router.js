import { createRouter, createWebHistory } from 'vue-router'
import Home from './views/h5/Home.vue'
import AdminLayout from './views/admin/AdminLayout.vue'
import Login from './views/admin/Login.vue'
import Dashboard from './views/admin/Dashboard.vue'
import Activities from './views/admin/Activities.vue'
import Strategies from './views/admin/Strategies.vue'
import StrategyDetail from './views/admin/StrategyDetail.vue'
import Awards from './views/admin/Awards.vue'
import RuleTrees from './views/admin/RuleTrees.vue'
import Orders from './views/admin/Orders.vue'
import Users from './views/admin/Users.vue'
import Tasks from './views/admin/Tasks.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/admin/login', component: Login },
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      { path: '', redirect: '/admin/dashboard' },
      { path: 'dashboard', component: Dashboard },
      { path: 'activities', component: Activities },
      { path: 'strategies', component: Strategies },
      { path: 'strategies/:strategyId', component: StrategyDetail },
      { path: 'awards', component: Awards },
      { path: 'rule-trees', component: RuleTrees },
      { path: 'orders', component: Orders },
      { path: 'users', component: Users },
      { path: 'tasks', component: Tasks }
    ]
  }
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to) => {
  if (to.path.startsWith('/admin') && to.path !== '/admin/login') {
    if (!localStorage.getItem('bm_admin_token')) return '/admin/login'
  }
})

export default router

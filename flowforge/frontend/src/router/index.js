import { createRouter, createWebHistory } from 'vue-router'
import { authService } from '../services/api'
import Login from '../views/Login.vue'
import Dashboard from '../views/Dashboard.vue'
import Templates from '../views/Templates.vue'
import WorkflowDetail from '../views/WorkflowDetail.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard, meta: { requiresAuth: true } },
  { path: '/templates', component: Templates, meta: { requiresAuth: true } },
  { path: '/workflow/:id', component: WorkflowDetail, meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(r => r.meta.requiresAuth)
  if (requiresAuth && !authService.getToken()) {
    next('/login')
  } else if (to.path === '/login' && authService.getToken()) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router

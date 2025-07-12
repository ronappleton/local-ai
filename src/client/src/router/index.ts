import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import AdminLayout from '../layouts/AdminLayout.vue'
import ChatLayout from '../components/ChatLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Chat',
    component: ChatLayout
  },
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('../views/admin/AdminDashboard.vue')
      },
      {
        path: 'models',
        name: 'ModelList',
        component: () => import('../views/admin/ModelList.vue')
      },
      {
        path: 'models/active',
        name: 'ActiveModels',
        component: () => import('../views/admin/ActiveModels.vue')
      },
      {
        path: 'models/disabled',
        name: 'DisabledModels',
        component: () => import('../views/admin/DisabledModels.vue')
      },
      {
        path: 'models/import',
        name: 'ModelImport',
        component: () => import('../views/admin/ModelImport.vue')
      },
      {
        path: 'models/stats/global',
        name: 'GlobalModelStats',
        component: () => import('../views/admin/GlobalModelStats.vue')
      },
      {
        path: 'models/:id',
        name: 'ModelDetail',
        component: () => import('../views/admin/ModelDetail.vue')
      },
      {
        path: 'models/:id/stats',
        name: 'ModelStats',
        component: () => import('../views/admin/ModelStats.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Placeholder for auth guard
router.beforeEach((to, from, next) => {
  // TODO: integrate authentication
  next()
})

export default router

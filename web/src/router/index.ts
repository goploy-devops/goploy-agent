import { RouteRecordRaw, createRouter, createWebHashHistory } from 'vue-router'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
/* Layout */
import Layout from '@/layout/index.vue'

export const navbarRoutes: RouteRecordRaw = {
  path: '/',
  redirect: '/general',
  component: Layout,
  meta: { hidden: true },
  children: [
    {
      path: 'general',
      name: 'General',
      component: () => import('@/views/index.vue'),
      meta: {
        title: 'General',
        icon: 'el-icon-s-grid',
      },
    },
    {
      path: 'loadavg',
      name: 'Loadavg',
      component: () => import('@/views/index.vue'),
      meta: {
        title: 'Loadavg',
        icon: 'el-icon-odometer',
      },
    },
    {
      path: 'ram',
      name: 'Ram',
      component: () => import('@/views/index.vue'),
      meta: {
        title: 'Ram',
        icon: 'el-icon-help',
      },
    },
    {
      path: 'disk',
      name: 'Disk',
      component: () => import('@/views/index.vue'),
      meta: {
        title: 'Disk',
        icon: 'el-icon-postcard',
      },
    },
    {
      path: 'disk-io',
      name: 'DiskIO',
      component: () => import('@/views/index.vue'),
      meta: {
        title: 'DiskIO',
        icon: 'el-icon-bank-card',
      },
    },
    {
      path: 'cpu',
      name: 'Cpu',
      component: () => import('@/views/index.vue'),
      meta: {
        title: 'Cpu',
        icon: 'el-icon-cpu',
      },
    },
  ],
}

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes: RouteRecordRaw[] = [
  navbarRoutes,
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/404.vue'),
    meta: { hidden: true },
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior() {
    return {
      el: '#app',
      left: 0,
      behavior: 'smooth',
    }
  },
  routes: constantRoutes,
})

NProgress.configure({ showSpinner: false }) // NProgress Configuration

router.beforeEach(() => {
  // start progress bar
  NProgress.start()
  return true
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})

export default router

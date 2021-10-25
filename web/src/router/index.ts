import { RouteRecordRaw, createRouter, createWebHashHistory } from 'vue-router'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin', 'manager', 'group-manager', 'member']   control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */
export const homeRoutes: RouteRecordRaw[] = [
  // 预留常量 permission.js 会修改权限的第一条
  {
    path: '/',
    name: 'Index',
    component: () => import('@/views/index.vue'),
    meta: { hidden: true },
  },
]
/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/404.vue'),
    meta: { hidden: true },
  },
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user permission_uri
 */
export const asyncRoutes: RouteRecordRaw[] = []

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

export default router

import router, { homeRoutes, constantRoutes } from './router'
import { RouteRecordRaw } from 'vue-router'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style

NProgress.configure({ showSpinner: false }) // NProgress Configuration
homeRoutes.forEach((route: RouteRecordRaw) => router.addRoute(route))
constantRoutes.forEach((route: RouteRecordRaw) => router.addRoute(route))
router.beforeEach(() => {
  // start progress bar
  NProgress.start()
  return true
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})

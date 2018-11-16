import Vue from 'vue'
import Router from 'vue-router'
import Main from '@/components/Main'

Vue.use(Router)

const routes = [{
  path: '/',
  name: 'Main',
  component: Main
}]

export default new Router({
  hashbang: false,
  mode: 'history',
  linkActiveClass: 'active',
  routes
})

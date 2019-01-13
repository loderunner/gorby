import Vue from 'vue'
import Router from 'vue-router'
import Requests from '@/views/Requests'

Vue.use(Router)

const routes = [{
  path: '/',
  name: 'Main',
  component: Requests
}]

export default new Router({
  hashbang: false,
  mode: 'history',
  linkActiveClass: 'active',
  routes
})

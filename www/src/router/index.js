import Vue from 'vue'
import Router from 'vue-router'
import RequestsTable from '@/components/RequestsTable'

Vue.use(Router)

const routes = [{
  path: '/',
  name: 'Main',
  component: RequestsTable
}]

export default new Router({
  hashbang: false,
  mode: 'history',
  linkActiveClass: 'active',
  routes
})

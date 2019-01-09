import Vue from 'vue'
import App from './App'
import Buefy from 'buefy'
import moment from 'vue-moment'
import router from './router'
import store from './store'

Vue.use(Buefy)
Vue.use(moment)

new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
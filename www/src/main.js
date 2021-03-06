import Vue from 'vue'
import App from './App'
import Buefy from 'buefy'
import VueNumerals from 'vue-numerals'
import Moment from 'vue-moment'
import router from '@/router'
import store from '@/store'


Vue.use(Buefy)
Vue.use(VueNumerals, { locale: 'en' })
Vue.use(Moment)


new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
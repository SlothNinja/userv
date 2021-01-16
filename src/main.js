import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import router from './router/router'
import axios from 'axios'
import { Plugin } from 'vue-fragment'

Vue.use(Plugin)

const _ = require('lodash')

Vue.config.productionTip = false

new Vue({
  vuetify,
  data () {
    return {
      cu: null,
      nav: false,
      cuLoading: true,
      snackbar: { open: false, message: '' }
    }
  },
  created () {
      let self = this
      self.fetchCurrentUser()
  },
  watch: {
    '$route': 'fetchCurrentUser'
  },
  methods: {
    fetchCurrentUser () {
      let self = this
      axios.get('/current')
        .then(function (response) {
          let cu = _.get(response, 'data.cu', false)
          if (cu) {
            self.cu = cu
          }
                
          self.cuLoading = false
        })
        .catch(function (response) {
          let msg = _.get(response, 'data.message', false)
          if (msg) {
            self.snackbar.message = 'Server Error.  Try again.'
            self.snackbar.open = true
          }
          self.cuLoading = false
        })
    },
  },
  router,
  render: h => h(App),
}).$mount('#app')

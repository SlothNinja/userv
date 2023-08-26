// Components
import App from './App.vue'

// Composables
import { createApp, computed, readonly, unref, ref, watch } from 'vue'
import { VueFire, VueFireAuth } from 'vuefire'
import { firebaseApp } from '@/composables/firebase'

// Plugins
import { registerPlugins } from '@/plugins'

const app = createApp(App)

// creat app
app
  .use(VueFire, {
    // imported above but could also just be created here
    firebaseApp,
    modules: [
      // we will see other modules later on
      VueFireAuth(),
    ],
  })

/////////////////////////////////////////////////////
// get and provide current user
import  { useFetch } from '@/composables/fetch'
import _get from 'lodash/get'
import { cuKey } from '@/composables/keys'

const cuURL = `${import.meta.env.VITE_USER_BACKEND}sn/user/current`
const { data, error } = useFetch(cuURL)
const cu = ref({})

watch( data, () => { cu.value = _get(unref(data), 'CU', {}) })

function updateCU(user) {
  cu.value = unref(user)
}

app.provide( cuKey, { cu, updateCU })
////////////////////////////////////////////////////

registerPlugins(app)

app.mount('#app')

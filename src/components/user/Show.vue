<template>
  <v-app id='app'>
    <sn-toolbar v-model='nav'></sn-toolbar>
    <sn-nav-drawer v-model='nav' app></sn-nav-drawer>
    <sn-snackbar v-model='snackbar.open'>
      <div class='text-center'>
        {{snackbar.message}}
      </div>
    </sn-snackbar>
    <v-main>
      <v-container fluid>
        <v-card>
          <v-container>
            <v-row>
              <v-col cols='6'>
                <v-card-text v-if="loading" class="text-xs-center">
                  <v-progress-circular
                    indeterminate
                    color="green"
                    size="128"
                    width="10"
                    >Loading...</v-progress-circular>
                </v-card-text>
                <template v-else>
                  <v-card-title primary-title>
                    <div class="font-weight-bold title">
                      <sn-user-btn medium :user="u" > {{u.name}} </sn-user-btn>
                    </div>
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-text-field
                        name="name"
                        label="Screen Name"
                        v-model="u.name"
                        id="name"
                        readonly
                        >
                      </v-text-field>
                    </v-row>
                    <template v-if="isCuOrAdmin(u)">
                      <v-row>
                        <v-text-field
                          name="email"
                          label="Email"
                          v-model="u.email"
                          id="email"
                          readonly
                          >
                        </v-text-field>
                      </v-row>
                      <v-row>
                        <v-col>
                          <v-checkbox
                            v-model="u.emailReminders"
                            label="Email Reminders"
                            readonly
                            color="green"
                            ></v-checkbox>
                        </v-col>
                        <v-col>
                          <v-checkbox
                            v-model="u.emailNotifications"
                            label="Email Notifications"
                            color="green"
                            readonly
                            ></v-checkbox>
                        </v-col>
                      </v-row>
                      <v-row v-if="isAdmin(this.cu)">
                        <v-col>
                          <v-btn color="green" dark :to="{ name: 'edit', params: { id: $route.params.id }}">Edit</v-btn>
                        </v-col>
                        <v-col>
                          <v-btn @click="asUser" color="green" dark>As ({{u.name}})</v-btn>
                        </v-col>
                      </v-row>
                      <v-row v-if="isCu(u)">
                        <v-col>
                          <v-btn color="green" dark :to="{ name: 'edit', params: { id: $route.params.id }}">Edit</v-btn>
                        </v-col>
                      </v-row>
                    </template>
                  </v-card-text>
                </template>
              </v-col>
              <v-col cols='6'>
                <h1>Welcome to SlothNinja Games</h1>
                <p>
                SlothNinja Games is a play-by-web (PBW) site that permits registered members to play boardgames
                with other members via the Internet, in a turn-based manner. Registration is required in order
                to play. Registration currently requires a Google Account, but is free.
                </p>
                <p>
                Please specify the screen name for you account.  If you have multiple Google Accounts, please
                verify that the displayed email address is for the Google Account that you wish to register.  If the
                wrong Google Account, select the login button in the toolbar and login in with the correct
                Google Account.
                </p>
                <p>
                Select whether the system should send an email notification when it's your turn.  Regardless
                of selection, the site will send a daily email reminder if any games are waiting for you to
                make a move.
                </p>
                <p>
                Finally, you can also select your avatar from default avatars.  Avatars are provided by
                gravatar.com.  See, gravatar.com to personalize your avatar.
                </p>
              </v-col>
            </v-row>
          </v-container>
        </v-card>
      </v-container>
    </v-main>
    <sn-footer app></sn-footer>
  </v-app>
</template>

<script>
import Toolbar from '@/components/lib/Toolbar'
import NavDrawer from '@/components/lib/NavDrawer'
import Snackbar from '@/components/lib/Snackbar'
import Footer from '@/components/lib/Footer'
import UserButton from '@/components/lib/user/Button'
import CurrentUser from '@/components/lib/mixins/CurrentUser'

const axios = require('axios')
const _ = require('lodash')

export default {
  mixins: [ CurrentUser ],
  data () {
    return {
      u: { name: '', emailReminders: 'true', emailNotifications: 'true', gravType: 'monsterid' },
      auto: true,
      nav: false,
      idToken: '',
      loading: true
    }
  },
  components: {
    'sn-toolbar': Toolbar,
    'sn-nav-drawer': NavDrawer,
    'sn-snackbar': Snackbar,
    'sn-footer': Footer,
    'sn-user-btn': UserButton
  },
  computed: {
    snackbar: {
      get: function () {
        return this.$root.snackbar
      },
      set: function (value) {
        this.$root.snackbar = value
      }
    },
  },
  created () {
    let self = this
    self.fetchData()
  },
  watch: {
    '$route': 'fetchData'
  },
  methods: {
    fetchData () {
      let self = this
      axios.get(`json/${self.$route.params.id}`)
        .then(function (response) {
          let u = _.get(response, 'data.user', false)
          if (u) {
            self.u = u
          }
          self.loading = false
        })
        .catch(function (response) {
          let msg = _.get(response, 'data.message', false)
          if (msg) {
            self.snackbar.message = msg
            self.snackbar.open = true
          }
          self.loading = false
          self.$router.push({ name: 'show', params: { id: self.$route.params.id}})
        })
    },
    asUser () {
      let self = this
      let path = `/as/${self.$route.params.id}`
      self.loading = true
      axios.get(path)
        .then(function (response) {
          let msg = _.get(response, 'data.message', false)
          if (msg) {
            self.snackbar.message = msg
            self.snackbar.open = true
          }
          self.loading = false
          self.cu = null
          self.$router.push({ name: 'user-home' })
        })
        .catch(function (response) {
          let msg = _.get(response, 'data.message', false)
          if (msg) {
            self.snackbar.message = msg
            self.snackbar.open = true
          }
          self.loading = false
          self.$router.push({ name: 'show', params: { id: self.$route.params.id}})
        })
    }
  }
}
</script>

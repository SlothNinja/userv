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
                      <sn-user-btn size="medium" :user="u" ></sn-user-btn>
                      {{u.name}}
                    </div>
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-text-field
                        name="user-name"
                        label="Screen Name"
                        v-model="u.name"
                        id="user-name"
                        ></v-text-field>
                    </v-row>
                    <v-row>
                      <v-text-field
                        name="user-email"
                        label="Email"
                        v-model="u.email"
                        id="user-email"
                        disabled
                        ></v-text-field>
                    </v-row>
                    <v-row>
                      <v-col>
                        <v-checkbox
                          v-model="u.emailReminders"
                          label="Email Reminders"
                          color="green"
                          disabled
                          ></v-checkbox>
                      </v-col>
                      <v-col>
                        <v-checkbox
                          v-model="u.emailNotifications"
                          label="Email Notifications"
                          color="green"
                          ></v-checkbox>
                      </v-col>
                    </v-row>
                    <v-radio-group v-model="u.gravType" row label='Gravatar'>
                      <v-row>
                        <v-col v-for="t in gravTypes()" :key="t">
                          <v-radio :value="t" color="green">
                            <template slot='label'>
                              <v-avatar size='48px' >
                                <img :src="gravatar(u.emailHash, 48, t)" />
                              </v-avatar>
                            </template>
                          </v-radio>
                        </v-col>
                      </v-row>
                    </v-radio-group>
                    <v-row>
                      <v-col>
                        <v-btn @click="putData" color="green" dark>Register</v-btn>
                      </v-col>
                    </v-row>
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
import Gravatar from '@/components/lib/mixins/Gravatar'
import UserButton from '@/components/lib/user/Button'
import CurrentUser from '@/components/lib/mixins/CurrentUser'

const _ = require('lodash')
const axios = require('axios')

export default {
  mixins: [ Gravatar, CurrentUser ],
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
    'sn-user-btn': UserButton,
    'sn-toolbar': Toolbar,
    'sn-nav-drawer': NavDrawer,
    'sn-snackbar': Snackbar,
    'sn-footer': Footer
  },
  computed: {
    snackbar: {
      get: function () {
        return this.$root.snackbar
      },
      set: function (value) {
        this.$root.snackbar = value
      }
    }
  },
  created () {
    let self = this
    self.fetchData()
  },
  methods: {
    fetchData () {
      var self = this
      axios.get('/new')
        .then(function (response) {
          var msg = _.get(response, 'data.msg', false)
          if (msg) {
            self.snackbar.message = msg
            self.snackbar.open = true
          }
          let u = _.get(response, 'data.user', false)
          if (u) {
            self.u = u
          } else {
            self.$router.push({ name: 'user-home'})
          }
          self.loading = false
        })
        .catch(function () {
          self.loading = false
          self.snackbar.message = 'Server Error.  Try again.'
          self.snackbar.open = true
          self.$router.push({ name: 'user-home'})
        })
    },
    putData () {
      var self = this
      self.loading = true
      axios.put('/new', self.u)
        .then(function (response) {
          let u = _.get(response, 'data.user', false)
          if (u) {
            self.u = u
            if (_.isNil(self.cu)) {
              self.cu = u
            }
          }
          let msg = _.get(response, 'data.message', false)
          if (msg) {
            self.snackbar.message = msg
            self.snackbar.open = true
          }
          self.$router.push({ name: 'show', params: { id: self.u.id }})
          self.loading = false
        })
        .catch(function (response) {
          let msg = _.get(response, 'data.message', false)
          if (msg) {
            self.snackbar.message = msg
            self.snackbar.open = true
          }
          self.loading = false
          self.$router.push({ name: 'user-home'})
        })
    }
  }
}
</script>

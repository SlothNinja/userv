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
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols='6'>
                <v-img
                  :src="require('@/assets/slothninja_logo_fullsize.png')"
                  aspect-ratio='1'
                  height='480'
                  contain
                ></v-img>
              </v-col>

              <v-col cols='6'>
                <div class='title font-weight-black text-center'>Welcome to SlothNinja Games</div>
                <div class='subtitle-1 font-weight-bold text-center'>User Service Subdomain</div>
                <div class='body-1 pb-2'>
                  SlothNinja Games is a play-by-web (PBW) site that permits registered members to play board
                  and card games with other members in a turn-based manner. The site provides a graphical
                  interface where players can see the current game state and execute their actions
                  via a web browser. The site keeps track of the game state, enforces legal moves, and optionally
                  sends email notifications to players when it is their turn to move.
                </div>
                <div class='body-1 pb-2'>
                  Registration is required in order to play. Registration is free, but requires a Google Account.
                </div>
                <div class='body-1 pb-2'>
                  Due to the PBW nature of the site, games typically take days, weeks, or even months to complete
                  depending on how frequently players take their turns. 
                </div>
                <div class='body-1'>
                  There is a general expectation among the community that players will make
                  at least one move each day. If you plan to consistently take longer than a day per move,
                  please create a new game and clearly indicate your desired playing speed in the title.
                </div>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-card>
      </v-container>
    </v-main>
    <sn-footer></sn-footer>
  </v-app>
</template>

<script>
  import Toolbar from '@/components/lib/Toolbar'
  import NavDrawer from '@/components/lib/NavDrawer'
  import Snackbar from '@/components/lib/Snackbar'
  import Footer from '@/components/lib/Footer'

  export default {
    name: 'home',
    components: {
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
      },
      nav: {
        get: function () {
          return this.$root.nav
        },
        set: function (value) {
          this.$root.nav = value
        }
      }
    }
  }
</script>

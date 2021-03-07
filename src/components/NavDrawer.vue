<template>
  <v-navigation-drawer
    clipped
    v-model='drawer'
    light
    app
  >
    <v-list-item>
      <v-list-item-content>
        <v-list-item-title class='title font-weight-black text-center'>
          SlothNinja Games
        </v-list-item-title>
        <v-list-item-subtitle class='subtitle-1 font-weight-bold text-center'>
          Users Service
        </v-list-item-subtitle>
      </v-list-item-content>
    </v-list-item>

    <v-divider></v-divider>

    <v-list
      dense
      nav
    >
      <v-list-item :to="{ name: 'sng-home' }" exact>
        <v-list-item-icon>
          <v-icon>mdi-home</v-icon>
        </v-list-item-icon>
        <v-list-item-content>
          <v-list-item-title>Home</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <template v-if='cu'>
        <!--
        <v-list-item :to="{ name: 'new' }">
          <v-list-item-icon>
            <v-icon>mdi-pencil</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Create</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item :to="{ name: 'invitations' }">
          <v-list-item-icon>
            <v-icon>mdi-playlist-plus</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Join</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item :to="{ name: 'games', params: {status: 'running' } }">
          <v-list-item-icon>
            <v-icon>mdi-playlist-play</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Play</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item @click='logout'>
          <v-list-item-icon>
            <v-icon>mdi-logout</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Logout</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        -->
        <v-list-group
          no-action
          prepend-icon='mdi-pencil'
          >
          <template v-slot:activator>
            <v-list-item-title>Create</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.createlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-plus'
          >
          <template v-slot:activator>
            <v-list-item-title>Join</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.joinlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-play'
          >
          <template v-slot:activator>
            <v-list-item-title>Play</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.playlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-check'
          >
          <template v-slot:activator>
            <v-list-item-title>Completed</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.completedlink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-group
          no-action
          prepend-icon='mdi-star'
          >
          <template v-slot:activator>
            <v-list-item-title>Top Players</v-list-item-title>
          </template>
          <v-list-item v-for="(item, index) in items" :key='index' :to="item.ratingslink" >
            <v-list-item-title>{{item.title}}</v-list-item-title>
          </v-list-item>
        </v-list-group>
        <v-list-item :to="{ name: 'logout'}" >
          <v-list-item-icon>
            <v-icon>mdi-logout</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>Logout</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
  import CurrentUser from '@/components/mixins/CurrentUser'

  export default {
    mixins: [ CurrentUser ],
    name: 'nav-drawer',
    props: [ 'value' ],
    methods: {
      logout: function () {
        var self = this
        self.delete_cookie('sng-oauth')
        self.cu = null
        if (self.$route.name != 'user-home') {
          self.$router.push({ name: 'user-home'})
        }
      },
      delete_cookie: function (name) {
        document.cookie = `${name}= ; domain = ${process.env.VUE_APP_COOKIE_DOMAIN} ; expires = Thu, 01 Jan 1970 00:00:00 GMT`
      },
    },
    computed: {
      items: function () {
        return [
          { 
            createlink: { name: 'sng-new-game', params: { type: 'atf' } },
            joinlink: { name: 'sng-games', params: { type: 'atf', status: 'recruiting' } },
            playlink: { name: 'sng-games', params: { type: 'atf', status: 'running' } },
            completedlink: { name: 'sng-games', params: { type: 'atf', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'atf' } },
            title: "After the Flood"
          },
          { 
            createlink: { name: 'sng-new-game', params: { type: 'confucius' } },
            joinlink: { name: 'sng-games', params: { type: 'confucius', status: 'recruiting' } },
            playlink: { name: 'sng-games', params: { type: 'confucius', status: 'running' } },
            completedlink: { name: 'sng-games', params: { type: 'confucius', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'confucius' } },
            title: "Confucius"
          },
          { 
            createlink: { name: 'got-new-game' },
            joinlink: { name: 'got-join-game' },
            playlink: { name: 'got-games', params: { status: 'running' } },
            completedlink: { name: 'got-games', params: { status: 'completedlink' } },
            ratingslink: { name: 'got-ratings' },
            title: "Guild of Thieves"
          },
          { 
            createlink: { name: 'sng-new-game', params: { type: 'indonesia' } },
            joinlink: { name: 'sng-games', params: { type: 'indonesia', status: 'recruiting' } },
            playlink: { name: 'sng-games', params: { type: 'indonesia', status: 'running' } },
            completedlink: { name: 'sng-games', params: { type: 'indonesia', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'indonesia' } },
            title: "Indonesia"
          },
          { 
            createlink: { name: 'sng-new-game', params: { type: 'tammany' } },
            joinlink: { name: 'sng-games', params: { type: 'tammany', status: 'recruiting' } },
            playlink: { name: 'sng-games', params: { type: 'tammany', status: 'running' } },
            completedlink: { name: 'sng-games', params: { type: 'tammany', status: 'completed' } },
            ratingslink: { name: 'sng-ratings', params: { type: 'tammany' } },
            title: "Tammany Hall"
          },
        ]
      },
      drawer: {
        get: function () {
          var self = this
          return self.value
        },
        set: function (value) {
          var self = this
          self.$emit('input', value)
        }
      }
    }
  }
</script>

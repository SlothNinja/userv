<template>
  <v-navigation-drawer v-model='value'>

    <v-list v-if='isLoggedIn'>
      <v-list-item>
        <UserButton :user='cu' :size='32' />
      </v-list-item>
    </v-list>

    <v-divider></v-divider>

    <v-list nav density='compact' v-model:opened="open">

      <v-list-item prepend-icon="mdi-home" title="Home" :to="{ name: 'sng-home'}" ></v-list-item>

      <v-list-item
          prepend-icon='mdi-account-details'
          v-if="isLoggedIn"
          :to="{ name: 'sng-ugames', params: { status: 'running', type: 'all', uid: cuid} }"
          exact
          title='Your Games'
          />

        <v-list-group v-if='isLoggedIn' prepend-icon='mdi-pencil' >
          <template v-slot:activator='{ props }'>
            <v-list-item v-bind='props' title='Create' />
          </template>
          <v-list-item v-for='(item, index) in items' :key='index' :to='item.createlink' :title='item.title' />
        </v-list-group>

        <v-list-group v-if='cu' prepend-icon='mdi-plus' value='Join' >
          <template v-slot:activator='{ props }'>
            <v-list-item v-bind='props' title='Join' />
          </template>
          <v-list-item v-for='(item, index) in items' :key='index' :to='item.joinlink' :title='item.title' />
        </v-list-group>

        <v-list-group prepend-icon='mdi-play' value='Play' >
          <template v-slot:activator='{ props }'>
            <v-list-item v-bind='props' title='Play' />
          </template>
          <v-list-item v-for='(item, index) in items' :key='index' :to='item.playlink' :title='item.title' />
        </v-list-group>

        <v-list-group prepend-icon='mdi-check' value='Completed' >
          <template v-slot:activator='{ props }'>
            <v-list-item v-bind='props' title='Completed' />
          </template>
          <v-list-item v-for='(item, index) in items' :key='index' :to='item.completedlink' :title='item.title' />
        </v-list-group>

        <v-list-group prepend-icon='mdi-star' value='Top Players' >
          <template v-slot:activator='{ props }'>
            <v-list-item v-bind='props' title='Top Players' />
          </template>
          <v-list-item v-for='(item, index) in items' :key='index' :to='item.ratingslink' :title='item.title' />
        </v-list-group>
 
      <v-divider></v-divider>

      <v-list-item v-if='!isLoggedIn' title='Login' :to="{ name: 'Login' }" prependIcon='mdi-login' ></v-list-item>
      <v-list-item v-if='isLoggedIn' title='Logout' :to="{ name: 'Logout' }" prependIcon='mdi-logout' ></v-list-item>
    </v-list>

  </v-navigation-drawer>
</template>

<script setup>
import { computed, ref, onMounted, inject, unref } from 'vue'
import UserButton from '@/components/Common/UserButton.vue'
import { cuKey } from '@/composables/keys'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import _map from 'lodash/map'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])

const value = computed({
  get() {
    return props.modelValue
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

const types = computed(
  () => {
    return [
      { type: 'plateau', title: 'Le Plateau' },
      { type: 'atf', title: 'After the Flood' },
      { type: 'confucius', title: 'Confucius' },
      { type: 'got', title: 'Guild of Thieves' },
      { type: 'indonesia', title: 'Indonesia' },
      { type: 'tammany', title: 'Tammany Hall' }
    ]
  }
)

const items = computed(
  () => {
    return _map(unref(types), game => {
      return { 
        createlink: { name: 'sng-new-game', params: { type: game.type } },
        joinlink: { name: 'sng-join-game', params: { type: game.type } },
        playlink: { name: 'sng-games', params: { type: game.type, status: 'running' } },
        completedlink: { name: 'sng-games', params: { type: game.type, status: 'completed' } },
        ratingslink: { name: 'sng-ratings', params: { type: game.type } },
        title: game.title
      }
    })
  }
)

const { cu } = inject(cuKey)
const cuid = computed(() => (_get(unref(cu), 'ID', -1)))

const name = computed( () => {
  return _get(cu, 'value.name', '')
})

const open = ref( ['Create'] )

const isLoggedIn = computed(() => !_isEmpty(unref(cu)))

</script>

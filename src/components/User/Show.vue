<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <Form v-if='user' v-model='user' :cu='cu'>
              <v-row>
                <v-col v-if='cu.Admin'>
                  <v-btn small color='green' dark
                    :to="{ name: 'AsUser', params: { id: $route.params.id }}">
                    As ({{user.Name}})
                  </v-btn>
                </v-col>
                <v-col v-if='isCUOrAdmin'>
                  <v-btn small color='green' dark :to="{ name: 'Edit', params: { id: $route.params.id }}">Edit</v-btn>
                </v-col>
              </v-row>
            </Form>
          </v-col>
          <v-col cols='6'>
            <Greeting />
          </v-col>
        </v-row>
      </v-container>
    </v-card>
  </v-container>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton'
import Greeting from '@/components/Greeting'
import Form from '@/components/User/Form'

import { computed, inject, unref } from 'vue'
import { cuKey } from '@/composables/keys'
import { useRoute } from 'vue-router'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

const { cu, fetchCU } = inject(cuKey)
const route = useRoute()
const id = computed(() => _get(route, 'params.id', 0))

import { useFetch } from '@/snvue/fetch'

const url = computed(() => `${import.meta.env.VITE_USER_BACKEND}sn/user/${unref(id)}/json`)
const { data } = useFetch(url).json()

import { useIsCUOrAdmin } from '@/composables/user'
const isCUOrAdmin = computed(() => useIsCUOrAdmin(cu, user))

const user = computed(
  () => {
    if (unref(id) != 0) {
      return _get(unref(data), 'User', null)
    }
    return null
  }
)

</script>

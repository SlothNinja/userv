<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <Form v-if='user' v-model='user' :cu='cu' edit >
              <v-row v-if='isCUOrAdmin'>
                <v-col>
                  <v-radio-group v-model='user.GravType' inline label='Gravatar'>
                    <v-radio v-for='t in gravTypes' :key='t' :value='t' color='green'>
                      <template v-slot:label>
                        <Avatar :hash='user.EmailHash' :size='48' :variant='t' />
                      </template>
                    </v-radio>
                  </v-radio-group>
                </v-col>
              </v-row>
              <v-row v-if='isCUOrAdmin'>
                <v-col>
                  <v-btn @click='putData' color='green' dark>Update</v-btn>
                </v-col>
                <v-col class='text-xs-right'>
                  <v-btn :to="{name: 'User', params: { id: $route.params.id }}" color='green' dark>Cancel</v-btn>
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
import Avatar from '@/components/Common/Avatar'

import { computed, inject, unref, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { cuKey, snackKey } from '@/composables/keys'
import { useGravTypes } from '@/composables/gravatar'

import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

const { cu, updateCU } = inject(cuKey)
const route = useRoute()
const gravTypes = useGravTypes()
const id = computed(() => _get(route, 'params.id', 0))

import { useIsCUOrAdmin } from '@/composables/user'
const isCUOrAdmin = computed(() => useIsCUOrAdmin(cu, user))

import { useFetch, usePut } from '@/composables/fetch'

const getPath = computed(() => `${import.meta.env.VITE_USER_BACKEND}sn/user/${unref(id)}/json`)
const { data, error } = useFetch(getPath)

const user = computed(
  () => {
    if (unref(id) != 0) {
      return _get(unref(data), 'User', null)
    }
    return null
  }
)

const loading = computed(() => _isEmpty(unref(user)))

const putPath = computed(() => `${import.meta.env.VITE_USER_BACKEND}sn/user/${unref(id)}/update`)

// Inject snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

function putData() {
  const { response, error } = usePut(putPath, user)
  watch(response, () => update(response))
}

function update(response) {
  const u = _get(unref(response), 'CU', {})
  if (!_isEmpty(u)) {
    updateCU(u)
  }
  const msg = _get(unref(response), 'Message', '')
  if (!_isEmpty(msg)) {
    updateSnackbar(msg)
  }
}
</script>

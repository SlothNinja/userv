<template>
  <v-container fluid>
    <v-row>
      <v-col cols='6' v-if='loading' class='text-xs-center'>
        <v-card>
          <v-progress-circular
              indeterminate
              color='green'
              size='128'
              width='10'
              >Loading...</v-progress-circular>
        </v-card>
      </v-col>

      <v-col cols='6' v-else>
        <v-card>

          <v-card-title primary-title>
            <div class='my-4 font-weight-bold title'>
              <UserButton :user='user' :size='32' />
            </div>
          </v-card-title>

          <v-card-text>
            <v-row>
              <v-text-field
                  name='user-name'
                  label='Screen Name'
                  v-model='user.Name'
                  id='user-name'
                  readonly
                  ></v-text-field>
            </v-row>
            <v-row v-if='isCUOrAdmin'>
              <v-text-field
                  name='user-email'
                  label='Email'
                  v-model='user.Email'
                  id='user-email'
                  readonly
                  ></v-text-field>
            </v-row>
            <v-row v-if='isCUOrAdmin'>
              <v-col>
                <v-checkbox
                    v-model='user.EmailReminders'
                    label='Email Reminders'
                    color='green'
                    ></v-checkbox>
              </v-col>
              <v-col>
                <v-checkbox
                    v-model='user.EmailNotifications'
                    label='Email Notifications'
                    color='green'
                    ></v-checkbox>
              </v-col>
            </v-row>
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
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols='6'>
        <Greeting />
      </v-col>

    </v-row>
  </v-container>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton'
import Greeting from '@/components/Greeting'
import Avatar from '@/components/Common/Avatar'

import { computed, inject, unref, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { cuKey, snackKey } from '@/composables/keys'
import { useGravTypes, useGravatar } from '@/composables/gravatar'

import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

const { cu, updateCU } = inject(cuKey)
const route = useRoute()
const router = useRouter()
const gravTypes = useGravTypes()
const gravatar = useGravatar
const id = computed(() => _get(route, 'params.id', 0))

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

const cuid = computed(() => _get(unref(cu), 'ID', -1))
const uid = computed(() => _get(unref(user), 'ID', -2))
const isCU = computed(() => unref(cuid) == unref(uid))
const isAdmin = computed(() => _get(unref(cu), 'Admin', false))
const isCUOrAdmin = computed(() => unref(isCU) || unref(isAdmin))
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

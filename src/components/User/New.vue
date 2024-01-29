<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <Form v-if='user' v-model='user' :cu='cu' :loading='isFetching' create >
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
                  <v-btn @click='putData' color='green' dark>Create</v-btn>
                </v-col>
                <v-col class='text-xs-right'>
                  <v-btn :to="{name: 'Logout'}" color='green' dark>Cancel</v-btn>
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

import { computed, inject, unref, ref, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'

import { cuKey, snackKey } from '@/composables/keys'
import { useIsCUOrAdmin } from '@/composables/user'
import { useGravTypes } from '@/composables/gravatar'

import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

const { cu, updateCU } = inject(cuKey)
const route = useRoute()
const gravTypes = useGravTypes()

const isCUOrAdmin = computed(() => useIsCUOrAdmin(cu, user))

import { useFetch, usePut } from '@/snvue/fetch'

const getPath = computed(() => `${import.meta.env.VITE_USER_BACKEND}sn/user/new`)
const { data, isFinished, isFetching } = useFetch(getPath).json()
const user = ref(null)

watch(
  isFinished,
  () => {
    if(unref(isFinished)) {
      user.value = _get(unref(data), 'User', null)
    }
  }
)

watch(data, () =>
  {
    const error = _get(unref(data), 'Error', '')
    if (!_isEmpty(error)) {
      router.push({name: 'Home'})
    }
  }
)

const putPath = computed(() => `${import.meta.env.VITE_USER_BACKEND}sn/user/new`)

// Inject snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

function putData() {
  const { data } = usePut(putPath, user).json()
  watch(data, () => update(data))
}

function update(response) {
  const cu = _get(unref(response), 'CU', {})
  if (!_isEmpty(cu)) {
    updateCU(cu)
  }

  user.value = _get(unref(response), 'User', user.value)

  const msg = _get(unref(response), 'Message', '')
  if (!_isEmpty(msg)) {
    updateSnackbar(msg)
  }
}
</script>

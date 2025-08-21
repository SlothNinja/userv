<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <Form v-if='user' v-model='user' :cu :size :loading edit >
              <v-row v-if='isAdminOrUser'>
                <v-col>
                  <v-radio-group v-model='user.GravType' inline label='Gravatar'>
                    <v-radio v-for='t in useGravTypes' :key='t' :value='t' color='green'>
                      <template v-slot:label>
                        <Avatar :hash='user.EmailHash' :size='48' :variant='t' />
                      </template>
                    </v-radio>
                  </v-radio-group>
                </v-col>
              </v-row>
              <v-row v-if='isAdminOrUser'>
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

<script setup lang='ts'>
import Greeting from '@/components/Greeting.vue'
import Form from '@/components/User/Form.vue'
import Avatar from '@/snvue/components/Common/Avatar.vue'

import { UserResponse } from '@/composables/types'

import { computed, ref, Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useGravTypes } from '@/snvue/composables/gravatar'
import { useURLPath } from '@/composables/urlPaths'
import { useFetch, usePut } from '@/snvue/composables/fetch'
import { User } from '@/snvue/composables/types'
import { PathName } from '@/composables/types'
import { useIsAdminOrUser } from '@/snvue/composables/user'
import { Snackbar } from '@/snvue/composables/types'
import { updateMessageOrError } from '@/snvue/composables/snackbar'

import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

interface Props {
  size: number
}
defineProps<Props>()
const snackbar = defineModel<Snackbar>('snackbar', { required: true })
const cu = defineModel<User | null>({ required: true })

const user:Ref<User | null> = ref(null)
const route = useRoute()
const id = computed<string>(() => route.params.id as string)

const router = useRouter()
const isAdminOrUser = computed<boolean>(() => useIsAdminOrUser(cu.value, user.value))

const { onFetchResponse } = useFetch(useURLPath(PathName.Edit, id.value))
onFetchResponse(response => response.json().then((data:UserResponse) => update(data)))

function putData() {
  console.log(`user: ${JSON.stringify(user.value)}`)
  const { onFetchResponse } = usePut(useURLPath(PathName.Update, id.value), { User: user.value })
  onFetchResponse(response => response.json().then((data:UserResponse) => {
    update(data)
    router.push({name: 'User', params: { id: id.value}})
  }))
}

function update(
  response?: UserResponse
): void {
  if (response === undefined) {
    return
  }
  updateMessageOrError(snackbar, response)
  if ("User" in response) {
    user.value = response.User
    if (cu.value !== null && (cu.value.ID === user.value?.ID)) {
      cu.value = user.value
    }
  }
}

const loading = computed<boolean>(() => user === null)
</script>

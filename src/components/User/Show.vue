<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <Form v-if='user' v-model='user' :cu :size :loading >
              <v-row>
                <v-col v-if='isAdmin && isNotUser'>
                  <v-btn @click='asUser' small color='green' dark>
                    As ({{name}})
                  </v-btn>
                </v-col>
                <v-col v-if='isAdminOrUser'>
                  <v-btn small color='green' dark @click='edit'>Edit</v-btn>
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

import { computed, ref, Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CUResponse, PathName } from '@/composables/types'
import { Snackbar, User } from '@/snvue/composables/types'
import { useIsAdmin, useIsAdminOrUser } from '@/snvue/composables/user'
import { updateCU } from '@/composables/user'
import { useFetch, usePut } from '@/snvue/composables/fetch'
import { useURLPath } from '@/composables/urlPaths'
import { updateMessageOrError } from '@/snvue/composables/snackbar'


interface Props {
  size: number
}
defineProps<Props>()
const cu = defineModel<User | null>({required: true})
const snackbar = defineModel<Snackbar>('snackbar', { required: true })

const user:Ref<User | null> = ref(null)
const route = useRoute()
const router = useRouter()
const id = computed<string>(() => route.params.id as string)

const { onFetchResponse } = useFetch(useURLPath(PathName.Show, id.value))
onFetchResponse(response => response.json().then(data => user.value = data.User))

const isAdmin = computed<boolean>(() => useIsAdmin(cu.value))
const isNotUser = computed<boolean>(() => cu.value?.ID !== user.value?.ID)
const isAdminOrUser = computed<boolean>(() => useIsAdminOrUser(cu.value, user.value))

const name = computed<string>(() => user.value?.Name ?? '')
const loading = computed<boolean>(() => user === null)

function asUser() {
  const { onFetchResponse } = usePut(useURLPath(PathName.AsUser, id.value), { User: user.value })
  onFetchResponse(response => response.json().then((data:CUResponse) => {
    updateCU(cu, data)
    updateMessageOrError(snackbar, data)
    router.push({name: 'User', params: { id: id.value}})
  }))
}

function edit() {
  router.push({ name: 'Edit', params: { id: id.value}})
}
</script>

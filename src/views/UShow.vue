<template>
  <v-row>
    <v-col cols='6'>
      <UForm v-if='user' v-model:user='user' v-model:cu='cu' :loading>
        <v-row>
          <v-col v-if='isAdmin && isNotUser'>
            <v-btn :height :width class='inherit text-none' @click='asUser' color='primary'>
              As ({{ name }})
            </v-btn>
          </v-col>
          <v-col v-if='isAdminOrUser'>
            <v-btn :height :width class='inherit text-none' color='primary' @click='edit'>Edit</v-btn>
          </v-col>
        </v-row>
      </UForm>
    </v-col>
    <v-col cols='6'>
      <UGreeting />
    </v-col>
  </v-row>
</template>

<script setup lang='ts'>
import UGreeting from '@/components/UGreeting.vue'
import UForm from '@/components/UForm.vue'

import { computed, ref, Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { CUResponse, PathName } from '@/composables/types'
import { useURLPath } from '@/composables/urlPaths'
import { updateCU } from '@/composables/user'
import { useFetch, usePut } from '@/snvue/composables/fetch'
import { updateMessageOrError } from '@/snvue/composables/snackbar'
import { Snackbar, User } from '@/snvue/composables/types'
import { useIsAdmin, useIsAdminOrUser } from '@/snvue/composables/user'

const cu = defineModel<User | undefined | null>('cu', { required: true })
const snackbar = defineModel<Snackbar>('snackbar', { required: true })

const height = '2em'
const width = '8em'

const user: Ref<User | null> = ref(null)
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
  onFetchResponse(response => response.json().then((data: CUResponse) => {
    updateCU(cu, data)
    updateMessageOrError(snackbar, data)
    router.push({ name: 'User', params: { id: id.value } })
  }))
}

function edit() {
  router.push({ name: 'Edit', params: { id: id.value } })
}
</script>

<template>
  <v-row>
    <v-col cols='6'>
      <UForm v-if='user' v-model:user='user' v-model:cu='cu' :loading edit>
        <v-row v-if='isAdminOrUser'>
          <v-col>
            <v-radio-group v-model='user.GravType' inline label='Gravatar'>
              <v-radio v-for='t in useGravTypes' :key='t' :value='t' color='primary'>
                <template v-slot:label class='20em'>
                  <SNUserButton class='ma-4 inherit' noname noclick :user size='medium' :variant='t' />
                </template>
              </v-radio>
            </v-radio-group>
          </v-col>
        </v-row>
        <v-row v-if='isAdminOrUser'>
          <v-col>
            <v-btn :height :width class='inherit text-none' @click='putData' color='primary'>Update</v-btn>
          </v-col>
          <v-col class='text-xs-right'>
            <v-btn :height :width class='inherit text-none' @click='cancel' color='primary'>Cancel</v-btn>
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
import SNUserButton from '@/snvue/components/Common/SNUserButton.vue'

import { computed, ref, Ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { CUResponse, PathName, UserResponse } from '@/composables/types'
import { useURLPath } from '@/composables/urlPaths'
import { updateCU, updateUser } from '@/composables/user'
import { useFetch, usePut } from '@/snvue/composables/fetch'
import { useGravTypes } from '@/snvue/composables/gravatar'
import { updateMessageOrError } from '@/snvue/composables/snackbar'
import { Snackbar, User } from '@/snvue/composables/types'
import { useIsAdminOrUser } from '@/snvue/composables/user'

const cu = defineModel<User | null | undefined>('cu', { required: true })
const snackbar = defineModel<Snackbar>('snackbar', { required: true })

const user: Ref<User | null> = ref(null)
const route = useRoute()

const height = '2em'
const width = '8em'

const router = useRouter()
const id = computed<string>(() => route.params.id as string)
const isAdminOrUser = computed<boolean>(() => useIsAdminOrUser(cu.value, user.value))
const loading = computed<boolean>(() => user === null)

const { onFetchResponse } = useFetch(useURLPath(PathName.Edit, id.value))
onFetchResponse(response => response.json().then((data: UserResponse) => {
  updateUser(cu, user, data)
  updateMessageOrError(snackbar, data)
}))

function putData() {
  const { onFetchResponse } = usePut(useURLPath(PathName.Update, id.value), { User: user.value })
  onFetchResponse(response => response.json().then((data: UserResponse) => {
    updateUser(cu, user, data)
    updateMessageOrError(snackbar, data)
    router.push({ name: 'User', params: { id: id.value } })
  }))
}

function cancel() {
  const { onFetchResponse } = useFetch(useURLPath(PathName.CurrentUser))
  onFetchResponse(response => response.json().then((data: CUResponse) => {
    updateCU(cu, data)
    updateMessageOrError(snackbar, data)
    router.push({ name: 'User', params: { id: id.value } })
  }))
}
</script>

<style lang='sass' scoped>
:deep(.v-selection-control__wrapper)
  height: 1em
  width: 1em

:deep(.v-selection-control__input)
  height: 1em
  width: 1em

:deep(.v-selection-control.v-selection-control--density-default.v-radio)
  font-size: inherit !important
  width: 5.5em !important

:deep(.mdi-radiobox-blank), :deep(.mdi-radiobox-marked)
  font-size: inherit
  width: 0.25em
  margin: 0.25em

:deep(.mdi-checkbox-blank-outline), :deep(.mdi-checkbox-marked)
  font-size: inherit
  width: 0.25em
</style>

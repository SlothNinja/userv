<template>
  <v-row>
    <v-col cols='6'>
      <UForm v-if='user' v-model:user='user' v-model:cu='cu' :loading='isFetching' :size create>
        <v-row v-if='isAdminOrUser'>
          <v-col>
            <v-radio-group v-model='user.GravType' inline label='Gravatar'>
              <v-radio v-for='t in useGravTypes' :key='t' :value='t' color='primary'>
                <template v-slot:label>
                  <SNAvatar :hash='user.EmailHash' size='small' :variant='t' />
                </template>
              </v-radio>
            </v-radio-group>
          </v-col>
        </v-row>
        <v-row v-if='isAdminOrUser'>
          <v-col>
            <v-btn :height :width class='inherit text-none' @click='create' color='primary'>Create</v-btn>
          </v-col>
          <v-col class='text-xs-right'>
            <v-btn :height :width class='inherit text-none' :to="{ name: 'Logout' }" color='primary'>Cancel</v-btn>
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
import SNAvatar from '@/snvue/components/Common/SNAvatar.vue'

import { computed, ref, Ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useFetch, usePut } from '@/snvue/composables/fetch'
import { useGravTypes } from '@/snvue/composables/gravatar'
import { updateMessageOrError } from '@/snvue/composables/snackbar'
import { Snackbar, User } from '@/snvue/composables/types'
import { useIsAdminOrUser } from '@/snvue/composables/user'

import { PathName, UserResponse } from '@/composables/types'
import { useURLPath } from '@/composables/urlPaths'

interface Props {
  size: number
}
defineProps<Props>()

const snackbar = defineModel<Snackbar>('snackbar', { required: true })
const cu = defineModel<User | null | undefined>('cu', { required: true })

const user: Ref<User | null> = ref(null)

const router = useRouter()
const route = useRoute()
const id = computed<string>(() => route.params.id as string)

const height = '2em'
const width = '8em'

const isAdminOrUser = computed(() => useIsAdminOrUser(cu.value, user.value))
const { onFetchResponse, isFetching } = useFetch(useURLPath(PathName.GetNew))
onFetchResponse(response => response.json().then((data: UserResponse) => update(data)))

function create() {
  const { onFetchResponse } = usePut(useURLPath(PathName.Create, id.value), { User: user.value })
  onFetchResponse(response => response.json().then((data: UserResponse) => {
    update(data)
    router.push({ name: 'Home' })
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
  }
}

watch(user, () => cu.value = user.value)

</script>

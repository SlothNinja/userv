<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <Form v-if='user' v-model='user' :cu :loading='isFetching' :size create >
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
                  <v-btn @click='create' color='green' dark>Create</v-btn>
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

<script setup lang='ts'>
import Greeting from '@/components/Greeting.vue'
import Form from '@/components/User/Form.vue'
import Avatar from '@/snvue/components/Common/Avatar.vue'

import { computed, ref, watch, Ref } from 'vue'
import { useRoute } from 'vue-router'
import { useIsAdminOrUser } from '@/snvue/composables/user'
import { useGravTypes } from '@/snvue/composables/gravatar'
import { useFetch, usePut } from '@/snvue/composables/fetch'
import { useURLPath } from '@/composables/urlPaths'
import { User } from '@/snvue/composables/types'
import { Snackbar } from '@/snvue/composables/types'
import { PathName } from '@/composables/types'
import { UserResponse } from '@/composables/types'
import { updateMessageOrError } from '@/snvue/composables/snackbar'
import { useRouter } from 'vue-router'

interface Props {
  size: number
}
defineProps<Props>()
const snackbar = defineModel<Snackbar>('snackbar', { required: true })
const cu = defineModel<User | null>({required: true})

const user:Ref<User | null> = ref(null)

const router = useRouter()
const route = useRoute()
const id = computed<string>(() => route.params.id as string)

const isAdminOrUser = computed(() => useIsAdminOrUser(cu.value, user.value))
const { onFetchResponse, isFetching } = useFetch(useURLPath(PathName.GetNew))
onFetchResponse(response => response.json().then((data:UserResponse) => update(data)))

function create() {
  const { onFetchResponse } = usePut(useURLPath(PathName.Create, id.value), { User: user.value })
  onFetchResponse(response => response.json().then((data:UserResponse) => {
    update(data)
    router.push({name: 'Home'})
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

<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <Form v-if='user' v-model='user' :cu :size :loading >
              <v-row>
                <v-col v-if='isAdmin'>
                  <v-btn small color='green' dark
                    :to="{ name: 'AsUser', params: { id: $route.params.id }}">
                    As ({{name}})
                  </v-btn>
                </v-col>
                <v-col v-if='isAdminOrUser'>
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

<script setup lang='ts'>
import Greeting from '@/components/Greeting.vue'
import Form from '@/components/User/Form.vue'

import { computed, ref, Ref } from 'vue'
import { useRoute } from 'vue-router'
import { PathName } from '@/composables/types'
import { User } from '@/snvue/composables/types'
import { useIsAdmin, useIsAdminOrUser } from '@/snvue/composables/user'
import { useFetch } from '@/snvue/composables/fetch'
import { useURLPath } from '@/composables/urlPaths'

interface Props {
  size: number
}
defineProps<Props>()
const cu = defineModel<User | null>({required: true})

const user:Ref<User | null> = ref(null)
const route = useRoute()
const id = computed<string>(() => route.params.id as string)

const { onFetchResponse } = useFetch(useURLPath(PathName.Show, id.value))
onFetchResponse(response => response.json().then(data => user.value = data.User))

const isAdmin = computed<boolean>(() => useIsAdmin(cu.value))
const isAdminOrUser = computed<boolean>(() => useIsAdminOrUser(cu.value, user.value))

const name = computed<string>(() => user.value?.Name ?? '')
const loading = computed<boolean>(() => user === null)

</script>

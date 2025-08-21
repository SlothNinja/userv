<template>
  <router-view v-model='cu' />
</template>

<script setup lang='ts'>
import { ref, Ref } from 'vue'
import { useURLPath } from '@/composables/urlPaths'
import { useFetch } from '@/snvue/composables/fetch'
import { User } from '@/snvue/composables/types'
import { PathName } from '@/composables/types'

// Current User stuff
type cuAndToken = { CU: User, fsTOKEN?: string }
const { onFetchResponse } = useFetch(useURLPath(PathName.CurrentUser)).json()
onFetchResponse(response => response.json().then((data:cuAndToken) => update(data)))

const token:Ref<string> = ref('')
const cu:Ref<User | null> = ref(null)

function update(data:cuAndToken) {
  cu.value = data.CU
  if (data.fsTOKEN !== undefined) {
    token.value = data.fsTOKEN
  }
}

</script>

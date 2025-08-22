<template>
  <router-view v-model='cu' />
</template>

<script setup lang='ts'>
import { ref, Ref } from 'vue'
import { useURLPath } from '@/composables/urlPaths'
import { useFetch } from '@/snvue/composables/fetch'
import { User } from '@/snvue/composables/types'
import { PathName } from '@/composables/types'
import { CUResponse } from '@/composables/types'
import { updateCU } from './composables/user'

// Current User stuff
const token:Ref<string> = ref('')
const cu:Ref<User | null> = ref(null)

const { onFetchResponse } = useFetch(useURLPath(PathName.CurrentUser))
onFetchResponse(response => response.json().then((data:CUResponse) => updateCU(cu, data, token)))

</script>

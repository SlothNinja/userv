<template>
  <SNDefaultLayout v-model:cu='cu' v-model:nav='nav' v-model:snackbar='snackbar'>
    <template #toolbar>
      <SNToolBar class='font-medium' v-model:nav='nav' />
    </template>

    <template #main>
      <SNDefaultMain v-model:cu='cu' v-model:snackbar='snackbar' />
    </template>
  </SNDefaultLayout>
</template>

<script setup lang='ts'>
import SNDefaultLayout from '@/snvue/layouts/SNDefaultLayout.vue'
import SNToolBar from '@/snvue/components/Default/SNToolBar.vue'
import SNDefaultMain from '@/snvue/components/Default/SNDefaultMain.vue'

import { ref } from 'vue'
import { Snackbar, User } from '@/snvue/composables/types'

const cu = defineModel<User | undefined | null>('cu', { required: true })
const nav = ref<boolean>(false)
const snackbar = ref<Snackbar>(new Snackbar)
</script>

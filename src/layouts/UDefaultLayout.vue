<template>
  <SNDefaultLayout v-model:cu='cu' v-model:nav='nav' v-model:snackbar='snackbar' nogod />
</template>

<script setup lang='ts'>
import { Snackbar, User } from '@/snvue/composables/types'
import { Ref, ref } from 'vue'

const cu = defineModel<User | null>({ required: true })
const nav = ref(false)
const snackbar: Ref<Snackbar> = ref(new Snackbar)
</script>

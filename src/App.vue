<template>
  <router-view />
</template>

<script setup>
/////////////////////////////////////////////////////
// get and provide current user
import { computed, provide, unref } from 'vue'
import  { useFetch } from '@/snvue/fetch'
import { cuKey } from '@/composables/keys'
import _get from 'lodash/get'

const cuURL = `${import.meta.env.VITE_USER_BACKEND}sn/user/current`
const { data } = useFetch(cuURL).json()
const cu = computed(() => _get(unref(data), 'CU', {}))

function updateCU(user) {
  cu.value = unref(user)
}

provide( cuKey, { cu, updateCU })
////////////////////////////////////////////////////
</script>

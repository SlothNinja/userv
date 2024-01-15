<template>
  <router-view />
</template>

<script setup>
/////////////////////////////////////////////////////
// get and provide current user
import { provide, unref, ref, watch } from 'vue'
import  { useFetch } from '@/composables/fetch'
import { cuKey } from '@/composables/keys'
import _get from 'lodash/get'

const cuURL = `${import.meta.env.VITE_USER_BACKEND}sn/user/current`
const { data, error } = useFetch(cuURL)
const cu = ref({})

watch( data, () => { cu.value = _get(unref(data), 'CU', {}) })

function updateCU(user) {
  cu.value = unref(user)
}

provide( cuKey, { cu, updateCU })
////////////////////////////////////////////////////
</script>

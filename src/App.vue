<template>
  <router-view v-model:cu='cu' />
</template>

<script setup lang='ts'>
import { ref, Ref } from 'vue'

import { CUResponse, PathName } from '@/composables/types'
import { useURLPath } from '@/composables/urlPaths'
import { useFetch } from '@/snvue/composables/fetch'
import { User } from '@/snvue/composables/types'
import { updateCU } from '@/composables/user'

// Current User stuff
const token: Ref<string> = ref('')
const cu: Ref<User | null> = ref(null)

const { onFetchResponse } = useFetch(useURLPath(PathName.CurrentUser))
onFetchResponse(response => response.json().then((data: CUResponse) => updateCU(cu, data, token)))

</script>
<style lang='sass'>
$spacer: 0.125em

@for $i from 1 through 32
  .ma-#{$i}
    margin: ($spacer * $i)

  .mx-#{$i}
    margin-left: $spacer * $i
    margin-right: $spacer * $i

  .my-#{$i}
    margin-top: $spacer * $i
    margin-bottom: $spacer * $i

  .ml-#{$i}
    margin-left: ($spacer * $i)

  .mr-#{$i}
    margin-right: ($spacer * $i)

  .mt-#{$i}
    margin-top: ($spacer * $i)

  .mb-#{$i}
    margin-bottom: ($spacer * $i)

  .pa-#{$i}
    padding: ($spacer * $i)

  .px-#{$i}
    padding-left: $spacer * $i
    padding-right: $spacer * $i

  .py-#{$i}
    padding-top: $spacer * $i
    padding-bottom: $spacer * $i

  .pl-#{$i}
    padding-left: ($spacer * $i)

  .pr-#{$i}
    padding-right: ($spacer * $i)

  .pt-#{$i}
    padding-top: ($spacer * $i)

  .pb-#{$i}
    padding-bottom: ($spacer * $i)

.sn-switch
  height: 2em;

.v-switch
  font-size: inherit !important

.inline
  display: flex
  flex-wrap: wrap
  vertical-align: middle
  line-height: 2em

.shaded
  background-color: rgba(0, 0, 0, 0.5)
  border: 1px solid black

.inherit
  font-size: inherit !important

.v-input__control
  .v-selection-control.v-selection-control--density-default
    font-size: inherit !important

.v-label
  opacity: 1 !important
</style>

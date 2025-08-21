<template>
  <v-app>
    <ToolBar :cu :size v-model:nav='nav' />

    <View ref='container' :size v-model='cu' v-model:snackbar='snackbar' />

    <Footer />

    <NavDrawer v-model='nav' :cu />

    <SnackBar v-model='snackbar' />
  </v-app>
</template>

<script setup lang='ts'>
import ToolBar from '@/snvue/components/Default/ToolBar.vue'
import NavDrawer from '@/snvue/components/Default/NavDrawer.vue'
import View from '@/snvue/components/Default/View.vue'
import Footer from '@/snvue/components/Default/Footer.vue'
import SnackBar from '@/snvue/components/Default/SnackBar.vue'
import { useElementSize, refThrottled, type MaybeElement } from '@vueuse/core'
import { computed, Ref } from 'vue'
import { User } from '@/snvue/composables/types'
import { Snackbar } from '@/snvue/composables/types'
import { ref } from 'vue'

const cu = defineModel<User | null>({required: true })
const nav = ref(false)
const snackbar:Ref<Snackbar> = ref(new Snackbar)

const container = ref<MaybeElement>(null)
const { width } = useElementSize(container)

const baseWidth = 1284.0
const percentage = computed<number>(() => width.value / baseWidth)

const bsize = 28
const size = refThrottled(computed<number>(() => percentage.value * bsize), 50)

</script>

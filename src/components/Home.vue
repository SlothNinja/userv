<template>
  <v-container fluid>
    <v-card>
      <v-card-text>
        <div class='title font-weight-black text-center'>Welcome to User Service Subdomain</div>
        <div class='title font-weight-black text-center'>For</div>
        <v-img
            :src='logoFull'
            aspect-ratio='1'
            height='480'
            contain
            ></v-img>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup lang='ts'>
import logoFull from '@/assets/slothninja_logo_fullsize.png'
</script>

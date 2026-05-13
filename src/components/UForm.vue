<template>
  <div>
    <v-card-text v-if='loading' class='text-xs-center'>
      <v-progress-circular indeterminate color='primary' size='10em' width='10em'>Loading...</v-progress-circular>
    </v-card-text>
    <template v-else>
      <v-card-title>
        <div class='font-weight-bold'>
          <SNUserButton :user size='xxlarge' />
        </div>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-text-field name='user-name' label='Screen Name' v-model='user.Name' id='user-name' :readonly='!create'>
          </v-text-field>
        </v-row>
        <template v-if='isCUOrAdmin'>
          <v-row>
            <v-text-field name='user-email' label='Email' v-model='user.Email' id='user-email' readonly>
            </v-text-field>
          </v-row>
          <v-row>
            <v-col>
              <v-checkbox v-model='user.EmailReminders' label='Email Reminders' :readonly='!edit && !create'
                color='primary'></v-checkbox>
            </v-col>
            <v-col>
              <v-checkbox v-model='user.EmailNotifications' label='Email Notifications' color='primary'
                :readonly='!edit && !create'></v-checkbox>
            </v-col>
          </v-row>
          <slot></slot>
        </template>
      </v-card-text>
    </template>
  </div>
</template>

<script setup lang='ts'>
import SNUserButton from '@/snvue/components/Common/SNUserButton.vue'

import { computed } from 'vue'
import { useIsAdmin } from '@/snvue/composables/user'
import { User } from '@/snvue/composables/types'

interface Props {
  edit?: Boolean,
  create?: Boolean,
  loading?: Boolean,
}
const props = defineProps<Props>()
const cu = defineModel<User | undefined | null>('cu', { required: true })
const user = defineModel<User>('user', { required: true })
const isCUOrAdmin = computed(() => (useIsAdmin(cu)) || (cu.value !== null && (cu.value?.ID === user.value.ID)))
</script>

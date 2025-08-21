<template>
  <v-card-text v-if='loading' class='text-xs-center'>
    <v-progress-circular
        indeterminate
        color='green'
        size='128'
        width='10'
        >Loading...</v-progress-circular>
  </v-card-text>
  <template v-else>
    <v-card-title primary-title>
      <div class='font-weight-bold title'>
        <UserButton :user='user' :size='32' />
      </div>
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-text-field
            name='user-name'
            label='Screen Name'
            v-model='user.Name'
            id='user-name'
            :readonly='!create'
            >
        </v-text-field>
      </v-row>
      <template v-if='isCUOrAdmin'>
        <v-row>
          <v-text-field
              name='user-email'
              label='Email'
              v-model='user.Email'
              id='user-email'
              readonly
              >
          </v-text-field>
        </v-row>
        <v-row>
          <v-col>
            <v-checkbox
                v-model='user.EmailReminders'
                label='Email Reminders'
                :readonly='!edit && !create'
                color='green'
                ></v-checkbox>
          </v-col>
          <v-col>
            <v-checkbox
                v-model='user.EmailNotifications'
                label='Email Notifications'
                color='green'
                :readonly='!edit && !create'
                ></v-checkbox>
          </v-col>
        </v-row>
        <slot></slot>
      </template>
    </v-card-text>
  </template>
</template>

<script setup lang='ts'>
import UserButton from '@/snvue/components/Common/UserButton.vue'

import { computed } from 'vue'
import { useIsAdmin } from '@/snvue/composables/user'
import { User } from '@/snvue/composables/types'

interface Props {
  cu: User | null
  size: number
  edit?: Boolean,
  create?: Boolean,
  loading: Boolean,
}
const props = defineProps<Props>()

const user = defineModel<User>({required: true})
const isCUOrAdmin = computed(() => (useIsAdmin(props.cu)) || (props.cu !== null && (props.cu.ID === user.value.ID)))
</script>

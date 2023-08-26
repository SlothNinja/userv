<template>
  <v-container fluid>
    <v-card>
      <v-container>
        <v-row>
          <v-col cols='6'>
            <v-card-text v-if="loading" class="text-xs-center">
              <v-progress-circular
                  indeterminate
                  color="green"
                  size="128"
                  width="10"
                  >Loading...</v-progress-circular>
            </v-card-text>
            <template v-else>
              <v-card-title primary-title>
                <div class="font-weight-bold title">
                  <UserButton :user='user' :size='32' />
                </div>
              </v-card-title>
              <v-card-text>
                <v-row>
                  <v-text-field
                      name="name"
                      label="Screen Name"
                      v-model="user.Name"
                      id="name"
                      readonly
                      >
                  </v-text-field>
                </v-row>
                <template v-if="(cu.Admin || cu.ID == user.ID)">
                  <v-row>
                    <v-text-field
                        name="email"
                        label="Email"
                        v-model="user.Email"
                        id="email"
                        readonly
                        >
                    </v-text-field>
                  </v-row>
                  <v-row>
                    <v-col>
                      <v-checkbox
                          v-model="user.EmailReminders"
                          label="Email Reminders"
                          readonly
                          color="green"
                          ></v-checkbox>
                    </v-col>
                    <v-col>
                      <v-checkbox
                          v-model="user.EmailNotifications"
                          label="Email Notifications"
                          color="green"
                          readonly
                          ></v-checkbox>
                    </v-col>
                  </v-row>
                  <v-row v-if="cu.Admin">
                    <v-col>
                      <v-btn color="green" dark :to="{ name: 'Edit', params: { id: $route.params.id }}">Edit</v-btn>
                    </v-col>
                    <v-col>
                      <v-btn @click="asUser" color="green" dark>As ({{user.Name}})</v-btn>
                    </v-col>
                  </v-row>
                  <v-row v-if="!cu.Admin && (cu.ID == user.ID)">
                    <v-col>
                      <v-btn color="green" dark :to="{ name: 'Edit', params: { id: $route.params.id }}">Edit</v-btn>
                    </v-col>
                  </v-row>
                </template>
              </v-card-text>
            </template>
          </v-col>
          <v-col cols='6'>
            <h1>Welcome to SlothNinja Games</h1>
            <p>
            SlothNinja Games is a play-by-web (PBW) site that permits registered members to play boardgames
            with other members via the Internet, in a turn-based manner. Registration is required in order
            to play. Registration currently requires a Google Account, but is free.
            </p>
            <p>
            Please specify the screen name for you account.  If you have multiple Google Accounts, please
            verify that the displayed email address is for the Google Account that you wish to register.  If the
            wrong Google Account, select the login button in the toolbar and login in with the correct
            Google Account.
            </p>
            <p>
            Select whether the system should send an email notification when it's your turn.  Regardless
            of selection, the site will send a daily email reminder if any games are waiting for you to
            make a move.
            </p>
            <p>
            Finally, you can also select your avatar from default avatars.  Avatars are provided by
            gravatar.com.  See, gravatar.com to personalize your avatar.
            </p>
          </v-col>
        </v-row>
      </v-container>
    </v-card>
  </v-container>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton.vue'
import { computed, inject, unref } from 'vue'
import { cuKey } from '@/composables/keys'
import { useRoute } from 'vue-router'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

const { cu, fetchCU } = inject(cuKey)
const route = useRoute()
const id = computed(() => _get(route, 'params.id', 0))

import { useFetch } from '@/composables/fetch'

const url = computed(() => `${import.meta.env.VITE_USER_BACKEND}sn/user/${unref(id)}/json`)
const { data, error } = useFetch(url)

const user = computed(
  () => {
    if (unref(id) != 0) {
      return _get(unref(data), 'User', null)
    }
    return null
  }
)

const loading = computed(() => _isEmpty(unref(user)))

</script>

import { Ref } from "vue"
import { User } from "@/snvue/composables/types"
import { CUResponse, UserResponse } from "@/composables/types"

export function updateUser(
  cu:Ref<User | null>,
  user:Ref<User | null>,
  response?: UserResponse
): void {
  if (response === undefined) {
    return
  }
  if ("User" in response) {
    user.value = response.User
    if (cu.value !== null && (cu.value.ID === user.value?.ID)) {
      cu.value = user.value
    }
  }
}

export function updateCU(
  cu:Ref<User | null>,
  response?: CUResponse,
  token?:Ref<string>
): void {
  if (response === undefined) {
    return
  }
  if ("CU" in response) {
    cu.value = response.CU
  }
  if ("fsToken" in response && token !== undefined) {
    token.value = response.fsToken ?? ''
  }
}

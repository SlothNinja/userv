import { User } from "@/snvue/composables/types"

export enum PathName {
  None = '',
  CurrentUser = 'current',
  Create = 'create',
  GetNew = 'getNew',
  Show = 'show',
  Edit = 'edit',
  Update = 'update',
  Send = 'send',
  UpdateRead = 'update-read',
}

export type UserResponse = { User: User, Message: string } | { User: User, Error: string }

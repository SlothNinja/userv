import { User } from "@/snvue/composables/types"

export enum PathName {
  None = '',
  CurrentUser = 'current',
  Create = 'create',
  GetNew = 'getNew',
  Show = 'show',
  Edit = 'edit',
  Update = 'update',
  AsUser = 'as-user',
  Send = 'send',
  UpdateRead = 'update-read',
}

export type UserResponse = { User: User, Message: string } | { User: User, Error: string }

export type CUResponse = { CU: User, Message: string, fsTOKEN?: string } | { CU: User, Error: string, fsToken?: string}

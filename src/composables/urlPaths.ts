import { PathName } from '@/composables/types'

export function useURLPath(
  name: PathName,
  id?: string
): string {
  const url = (process.env.NODE_ENV === 'development') ?  import.meta.env.VITE_BACKEND : ''
  switch (name) {
    case PathName.CurrentUser:
      return url + '/sn/user/current'
    case PathName.Create:
      return url + '/sn/user/new/'
    case PathName.GetNew:
      return url + '/sn/user/new/'
    case PathName.Show:
      return url + '/sn/user/' + id + '/json'
    case PathName.Edit:
      return url + '/sn/user/' + id + '/json'
    case PathName.Update:
      return url + '/sn/user/' + id + '/update'
    default:
      console.log(`useURLPath missing path for [${name}, ${id}]`)
      return ''
  }
}

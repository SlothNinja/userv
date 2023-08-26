import { ref } from 'vue'
import _get from 'lodash/get'
import _includes from 'lodash/includes'

export function useGravSizes() {
  return { 'x-small': '24', 'small': '30', 'medium': '48', 'large': '54' }
}

export function useGravTypes() {
  return [ 'personal', 'identicon', 'monsterid', 'retro', 'robohash' ]
}

export function useGravatar(hash, size, t) {
  const sz = _get(useGravSizes(), size, '64')

  if (t == 'personal') {
    return `https://www.gravatar.com/avatar/${hash}?s=${sz}`
  }

  if (!_includes(useGravTypes(), t)) {
    t = 'monsterid'
  }
  return `https://www.gravatar.com/avatar/${hash}?s=${sz}&d=${t}&f=y`
}

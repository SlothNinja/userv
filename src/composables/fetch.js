import { ref, unref } from 'vue'

export function useFetch(url) {
  const data = ref(null)
  const error = ref(null)

  if (process.env.NODE_ENV == 'development') {
    // url = 'https://plateau.fake-slothninja.com:8091' + url
    console.log('fetching: ' + unref(url))
  }

  fetch(unref(url), { credentials: 'include' } )
    .then((res) => res.json())
    .then((json) => (data.value = json))
    .catch((err) => (error.value = err))

  return { data, error }
}

export function usePut(url, data) {
  const response = ref(null)
  const error = ref(null)

  if (process.env.NODE_ENV == 'development') {
    // url = 'https://plateau.fake-slothninja.com:8091' + unref(url)
    console.log('putting: ' + unref(url))
  }

  fetch(unref(url), {
    method: 'PUT',
    credentials: 'include',
    headers: {
      'Content-type': 'application/json'
    },
    body: JSON.stringify(unref(data)),
  } )
    .then((res) => res.json())
    .then((json) => (response.value = json))
    .catch((err) => (error.value = err))

  return { response, error }
}

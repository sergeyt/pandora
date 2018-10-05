const stdHeaders = {
  Accept: 'application/json',
  'Content-Type': 'application/json',
}

function getToken() {
  return localStorage.getItem('access_token')
}

function setToken(value) {
  if (value) {
    localStorage.setItem('access_token', value)
  } else {
    localStorage.removeItem('access_token')
  }
}

export function fetchJSON(url, payload, options = {}) {
  const token = getToken()

  const headers = {
    ...stdHeaders,
    Authorization: token ? `Bearer ${token}` : undefined,
    ...(options.headers || {}),
  }

  return fetch(url, {
    method: options.method || (payload !== undefined ? 'POST' : 'GET'),
    headers,
    body: payload !== undefined ? JSON.stringify(payload) : undefined,
    ...options,
  }).then(resp => resp.json())
}

export function post(url, payload, options = {}) {
  return fetchJSON(url, payload, {
    ...options,
    method: 'POST',
  })
}

export function get(url, options = {}) {
  return fetchJSON(url, undefined, {
    ...options,
    method: 'GET',
  })
}

export function login(username, password) {
  const creds = btoa(`${username}:${password}`)
  return post('/api/login', undefined, {
    headers: {
      Authorization: `Basic ${creds}`,
    },
  }).then(resp => {
    setToken(resp.token)
    return resp
  })
}

export function me() {
  return get('/api/me')
}

export function sendMessage(msg) {
  return post('/api/data/message', msg)
}

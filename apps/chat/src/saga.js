import { isString } from 'lodash'
import { eventChannel } from 'redux-saga'
import { spawn, call, put, take } from 'redux-saga/effects'
import { login, me, sendMessage } from './api'

import { setCurrentUser, pushMessage } from './actions'

export function* autoLogin() {
  yield call(login, 'sergeyt', 'sergeyt123')
  const user = yield call(me)
  yield put(setCurrentUser(user))
}

export function* send(msg) {
  yield call(sendMessage, msg)
}

function eventSourceChannel(source) {
  return eventChannel(emit => {
    source.onerror = err => {
      console.log('SSE error:', err)
      emit(new Error(err))
    }

    source.onmessage = event => {
      const msg = isString(event.data) ? JSON.parse(event.data) : event.data
      console.log('SSE message:', msg)
      emit(msg)
    }

    return () => source.close()
  })
}

function* eventListener() {
  const chan = eventSourceChannel(
    new EventSource('http://localhost:4302/api/event/stream')
  )
  while (true) {
    const msg = yield take(chan)
    if (msg instanceof Error) {
      // TODO handle error
      continue
    }
    if (msg.resource_type === 'message') {
      yield put(pushMessage(msg.result))
    }
  }
}

export default function* root() {
  yield spawn(autoLogin)
  yield spawn(eventListener)
}

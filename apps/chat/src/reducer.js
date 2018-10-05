import Immutable from 'immutable'
import { combineReducers } from 'redux'
import { handleActions } from 'redux-actions'

import { setCurrentUser, loadMessages, pushMessage } from './actions'

const GlobalState = Immutable.Record({
  currentUser: undefined,
})

const commonReducer = handleActions(
  {
    [setCurrentUser]: (state, action) => {
      return state.set('currentUser', action.payload)
    },
  },
  GlobalState()
)

const ChatState = Immutable.Record({
  messages: Immutable.List(),
})

const chatReducer = handleActions(
  {
    [loadMessages]: (state, action) => {
      return state.set('messages', action.payload || Immutable.List())
    },

    [pushMessage]: (state, action) => {
      const messages = state.messages.push(action.payload)
      return state.set('messages', messages)
    },
  },
  ChatState()
)

export default combineReducers({
  common: commonReducer,
  chat: chatReducer,
})

import { createAction } from 'redux-actions'

export const loadMessages = createAction('LOAD_MESSAGES') // set messages on load
export const pushMessage = createAction('PUSH_MESSAGE') // append message

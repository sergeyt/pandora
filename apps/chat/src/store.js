import { createStore, applyMiddleware, compose } from 'redux'
import createSagaMiddleware, { END } from 'redux-saga'

import rootReducer from './reducer'

const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose

export function makeStore() {
  const sagaMiddleware = createSagaMiddleware()

  const middlewares = [sagaMiddleware]

  const initialState = rootReducer({}, { type: '@@INIT' })

  const store = createStore(
    rootReducer,
    initialState,
    composeEnhancers(applyMiddleware(...middlewares))
  )

  store.runSaga = sagaMiddleware.run
  store.close = () => store.dispatch(END)

  return store
}

const store = makeStore()

export default store

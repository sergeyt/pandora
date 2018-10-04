import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'

import './index.css'
import store from './store'
import rootSaga from './saga'
import App from './App'

store.runSaga(rootSaga)

const Root = (
  <Provider store={store}>
    <App />
  </Provider>
)

ReactDOM.render(Root, document.getElementById('root'))

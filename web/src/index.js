import React from "react";
import ReactDOM from "react-dom";
import {applyMiddleware, compose as reduxCompose, createStore} from "redux";
import {Provider} from "react-redux";
import createSagaMiddleware from "redux-saga";
import {CssBaseline} from "@material-ui/core";
import {ThemeProvider} from "@material-ui/styles";

import theme from "./theme";
import {appReducer, appSaga, initialState} from "./application/state";
import Application from "./application";


const sagaMiddleware = createSagaMiddleware();
const compose = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || reduxCompose;


// Initialize application state
const store = createStore(
    appReducer,
    initialState,
    compose(applyMiddleware(sagaMiddleware))
);

sagaMiddleware.run(appSaga);


ReactDOM.render(
    <React.Fragment>
        <CssBaseline/>
        <ThemeProvider theme={theme}>
            <Provider store={store}>
                <Application/>
            </Provider>
        </ThemeProvider>
    </React.Fragment>,
    document.querySelector("#root")
);

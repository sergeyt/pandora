import React from "react";
import ReactDOM from "react-dom";
import {CssBaseline} from "@material-ui/core";
import {ThemeProvider} from "@material-ui/styles";
import Application from "./application";
import theme from "./theme";

ReactDOM.render(
    <React.Fragment>
        <CssBaseline/>
        <ThemeProvider theme={theme}>
            <Application/>
        </ThemeProvider>
    </React.Fragment>,
    document.querySelector("#root")
);

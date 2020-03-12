import React, {useState} from "react";
import {makeStyles} from "@material-ui/styles";
import {TextField, Typography} from "@material-ui/core";
import {useDispatch, useSelector} from "react-redux";
import {query} from "../../state";

const useStyles = makeStyles(() => ({
    container: {
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        height: "80%",
    },
    form: {
        width: "100%",
    },
}));

function QueryInput() {
    const classes = useStyles();
    const [queryString, setQueryString] = useState("");

    const dispatch = useDispatch();
    const loading = useSelector(state => state.search.loading);

    function handleSubmit(event) {
        dispatch(query(queryString));
        event.preventDefault();
    }

    return (
        <div className={classes.container}>
            <Typography component="h2" variant="h2" align="center">
                Open Pandora!
            </Typography>
            <form className={classes.form} onSubmit={handleSubmit}>
                <TextField
                    variant="outlined"
                    margin="normal"
                    fullWidth
                    label="Query Documents"
                    name="query"
                    autoComplete="Query Documents"
                    autoFocus
                    onChange={event => setQueryString(event.target.value)}
                    value={queryString}
                    disabled={loading}
                />
            </form>
        </div>
    );
}

export default QueryInput;
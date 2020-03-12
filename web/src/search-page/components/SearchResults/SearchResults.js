import React, {useState} from "react";
import {makeStyles} from "@material-ui/styles";
import {TextField} from "@material-ui/core";
import {useDispatch, useSelector} from "react-redux";
import {query} from "../../state";
import DocumentPreview from "../DocumentPreview";

const useStyles = makeStyles(theme => ({
    container: {
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        height: "100%",
        width: "100%",
    },
    form: {
        width: "100%",
    },
    results: {
        display: "flex",
        // justifyContent: 'flex-start',
        flexWrap: 'wrap',
        '& > *': {
            margin: theme.spacing(1),
        },
    }
}));

function SearchResults() {
    const classes = useStyles();
    const [queryString, setQueryString] = useState("");

    const dispatch = useDispatch();
    const loading = useSelector(state => state.search.loading);
    const documents = useSelector(state => state.search.documents);

    function handleSubmit(event) {
        dispatch(query(queryString));
        event.preventDefault();
    }

    return (
        <div className={classes.container}>
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
            <div className={classes.results}>
                {
                    documents.map(document => (
                        <DocumentPreview document={document}/>
                    ))
                }
            </div>
        </div>
    );
}

export default SearchResults;
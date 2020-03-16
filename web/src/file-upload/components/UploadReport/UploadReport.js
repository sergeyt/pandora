import React, {useEffect, useState} from "react";
import Snackbar from "@material-ui/core/Snackbar";
import Collapse from "@material-ui/core/Collapse";
import Paper from "@material-ui/core/Paper";
import {makeStyles} from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import FileList from "./FileList";
import Header from "./Header";
import {useDispatch, useSelector} from "react-redux";
import {cancelAll, cancelFile} from "../..";
import {uploadDone} from "../../state/reducers";
import Alert from "@material-ui/lab/Alert";
import Typography from "@material-ui/core/Typography";

const useStyles = makeStyles(theme => ({
    snackbar: {
        bottom: theme.spacing(2),
        right: theme.spacing(2),

    },
    header: {
        // backgroundColor: theme.palette.grey[900],
    },
    list: {
        minWidth: 300,
        maxHeight: "300px",
        overflow: "auto"
    },
    notify: {
        display: "flex",
    }
}));


function UploadReport() {
    const classes = useStyles();
    const dispatch = useDispatch();
    const files = useSelector(state => state.upload.files);

    const [notify, setNotify] = useState(false);
    const [shown, setShown] = useState(files.length > 0);
    const [expanded, setExpanded] = useState(true);
    const done = files.every(uploadDone);

    useEffect(() => {
        if (!done && files.length > 0) {
            // New files were submitted to upload
            setShown(true);
            setExpanded(true);
            setNotify(false);
        }
        if (done && !shown && files.length > 0) {
            // Progress popover is closed and
            // upload is finished in background
            setNotify(true);
        }
    }, [done]);

    const clearAll = () => {
        // clear state if done
        setNotify(false);
        setShown(false);
        dispatch(cancelAll());
    };

    const handleClose = () => {
        setShown(false);
        setNotify(false);
    };

    const handleShow = () => {
        setShown(true);
        setExpanded(true);
    };

    return (
        <React.Fragment>
            <Snackbar
                open={shown && files.length > 0}
                anchorOrigin={{vertical: "bottom", horizontal: "right"}}
                className={classes.snackbar}
            >
                <Paper elevation={4}>
                    <Header
                        count={files.length}
                        done={done}
                        expanded={expanded}
                        onCancel={clearAll}
                        onClose={handleClose}
                        onExpand={setExpanded}
                    />
                    <Collapse in={expanded}>
                        <FileList files={files} onCancel={file => dispatch(cancelFile(file))}/>
                    </Collapse>
                </Paper>
            </Snackbar>
            <Snackbar
                open={!shown && notify && files.length > 0}
                autoHideDuration={6000}
                message={`${files.length} uploads done`}
                anchorOrigin={{vertical: "bottom", horizontal: "right"}}
            >
                <Alert onClose={handleClose} severity="success">
                    <div className={classes.notify}>
                        <Typography>{files.length} uploads done</Typography>
                        <Button color="secondary" size="small" onClick={handleShow}>
                            Show
                        </Button>
                    </div>
                </Alert>
            </Snackbar>
        </React.Fragment>
    );
}

export default UploadReport;

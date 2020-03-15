import React, {useEffect, useState} from "react";
import PropTypes from "prop-types";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogTitle from "@material-ui/core/DialogTitle";
import {RootRef} from "@material-ui/core";
import {useDropzone} from "react-dropzone";
import {makeStyles} from "@material-ui/styles";
import Typography from "@material-ui/core/Typography";
import useMediaQuery from "@material-ui/core/useMediaQuery";
import {useTheme} from "@material-ui/core/styles";
import clsx from "clsx";
import UploadFileList from "./UploadFileList";

const useStyles = makeStyles(theme => ({
    dropzone: {
        minHeight: "100px",
        padding: theme.spacing(2),
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        borderWidth: theme.spacing(0.5),
        backgroundColor: theme.palette.grey[100],
    },
    active: {
        borderStyle: "dashed",
        borderColor: theme.palette.primary.main,
        backgroundColor: theme.palette.grey[200],
    },
    accept: {},
    reject: {}
}));


function UploadDialog(props) {
    const {
        open,
        onClose,
        onFileUpload,
    } = props;

    // Retrieve styles
    const classes = useStyles();
    const theme = useTheme();
    const fullScreen = useMediaQuery(theme.breakpoints.down("xs"));

    // Setup dialog state
    const [files, setFiles] = useState([]);
    const onDrop = (acceptedFiles) => {
        const paths = new Set();
        files.forEach(file => paths.add(file.path));
        setFiles(files.concat(acceptedFiles.filter(file => !paths.has(file.path))));
    };
    const onDelete = (file, i) => {
        const newFiles = [...files];
        newFiles.splice(i, 1);
        setFiles(newFiles);
    };

    // Clear dialog state on open/close
    useEffect(() => {
        setFiles([]);
    }, [open]);

    // Get dropzone properties
    const {
        getRootProps,
        getInputProps,
        isDragActive,
        isDragAccept,
        isDragReject
    } = useDropzone({onDrop});
    const {ref, ...rootProps} = getRootProps();
    const dropzoneClass = clsx(
        classes.dropzone,
        isDragActive && classes.active,
        isDragAccept && classes.accept,
        isDragReject && classes.reject
    );

    return (
        <Dialog
            fullScreen={fullScreen}
            open={open}
            onClose={onClose}
            aria-labelledby="responsive-dialog-title"
        >
            <DialogTitle id="responsive-dialog-title">Upload Files</DialogTitle>
            <DialogContent>
                <RootRef rootRef={ref}>
                    <div {...rootProps} className={dropzoneClass}>
                        <input {...getInputProps()}/>
                        <Typography align="center">
                            Drag files here, or click to select files
                        </Typography>
                    </div>
                </RootRef>
                <UploadFileList files={files} onDelete={onDelete}/>
            </DialogContent>
            <DialogActions>
                <Button autoFocus onClick={onClose} color="primary">
                    Cancel
                </Button>
                <Button onClick={() => onFileUpload(files)} color="primary">
                    Upload
                </Button>
            </DialogActions>
        </Dialog>
    );
}

UploadDialog.propTypes = {
    open: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    onFileUpload: PropTypes.func.isRequired,
};


export default UploadDialog;
import {ACTION_CANCEL_ALL, ACTION_CANCEL_FILE, ACTION_UPDATE_STATUS, ACTION_UPLOAD, UploadStatus} from "./actions";

export const initialState = {
    files: [],
};


export function uploadDone(file) {
    return file.status === UploadStatus.FAILURE || file.status === UploadStatus.SUCCESS;
}


export function uploadReducer(state = initialState, action) {
    switch (action.type) {
        case ACTION_UPLOAD: {
            let shouldClear = state.files.every(uploadDone);
            let files = (shouldClear ? [] : state.files).concat(action.files.map(file => ({
                path: file.path,
                status: UploadStatus.PENDING
            })));
            return {files};
        }
        case ACTION_CANCEL_ALL:
            return {files: []};
        case ACTION_CANCEL_FILE:
            return {files: state.files.filter(file => file.path !== action.file.path)};
        case ACTION_UPDATE_STATUS:
            return {
                files: state.files.map(file => ((file.path === action.file.path) ? action.file : file))
            };
        default:
            return state;
    }
}
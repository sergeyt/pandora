import {
    ACTION_CANCEL_ALL,
    ACTION_CANCEL_FILE,
    ACTION_UPDATE_STATUS,
    ACTION_UPLOAD,
    ACTION_UPLOAD_PROGRESS,
    UploadStatus
} from "./actions";

export const initialState = {
    files: [],
};


export function uploadDone(file) {
    return file.status === UploadStatus.FAILURE || file.status === UploadStatus.SUCCESS;
}


export function uploadReducer(state = initialState, action) {
    switch (action.type) {
        case ACTION_UPLOAD: {
            action.files.forEach(file => file.status = UploadStatus.PENDING);
            let shouldClear = state.files.every(uploadDone);
            let files = (shouldClear ? [] : state.files).concat(action.files);
            return {files};
        }
        case ACTION_CANCEL_ALL:
            return {files: []};
        case ACTION_CANCEL_FILE:
            return {files: state.files.filter(file => file.path !== action.file.path)};
        case ACTION_UPDATE_STATUS:
        case ACTION_UPLOAD_PROGRESS:
            return {
                files: state.files.map(file => {
                    if (file.path === action.file.path) {
                        return Object.assign({}, file, action.file);
                    }
                    return file
                })
            };
        default:
            return state;
    }
}
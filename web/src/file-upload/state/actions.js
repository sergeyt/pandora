export const ACTION_UPLOAD = "upload.UPLOAD";
export const ACTION_CANCEL_ALL = "upload.CANCEL_ALL";
export const ACTION_CANCEL_FILE = "upload.CANCEL_FILE";
export const ACTION_UPDATE_STATUS = "upload.UPDATE_STATUS";
export const ACTION_UPLOAD_PROGRESS = "upload.PROGRESS";

export const UploadStatus = {
    SUCCESS: "success",
    FAILURE: "failure",
    PENDING: "pending",
    ACTIVE: "active",
};


export function upload(files) {
    return {
        type: ACTION_UPLOAD,
        files
    };
}

export function cancelAll() {
    return {type: ACTION_CANCEL_ALL};
}


export function cancelFile(file) {
    return {type: ACTION_CANCEL_FILE, file};
}


export function updateStatus(file, status) {
    return {
        type: ACTION_UPDATE_STATUS,
        file: {
            path: file.path,
            status,
        }
    };
}

export function uploadSuccess(file) {
    return updateStatus(file, UploadStatus.SUCCESS);
}

export function uploadFailure(file) {
    return updateStatus(file, UploadStatus.FAILURE);
}

export function uploadProgress(file, progress) {
    return {
        type: ACTION_UPLOAD_PROGRESS,
        file: {
            path: file.path,
            progress,
        }
    };
}

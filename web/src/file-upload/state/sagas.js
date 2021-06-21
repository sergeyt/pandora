import {all, call, cancel, cancelled, fork, put, select, take, takeEvery} from "redux-saga/effects";
import {END, eventChannel} from "redux-saga";
import {
    ACTION_CANCEL_ALL,
    ACTION_CANCEL_FILE,
    ACTION_UPDATE_STATUS,
    ACTION_UPLOAD,
    updateStatus,
    uploadFailure,
    uploadProgress,
    UploadStatus,
    uploadSuccess
} from "./actions";
import Pandora from "../../server-api";


function makeUploadChannel({pandora, file}) {
    const uploadingOperation = Pandora.makeCancellableOperation();

    return eventChannel(emit => {
        const onProgress = (percent) => {
            emit(uploadProgress(file, percent));
        };
        pandora.uploadFile(file, {onProgress, cancel: uploadingOperation}).then(() => {
            emit(uploadSuccess(file));
            emit(END);
        }).catch((err) => {
            console.error(err);
            emit(uploadFailure(file));
            emit(END);
        });

        return () => {
            uploadingOperation.cancel("Cancelled by user");
        };
    });
}


export function* doUpload(file) {
    yield put(updateStatus(file, UploadStatus.ACTIVE));
    const uploading = yield call(makeUploadChannel, {pandora: Pandora, file});
    try {
        // Receive actions indicating uploading progress
        // and redirect them to the redux store...
        while (true) {
            const action = yield take(uploading);
            yield put(action);
        }
    } catch (err) {
        console.error(err);
    } finally {
        if (yield cancelled()) {
            // If cancel received simply close the channel
            // This will result in uploading cancellation.
            uploading.close();
        }
    }
}

function isCancelled(action, file) {
    return action.type === ACTION_CANCEL_ALL ||
        action.type === ACTION_CANCEL_FILE &&
        action.file.path === file.path;
}

function isDone(action, file) {
    return action.file.path === file.path &&
        action.type === ACTION_UPDATE_STATUS &&
        (action.file.status === UploadStatus.SUCCESS ||
            action.file.status === UploadStatus.FAILURE);
}

export function* manageUpload(file) {
    const task = yield fork(doUpload, file);

    // Monitor for cancellation until uploading complete.
    while (true) {
        // Receive uploading life-cycle action
        const action = yield take([
            ACTION_CANCEL_ALL,
            ACTION_CANCEL_FILE,
            ACTION_UPDATE_STATUS
        ]);

        // The operation is either cancelled or completed
        if (isCancelled(action, file)) {
            yield cancel(task);
            return;
        } else if (isDone(action, file)) {
            return;
        }

        // Skip if action is related to different file...
    }
}


export function* handleUploadQueue() {
    try {
        const files = yield select(state => state.upload.files);
        for (let file of files) {
            if (file.status === UploadStatus.ACTIVE) {
                // upload is already in progress
                return;
            }
            if (file.status === UploadStatus.PENDING) {
                yield all([manageUpload(file)]);
                return;
            }
        }
    } catch (err) {
        // TODO: implement appropriate error handling
        console.error(`handleUploadSaga error :${err}`);
    }
}

export function* uploadRootSaga() {
    yield takeEvery(ACTION_UPLOAD, handleUploadQueue);
    yield takeEvery(ACTION_UPDATE_STATUS, handleUploadQueue);
}

export default uploadRootSaga;

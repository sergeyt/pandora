import {all, call, cancelled, put, select, take, takeEvery} from "redux-saga/effects";
import {END, eventChannel} from 'redux-saga';
import {
    ACTION_UPDATE_STATUS,
    ACTION_UPLOAD,
    updateStatus,
    uploadFailure,
    uploadProgress,
    UploadStatus,
    uploadSuccess
} from "./actions";
import Pandora from "../../server-api";


function makeUploadChannel({pandora, file, cancel}) {
    return eventChannel(emit => {
        const onProgress = (percent) => {
            emit(uploadProgress(file, percent));
        };
        pandora.uploadFile(file, {onProgress, cancel}).then(() => {
            emit(uploadSuccess(file));
            emit(END);
        }).catch((err) => {
            console.error(err);
            emit(uploadFailure(file));
            emit(END);
        });

        return () => {
            cancel.cancel('cancelled');
        };
    });
}


export function* handleFileUploadSaga(file) {
    const uploading = Pandora.makeCancellableOperation();
    try {
        yield put(updateStatus(file, UploadStatus.ACTIVE));
        const channel = yield call(makeUploadChannel, {pandora: Pandora, file, cancel: uploading});

        while (true) {
            const action = yield take(channel);
            console.log(action);
            yield put(action);
        }
    } catch (err) {
        if (yield cancelled()) {
            uploading.cancel('cancelled');
        }
        console.error(err);
    }
}


export function* handleUploadUpdateSaga() {
    try {
        const files = yield select(state => state.upload.files);
        for (let file of files) {
            if (file.status === UploadStatus.ACTIVE) {
                // upload is already in progress
                return;
            }
            if (file.status === UploadStatus.PENDING) {
                yield all([handleFileUploadSaga(file)]);
                return;
            }
        }
    } catch (err) {
        // TODO: implement appropriate error handling
        console.error(`handleUploadSaga error :${err}`);
    }
}

export function* uploadRootSaga() {
    yield takeEvery(ACTION_UPLOAD, handleUploadUpdateSaga);
    yield takeEvery(ACTION_UPDATE_STATUS, handleUploadUpdateSaga);
}

export default uploadRootSaga;

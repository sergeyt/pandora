import {all, call, put, select, takeEvery} from "redux-saga/effects";
import {ACTION_UPDATE_STATUS, ACTION_UPLOAD, updateStatus, uploadFailure, UploadStatus, uploadSuccess} from "./actions";
import Pandora from "../../server-api";


export function* handleFileUploadSaga(file) {
    try {
        yield put(updateStatus(file, UploadStatus.ACTIVE));
        yield call(path => Pandora.uploadFile(path), file.path);
        yield put(uploadSuccess(file));
    } catch (err) {
        yield put(uploadFailure(file));
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

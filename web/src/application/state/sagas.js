import {all} from "redux-saga/effects";
import {searchSaga} from "../../search-page";
import {uploadRootSaga} from "../../file-upload";


export function* appSaga() {
    yield all([
        searchSaga(),
        uploadRootSaga(),
    ]);
}

export default appSaga;

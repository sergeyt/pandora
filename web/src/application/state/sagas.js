import {all} from "redux-saga/effects";
import searchSaga from "../../search-page/state/sagas";


export function* appSaga() {
    yield all([
        searchSaga()
    ]);
}

export default appSaga;

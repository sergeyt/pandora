import {combineReducers} from "redux";
import {initialState as searchInitialState, searchReducer,} from "../../search-page";
import {initialState as uploadInitialState, uploadReducer} from "../../file-upload";

export const initialState = {
    search: searchInitialState,
    upload: uploadInitialState,
};

export const appReducer = combineReducers({
    search: searchReducer,
    upload: uploadReducer,
});


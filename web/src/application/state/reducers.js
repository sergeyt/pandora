import {combineReducers} from "redux";
import {initialState as searchInitialState, searchReducer,} from "../../search-page/state";


export const initialState = {
    search: searchInitialState,
};

export const appReducer = combineReducers({
    search: searchReducer,
});


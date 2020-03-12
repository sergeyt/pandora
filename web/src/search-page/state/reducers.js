import {ACTION_QUERY, ACTION_QUERY_RESULTS} from "./actions";

export const initialState = {
    loading: false,
    documents: [],
    success: true,
};


export function searchReducer(state = initialState, action) {
    switch (action.type) {
    case ACTION_QUERY:
        return Object.assign({}, state, {loading: true});
    case ACTION_QUERY_RESULTS:
        return Object.assign({}, state, {
            loading: false,
            documents: action.documents,
            success: action.success,
        });
    default:
        return state;
    }
}

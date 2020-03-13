import {ACTION_CLEAR, ACTION_QUERY, ACTION_QUERY_RESULTS} from "./actions";

export const initialState = {
    loading: false,
    query: "",
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
            query: action.query,
            documents: action.documents,
            success: action.success,
        });
    case ACTION_CLEAR:
        return Object.assign({}, state, {documents: [], query: ""});
    default:
        return state;
    }
}

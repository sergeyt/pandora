import {delay, put, takeEvery} from "redux-saga/effects";
import {ACTION_QUERY, queryResults} from "../../search-page/state";

const randomText = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor 
incididunt ut labore et dolore magna aliqua. Rhoncus dolor purus non enim praesent 
elementum facilisis leo vel. Risus at ultrices mi tempus imperdiet. Semper risus 
in hendrerit gravida rutrum quisque non tellus. Convallis convallis tellus id 
interdum velit laoreet id donec ultrices. Odio morbi quis commodo odio aenean 
sed adipiscing. Amet nisl suscipit adipiscing bibendum est ultricies integer 
quis. Cursus euismod quis viverra nibh cras. Metus vulputate eu scelerisque 
felis imperdiet proin fermentum leo. Mauris commodo quis imperdiet massa tincidunt. 
`;

const stubDocuments = [
    {
        title: "Metus vulputate",
        size: "3M",
        preview: "https://images.unsplash.com/photo-1553463637-062ef4fe06da?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1350&q=80",
        tags: ["animals", "birds"],
        link: "https://images.unsplash.com/photo-1553463637-062ef4fe06da?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1350&q=80",
        previewText: randomText,
    },
    {
        title: "Lorem ipsum",
        size: "10GB",
        image: "https://images.unsplash.com/photo-1583762713699-7a6d1b8b6679?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=700&q=80",
        tags: ["image", "animals", "dogs"],
        link: "https://images.unsplash.com/photo-1583762713699-7a6d1b8b6679?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=700&q=80",
        previewText: randomText,
    }
];


export function* querySaga() {
    try {
        // Imitate server call
        yield delay(500 * Math.random() + 500);
        yield put(queryResults(stubDocuments, true));
    } catch (err) {
        yield put(queryResults([], false));
    }
}

export function* searchSaga() {
    yield takeEvery(ACTION_QUERY, querySaga);
}

export default searchSaga;

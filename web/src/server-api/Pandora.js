import sleep from 'sleep-promise'
import axios from "axios";

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
        date: "December 31, 2019",
        size: "3M",
        image: "https://images.unsplash.com/photo-1480555017593-6fb093e13f10?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=700&q=80",
        tags: ["animals", "birds"],
        link: "https://images.unsplash.com/photo-1480555017593-6fb093e13f10?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=700&q=80",
        previewText: randomText,
    },
    {
        title: "Lorem ipsum",
        date: "January 1, 2020",
        size: "10GB",
        image: "https://images.unsplash.com/photo-1583762713699-7a6d1b8b6679?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=700&q=80",
        tags: ["image", "animals", "dogs"],
        link: "https://images.unsplash.com/photo-1583762713699-7a6d1b8b6679?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=700&q=80",
        previewText: randomText,
    }
];

export class Pandora {
    constructor() {
        this.axios = axios.create({
            baseURL: '/api/',
            timeout: 10000,
            headers: {}
        })
    }

    async queryDocuments(queryString) {
        await sleep(500 + 500 * Math.random());
        return stubDocuments;
    }

    async uploadFile(file) {
        let data = new FormData();
        data.append('file', file);

        return await this.axios.post(`/file/${file.name}`, data, {
            onUploadProgress(progressEvent) {
                let percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                console.log(`${file.name}: ${percentCompleted}%`)
            }
        });
    }
}

export default new Pandora();
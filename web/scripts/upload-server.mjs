import DBG from "debug";
import fs from "fs";
import path from "path";
import lodash from "lodash";
import logger from "morgan";
import express from 'express';
import fileUpload from 'express-fileupload';

const debug = DBG("upload-server:server");
const directory = process.env.DIRECTORY || ".data";
const port = process.env.SERVER_PORT || 9123;

const app = express();

app.use(fileUpload({
    limits: { fileSize: 10 * 1024 * 1024 * 1024 }, // 10GB
}));
app.use(logger('dev', process.stdout));

app.post('/api/file/:name', function (req, res) {
    if (!req.files || Object.keys(req.files).length === 0) {
        return res.status(400).json({error: 'No files were uploaded.'});
    }

    const name = req.params.name;
    const file = req.files.file;

    // Use the mv() method to place the file somewhere on your server
    file.mv(path.join(directory, name), err => {
        if (err) {
            return res.status(500).json({error: err.message});
        }
        res.json({uid: lodash.uniqueId()});
    });
});


if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory);
}

app.listen(port, () => debug(`Upload server listening on port ${port}!`));
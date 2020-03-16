import DBG from "debug";
import fs from "fs";
import path from "path";
import lodash from "lodash";
import logger from "morgan";
import express from 'express';
import fileUpload from 'express-fileupload';

const debug = DBG("upload-server:server");
const error = DBG("upload-server:error");
const directory = process.env.DIRECTORY || ".data";
const port = process.env.SERVER_PORT || 9123;

const app = express();

app.use(fileUpload({}));
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
            return res.status(500).send(err);
        }
        res.json({uid: lodash.uniqueId()});
    });
});

app.use(function (err, req, res, next) {
    // set locals, only providing error in development
    res.locals.message = err.message;
    res.locals.error = req.app.get('env') === 'development' ? err : {};

    // render the error page
    res.status(err.status || 500);
    error(`${err.status || 500} ${err.message}`);
    next(error);
});


if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory);
}

app.listen(port, () => debug(`Upload server listening on port ${port}!`));
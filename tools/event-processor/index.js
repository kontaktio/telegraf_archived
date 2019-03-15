const express = require('express');
const bodyParser = require('body-parser');
const app = express();

const ApiCaller = require('./apiCaller');
const apiAddress = process.env.API_ADDRESS;
console.log(`Talking with API: ${apiAddress}`);
const apiCaller = new ApiCaller(apiAddress);

const TelegrafEmitter = require('./telegrafEmitter');
const emitter = new TelegrafEmitter('/tmp/telegraf.sock');

app.use(bodyParser.json());

app.post('/event/collect', async (req, res) => {
    let apiKey = req.get('Api-Key');
    if (apiKey === undefined) {
        res.status(401).end();
    }
    try {
        let companyId = await apiCaller.getCompanyId(apiKey);
        emitter.emit(companyId, req.body);
        res.status(202).end();
    } catch (e) {
        console.log(`${apiKey}: ${e.message}`);
        res.status(e.response.status).end();
    }
});

app.listen(8080, () => "Listening started");
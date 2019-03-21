const http = require('http');
const process = require('process');
const ApiCaller = require('./apiCaller');
const apiAddress = process.env.API_URL;
console.log(`Talking with API: ${apiAddress}`);
const apiCaller = new ApiCaller(apiAddress);

const TelegrafEmitter = require('./telegrafEmitter');
const emitter = new TelegrafEmitter('/tmp/telegraf.sock');

async function processRequest(req, res, requestContent) {
    console.log(requestContent);
    let apiKey = req.headers['api-key'];
    if (apiKey === undefined) {
        return;
    }
    try {
        let companyId = await apiCaller.getCompanyId(apiKey);
        emitter.emit(companyId, JSON.parse(requestContent));
    } catch (e) {
        console.log(`${apiKey}: ${e.message}`);
        res.status(e.response.status).end();
    }
}

function respond(res, code) {
    res.writeHead(code);
    res.end();
}

http.createServer((req, res) => {
    if(req.method === 'POST') {
        let requestContent = '';
        req.on('readable', () => {
            let read = req.read();
            if (read !== null) {
                requestContent += read;
            }

        });
        req.on('end', async () => {
            respond(res, 202);
            await processRequest(req, res, requestContent);
        });
        req.on('aborted', () => respond(res, 202));
        req.on('error', () => respond(res, 202));
    } else {
        respond(res, 202);
    }
}).listen(8080, () => {
    console.log("Listening started")
});
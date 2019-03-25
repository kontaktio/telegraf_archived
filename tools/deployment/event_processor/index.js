const http = require('http');
const process = require('process');
const ApiCaller = require('./apiCaller');

const apiAddress = process.env.API_URL;
console.log(`Talking with API: ${apiAddress}`);

if (process.argv.length !== 3) {
    throw new Error("Invalid arguments. Usage: node index.js /socket/address.sock")
}

const socketPath = process.argv[2];
console.log(`Socket address: ${socketPath}`);

const apiCaller = new ApiCaller(apiAddress);

const TelegrafEmitter = require('./telegrafEmitter');
const emitter = new TelegrafEmitter(socketPath); //from parameter

async function processRequest(req, res, requestContent) {
    let apiKey = req.headers['api-key'];
    if (apiKey === undefined) {
        return;
    }
    try {
        let companyId = await apiCaller.getCompanyId(apiKey);
        emitter.emit(companyId, JSON.parse(requestContent));
    } catch (e) {
        console.log(`${apiKey}: ${e.message}`);
        if(e.response) {
            respond(res, e.response.status);
        } else {
            respond(res, 500)
        }
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
        req.on('end', () => {
            respond(res, 202);
            processRequest(req, res, requestContent);
        });
        req.on('aborted', () => respond(res, 202));
        req.on('error', () => respond(res, 202));
    } else {
        respond(res, 202);
    }
}).listen(8080, () => {
    console.log("Listening started")
});
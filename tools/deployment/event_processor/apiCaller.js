const https = require('https');

module.exports = class ApiCaller {
    constructor(apiAddress) {
        this.baseUrl = apiAddress
        this.headers = {
            'Accept': 'application/vnd.com.kontakt+json;version=10',
            'X-Kontakt-Agent': 'web-panel-apicaller'
        };
        this.cache = {}
    }

    async getCompanyId(apiKey) {
        if (this.cache[apiKey]) {
            return this.cache[apiKey]
        }
        console.log(`Querying for apiKey: ${apiKey}`);
        const options = {
            headers: {
                ...this.headers,
                'Api-Key': apiKey
            }
        };
        return new Promise((resolve, reject) => {
            https.get(this.baseUrl + '/manager/me', options, (res) => {
                if (res.statusCode !== 200) {
                    reject(new Error(`${res.statusCode} response from API`));
                    return;
                }
                let result = '';
                res.on('data', d => result += d);
                res.on('end', () => {
                    const jsonResult = JSON.parse(result);
                    const companyId = jsonResult.companyId;
                    this.cache[apiKey] = companyId;
                    resolve(companyId);
                })
            }).on('error', e => reject(e))
        });
    }
};
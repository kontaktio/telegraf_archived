const https = require('https');

module.exports = class ApiCaller {
    constructor(apiAddress) {
        this.baseUrl = apiAddress;
        this.headers = {
            'Accept': 'application/vnd.com.kontakt+json;version=10',
            'X-Kontakt-Agent': 'web-panel-apicaller'
        };
        this.cache = {}
    }

    async getCompanyId(apiKey) {
        const cache = this.cache;
        if (cache[apiKey]) {
            return cache[apiKey]
        }
        console.log(`Querying for apiKey: ${apiKey}`);

        const url = this.baseUrl + '/manager/me';
        const options = {
            headers: {
                ...this.headers,
                'Api-Key': apiKey
            }
        };
        return new Promise((resolve, reject) => {
            try {
                https.get(url, options, res => {
                    if (res.statusCode !== 200) {
                        console.log(res);
                        reject(new Error(`${res.statusCode} response from API`));
                        return;
                    }
                    let result = '';
                    res.on('data', d => result += d);
                    res.on('end', () => {
                        const jsonResult = JSON.parse(result);
                        const companyId = jsonResult.companyId;
                        cache[apiKey] = companyId;
                        resolve(companyId.split('-').pop());
                    })
                }).on('error', (e) => {
                    reject(e);
                });
            } catch (e) {
                reject(e);
            }
        });
    }
};
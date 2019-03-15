const axios = require('axios');

module.exports = class ApiCaller {
    constructor(apiAddress) {
        this.axios = axios.create({
            baseURL: apiAddress,
            timeout: 5000,
            headers: {
                'Accept': 'application/vnd.com.kontakt+json;version=10',
                'X-Kontakt-Agent': 'web-panel-apicaller'
            }
        });
        this.cache = {}
    }

    async getCompanyId(apiKey) {
        if (this.cache[apiKey]) {
            return this.cache[apiKey]
        }
        console.log(`Querying for apiKey: ${apiKey}`);
        let res = await this.axios.get('/manager/me', {headers: {'Api-Key': apiKey}});
        console.log(`Got response: ${JSON.stringify(res.data)}`);
        let companyId = res.data.companyId;
        this.cache[apiKey] = companyId;
        return companyId;
    }
};
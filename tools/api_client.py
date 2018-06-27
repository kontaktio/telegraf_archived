import requests


class ApiClient:
    ACCEPT = 'application/vnd.com.kontakt+json;version=10'
    GET_DEVICES_PATH = '/device'
    GET_MANAGER_ME_PATH = '/manager/me'
    MODELS_NOT_SUPPORTING_TLM = ['SMART_BEACON', 'CARD_BEACON', 'USB_BEACON', 'SENSOR_BEACON']

    def __init__(self, api_url, api_key):
        self._api_url = api_url
        self._api_key = api_key

    def get_telemetry_unique_ids(self, api_venue_id=None):
        params = {
            'access': 'OWNER',
            'deviceType': 'BEACON'
        }
        if api_venue_id is not None:
            params['q'] = 'venue.id==%s' % api_venue_id

        result = self._get_collection(self.GET_DEVICES_PATH, params, 'devices')
        self._filter_telemetry_not_compatible(result)
        unique_ids = [r['uniqueId'] for r in result]
        print "Received uniqueIds: %s" % str(unique_ids)
        return unique_ids

    def get_company_id(self):
        response = requests.get(self._api_url + self.GET_MANAGER_ME_PATH, headers=self._get_headers()).json()
        company_id = str(response['company']['id'])
        print "Received companyId %s" % company_id
        return company_id


    def _filter_telemetry_not_compatible(self, collection):
        return [b for b in collection if b['model'] not in self.MODELS_NOT_SUPPORTING_TLM]

    def _get_headers(self):
        return {
            'Api-Key': self._api_key,
            'Accept': self.ACCEPT
        }

    def _get_collection(self, path, params, collection_name):
        result = []
        params['startIndex'] = 0
        params['maxResult'] = 500
        while True:
            response = requests.get(self._api_url + path, headers=self._get_headers(), params=params).json()
            for r in response[collection_name]:
                result.append(r)
            if response['searchMeta']['nextResults'] == '':
                break

            params['startIndex'] = params['startIndex'] + 500
        return result

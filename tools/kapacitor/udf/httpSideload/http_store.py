import os
from cachetools import cached, TTLCache
import requests


class HttpStore:

    def __init__(self):
        self._session = requests.Session()
        self._session.headers.update({'Api-Key': os.environ.get('API_KEY')})
        self._baseUrl = os.environ.get('API_URL')

    @cached(cache=TTLCache(maxsize=10000, ttl=60*15))
    def call_get(self, url):
        r = self._session.get(self._baseUrl + url)
        if r.ok:
            return r.json()
        else:
            r.raise_for_status()

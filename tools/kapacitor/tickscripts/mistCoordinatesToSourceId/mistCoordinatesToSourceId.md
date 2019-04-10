## How to configure Apache Proxy to cache requests for (lat, lng) to gateway mapping (/api/management/zone/find)

 1. Run `a2enmod cache cache_disk`
 2. Edit apache conf file and add:
    ```    
    <Location "/api/management/zone/find">
        ProxyPass https://<<BA_URL>>/api/management/zone/find retry=0
        ProxyPassReverse http://<<BA_URL>>/api/management/zone/find
    </Location>
    
    CacheEnable disk /
    CacheIgnoreNoLastMod On
    ```
 3. Instead of using @sideLoad directly to ba-api, make request to http://localhost/api/management/zone/find
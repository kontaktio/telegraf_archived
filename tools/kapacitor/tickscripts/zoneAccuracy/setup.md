## How to setup a Zone level accuracy

1. Add UDF to Kapacitor config
    ```
    [udf.functions.httpSideload]
       prog = "/etc/kapacitor/udf/http-sideload/kapacitor_udf/venv/bin/python"
       args = ["-u", "/etc/kapacitor/udf/http-sideload/kapacitor_udf/kapacitor_agent.py"]
       timeout = "10s"
       [udf.functions.httpSideload.env]
            API_URL = "http://localhost:8087"
            API_KEY = "SET_ME_MANUALLY"
    ```

2. Copy ``location.tick`` and ``udf/httpSideload`` to correct location on kapacitor server.
3. Register ``dwellTime-tpl`` task definition on for kapacitor instance.
    ```
    kapacitor define-template location-tpl -tick location.tick
    ```
    
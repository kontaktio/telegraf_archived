## How to setup a Zone level accuracy

1. Add UDF to Kapacitor config
    ```
    [udf.functions.httpSideload]
       prog = "/etc/kapacitor/udf/httpSideload/venv/bin/python"
       args = ["-u", "/etc/kapacitor/udf/httpSideload/kapacitor_agent.py"]
       timeout = "10s"
       [udf.functions.httpSideload.env]
            API_URL = "http://localhost:8087"
            API_KEY = "SET_ME_MANUALLY"
    ```

2. Clone/update telegraf to server. Enter telegraf dir.

3. Update tick scripts on server.
    ```bash
    cd /home/ec2-user/telegraf/tools/kapacitor/
    sudo cp -r * /etc/kapacitor
    ```
4. Install python libraries.
    ```bash
    cd /etc/kapacitor/udf/httpSideload
    sudo virtualenv venv
    sudo venv/bin/pip install .
    ```
    
3. Register ``dwellTime-tpl`` task definition on for kapacitor instance.
    ```
    kapacitor define-template location-tpl -tick /etc/kapacitor/tickscripts/zoneAccuracy/location.tick
    ```
    
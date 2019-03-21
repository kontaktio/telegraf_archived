## How to setup a Single Wire Click endpoint

1. Configure kapacitor.conf
    Add section for single wire http endpoint. Fill correct values for url.
    ```
    [[httppost]]
      endpoint = "singleWireNotification"
      url = "https://api.icmobile.singlewire.com/api/v1/notifications"
      headers = { Content-Type = "application/json;charset=UTF-8" }
      row-template-file = "/etc/kapacitor/singleWireClick-httppost.tmpl"
    ```
2. Copy ``singleWireClick-httppost.tmpl`` to correct location on kapacitor server.
3. Register ``singleWireClick-tpl`` task definition on for kapacitor instance.
    ```
    kapacitor define-template singleWireClick-tpl -tick singleWireClick-tpl.tick
    ``` 
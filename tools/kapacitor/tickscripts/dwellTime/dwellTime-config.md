## How to setup a Dwell Time report

1. Configure kapacitor.conf
    Add section for dwellTime http endpoint. Fill correct values for url and security header.
    ```
    [[httppost]]
      endpoint = "dwellTime"
      url = "https://testba-api.kontakt.io/api/reports/input/dwellTime"
      headers = { Api-Key = "" }
      row-template-file = "/etc/kapacitor/dwellTime-httppost.tmpl"
    ```
2. Copy ``dwellTime-httppost.tmpl`` to correct location on kapacitor server.
3. Register ``dwellTime-tpl`` task definition on for kapacitor instance.
    ```
    kapacitor define-template dwellTime-tpl -tick dwellTime-tpl.tick
    ``` 
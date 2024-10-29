import json
import re

import requests

resp = requests.post("http://192.168.123.161:10089/get_account_token",
                     json={
                         "account": "pmjb8ecp@vlook.cloud",
                         "scope": "audience:server:client_id:342877434807-ug3d7kknubivam3e8m52uc6heocs0g3a.apps.googleusercontent.com",
                         "pkg_name": "com.instagram.lite",
                     })
print(resp.text)
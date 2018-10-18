#!/usr/local/bin/python3
import json
import sys
import urllib.parse
import urllib.request


f = urllib.request.urlopen('http://127.0.0.1:8080/sources/datadog/?%s' % urllib.parse.urlencode({'query': sys.argv[1]}))
data = json.loads(f.read().decode('utf-8'))

items = []
for r in data['items']:
    attrs = r['attributes']
    items.append({
        "title": r['autocomplete'],
        "arg": attrs['url'],
        "autocomplete": r['autocomplete'],
        "subtitle": "Open Datadog dashboard %s" % r['autocomplete'],
        "icon": {
            "path": "icon.png"
        }
    })


print(json.dumps({"items": items}, indent=2))

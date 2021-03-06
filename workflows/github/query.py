#!/usr/local/bin/python3
import json
import sys
import urllib.parse
import urllib.request


f = urllib.request.urlopen('http://127.0.0.1:8080/sources/github/?%s' % urllib.parse.urlencode({'query': sys.argv[1]}))
data = json.loads(f.read().decode('utf-8'))

items = []
for r in data['items']:
    attrs = r['attributes']
    items.append({
        "title": attrs['fullname'],
        "arg": "https://github.com/%s" % attrs['fullname'],
        "autocomplete": r['autocomplete'],
        "subtitle": "Open Github Repository %s" % attrs['fullname'],
        "icon": {
            "path": "github.png"
        }
    })


print(json.dumps({"items": items}, indent=2))

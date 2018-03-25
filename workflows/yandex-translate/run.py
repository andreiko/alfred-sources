#!/usr/local/bin/python3
import json
import sys
import re
import urllib.parse
import urllib.request

query = sys.argv[1]

if re.match('^[\sa-zA-Z]+$', query):
    lang = 'en-ru'
else:
    lang = 'ru-en'

API_KEY = open('./api.key').read()

f = urllib.request.urlopen('https://dictionary.yandex.net/api/v1/dicservice.json/lookup?%s' % urllib.parse.urlencode({
    'key': API_KEY,
    'lang': lang,
    'text': query,
}))

data = json.loads(f.read().decode('utf-8'))

items = []
for d in data['def']:
    for t in d['tr']:
        syns = []
        means = []

        syns.append(t['text'])

        for syn in t.get('syn', []):
            syns.append(syn['text'])

        for mean in t.get('mean', []):
            means.append(mean['text'])

        items.append({
            "title": ", ".join(syns),
            "subtitle": "(%s) %s" % (t['pos'], ", ".join(means)),
            "arg": query,
        })

if not items:
    items.append({
        "title": "no translations found",
        "valid": False
    })

print(json.dumps({'items': items}))


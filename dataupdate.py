import requests as rq
import pandas as pd
import codecs as cd
from git import Repo

url = 'https://docs.google.com/spreadsheets/d/1gBftBUJqezG2QjUij7YBoZYQswpGYPYFmX9V3bfZues/gviz/tq?tqx=out:csv&sheet=Glossary'
r = rq.get(url, allow_redirects=True)
open('data.csv', 'wb').write(r.content)

with open('data.csv', 'r', encoding='utf8') as fin:
    data = fin.read().splitlines(True)
with open('data.csv', 'w', encoding='utf8') as fout:
    fout.write('devn,iast,part,gend,engl,pada\n')
    fout.writelines(data[2:])

df = pd.read_csv(cd.open('data.csv', 'r', 'utf-8'))
print(df.to_json(r'data.json', force_ascii=False, orient='records', lines=True))

with open('data.json', 'r', encoding='utf8') as fin:
    data = fin.read().splitlines(True)
len = len(df) - 1
with open('data.json', 'w', encoding='utf8') as fout:
    fout.write('[\n')
    for idx, line in enumerate(data):
        if idx == len:
            fout.write(line)
        else:
            fout.write(line.replace('\n', ',').strip() + '\n')
    fout.write(']\n')

rp = Repo('.')
rp.index.add(['data.json'])
rp.index.commit('update data')
origin = rp.remote('origin')
origin.push()



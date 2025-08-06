import requests as rq
import pandas as pd
import codecs as cd
from git import Repo

url = 'https://docs.google.com/spreadsheets/d/1jewXALlsOojcegEiGAWYvncWaUy1g8oRE5Oup68bfP0/gviz/tq?tqx=out:csv&sheet=Glossary'
r = rq.get(url, allow_redirects=True)
open('data.csv', 'wb').write(r.content)



new_columns = [
    'artist', 'song', 'userdifficulty', 'usernotes', 'mmmdifficulty',
    'mmmnotes', 'style', 'learn', 'tuning', 'mmmtutorial', 'timesignature'
]
df = pd.read_csv('data.csv', header=0)
print(f"Detected columns: {list(df.columns)}")
print(f"Number of columns: {len(df.columns)}")

df = df.iloc[:, :11]  # Select only the first 11 columns
df.columns = new_columns


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

#rp = Repo('.')
#rp.index.add(['data.json'])
#rp.index.commit('update data')
#origin = rp.remote('origin')
#origin.push()



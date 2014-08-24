# Work in progress
#	contains regexp to be used in receiving text messages from Twilio

import csv
tsvfile = open('usasexguide.tsv','rb')
tsvin = csv.reader(tsvfile, delimiter='\t’)
csvfile = open('price_date.csv','wb’)
csvout = csv.writer(csvfile)
import string
for row in tsvin:
	if string.count(row, '$') > 0:
		csvout.writerows(row)

 for row in tsvin:

...     if string.count(row[7], '$') > 0:

...             prices.append(row)

import re
pex = [re.findall('\$\d+\.?\d*', x[7]) for x in prices]
pex = [re.findall('\$\d+\.?\d*', x[7]) for x in prices if re.findall('\$\d+\.?\d*', x[7])]
pex_flat = [x for y in pex for x in y]
from re import sub
>>> from decimal import Decimal
>>> for x in pex_flat:
...     sum += Decimal(sub(r'[^\d.]','',x))

sum / Decimal(len(pex_flat))

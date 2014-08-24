import csv
import string

with open('usasexguide.tsv', 'rb') as tsvin, open('usasexguide_match.csv', 'wb') as csvout:
	tsvin = csv.reader(tsvin, delimiter='\t')
	csvout = csv.writer(csvout)
	#Example code to loop through entire data set and rewrite to a new file if a string is matched.
	for row in tsvin:
		if string.count(row[7], ’t') > 0:
			csvout.writerows(row)

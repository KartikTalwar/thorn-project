import json
import re
from collections import Counter

# Count the frequency of occurence of words in the freenet data set
words = re.findall(r'\w+', open('freenet_ascii.json').read().lower())
common_words = Counter(words).most_common()

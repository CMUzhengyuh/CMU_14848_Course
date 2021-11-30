from pyspark import SparkContext, SparkConf
import json

conf = SparkConf()
sc = SparkContext.getOrCreate(conf=conf)

stopWordFile = "/tmp/stopWord.txt"
stop_rdd = sc.textFile(stopWordFile)

signs = list("?!.,[]\t\():;")

###
file_dir = '/HW5/*/*'
output_dir = "output.txt"
rdd = sc.wholeTextFiles(file_dir)
stop_rdd = sc.textFile(stopWordFile)
stop_words = stop_rdd.map(lambda n : n.strip()).collect()

###
for sign in signs:
	rdd = rdd.replace(sign, " ")
 
output = rdd.flatMap(lambda content: ((word, [content[0]]) for word in content[1].lower())).filter(lambda w: w[0] not in stop_words).reduceByKey(lambda m,n: m+n).map(lambda w: format(w))

def format(item):
    key = item[0]
    files = item[1]
    word_freq = {}
    # Loop all possible files for all possible key words
    for source_file in files:
        if source_file in word_freq:
            word_freq[source_file] += 1
        else:
            word_freq[source_file] = 1
    value = []
    # Format the word frequency 
    for source_file in word_freq:
        value.append((source_file, word_freq[source_file]))
    return (key, value)

###
spark_output = output.collect()
spark_result = json.dumps(spark_output)

###
f = open(output_dir, "w")
f.write(spark_result)
f.close()



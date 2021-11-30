# CMU_14848_Course
Repo for CMU 14848 Cloud Infrastructure

Student: Zhengyu HU

Andrew ID: zhengyuh

Time: 2021 Fall

---

### Homework 2

Part A:

1. URL for your Docker image: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/hello-world/general

2.  Screenshot for the execution of your docker container on GCP:
    
    ![avatar](HW2/HW2/HW2_PartA.png)

3.  Dockerfile contents and the source code file:

    DockerFile:

    ```Dockerfile
    FROM openjdk:11
    RUN mkdir /app
    COPY out/production/HW2/ /app
    WORKDIR /app
    CMD java HelloWorld
    ```

    Source Code:
    ```Java
    public class HelloWorld {
        public static void main(String[] args) {
            int count = 0;
            try {
                while (true) {
                    Thread.sleep(2 * 1000);
                    System.out.println("Hello World from Docker! (zhengyuh)");
                }
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        }
    }

    ```

Part B: Screenshot for running Jupyter notebook:

![avatar](HW2/HW2/HW2_PartB.jpg)

---

### Homework 3

See: HW3/NoSQL

Source Code:
```Python
import boto3
import csv

# Part I: Create an S3 Instance

s3 = boto3.resource('s3', aws_access_key_id='AKIA2MPFUMSSFF5OALFU', aws_secret_access_key='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')
try:
    s3.create_bucket(Bucket='14848-bucket', CreateBucketConfiguration={'LocationConstraint': 'us-east-2'})
except Exception as e:
    print("-- You have already built the bucket: 14848-bucket")
bucket = s3.Bucket("14848-bucket")
bucket.Acl().put(ACL='public-read')

# Part II: Create the Dynamo Table

body = open('/Users/hobo/desktop/14848/HW3/exp1.csv', 'rb')
o = s3.Object('14848-bucket', 'test').put(Body=body)
dyndb = boto3.resource('dynamodb', region_name='us-east-2', aws_access_key_id='AKIA2MPFUMSSFF5OALFU', aws_secret_access_key='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')
try:
    table = dyndb.create_table(
        TableName='DataTable',
        KeySchema=[
            {
                'AttributeName': 'PartitionKey',
                'KeyType': 'HASH'
            },
            {
                'AttributeName': 'RowKey',
                'KeyType': 'RANGE'
            }
        ],
        AttributeDefinitions=[
            {
                'AttributeName': 'PartitionKey',
                'AttributeType': 'S'
            },
            {
                'AttributeName': 'RowKey',
                'AttributeType': 'S'
            },
        ],
        ProvisionedThroughput={
            'ReadCapacityUnits': 5,
            'WriteCapacityUnits': 5
        }
    )
except Exception as e:
    print("-- Table already exists")

table = dyndb.Table("DataTable")

table.meta.client.get_waiter('table_exists').wait(TableName='DataTable')
print(table.item_count)

# Part III: Read data from the database

with open('/Users/hobo/desktop/14848/HW3/experiments.csv', 'r') as csvfile:
    csvf = csv.reader(csvfile, delimiter=',', quotechar='|')
    next(csvf)
    for item in csvf:
        print(item)
        body = open('/Users/hobo/desktop/14848/HW3/'+item[4], 'rb')
        s3.Object('14848-bucket', item[4]).put(Body=body)
        md = s3.Object('14848-bucket', item[4]).Acl().put(ACL='public-read')
        url = " https://s3-us-west-2.amazonaws.com/14848-bucket/"+item[4]
        metadata_item = {'PartitionKey': item[4], 'RowKey': item[0], 'Temp': item[1],
                 'Conductivity': item[2], 'Concentration': item[3], 'url': url}
        try:
            table.put_item(Item=metadata_item)
        except:
            print("item may already be there or another failure")


response1 = table.get_item(
    Key={
        'PartitionKey': 'exp1.csv',
        'RowKey': '1'
    } 
)
response2 = table.get_item(
    Key={
        'PartitionKey': 'exp2.csv',
        'RowKey': '2'
    } 
)
response3 = table.get_item(
    Key={
        'PartitionKey': 'exp3.csv',
        'RowKey': '3'
    } 
)
print("Item result:")
print(response1['Item'])
print(response2['Item'])
print(response3['Item'])
print("Response is (Take Item 1 as an example):")
print(response1)


```
---

### Homework 4

See: HW4

1. Mapper source code:
``` Python
import sys

for line in sys.stdin:
    line = line.strip()
    temperature = int(line[87:92])
    q = int(line[92])
    if ((temperature != 9999)) and (q == 0 or q == 1 or q == 4 or q == 5 or q == 9) == True:
        print('%s\t%d' % (line[15:23], int(line[87:92])))

```

2. Reducer source code:
```Python
from operator import itemgetter
import sys

current_date = None
current_temperature = 0
date = None

for line in sys.stdin:
    line = line.strip()
    date, temperature = line.split('\t', 1)
    try:
        temperature = int(temperature)
    except ValueError:
        continue

    if current_date == date:
        if temperature > current_temperature:
            current_temperature = temperature
    else:
        if current_date:
            print('%s\t%d' % (current_date, current_temperature))
        current_temperature = temperature
        current_date = date

if current_date == date:
    print('%s\t%d' % (current_date, current_temperature))
```

3. Screenshot of Hadoop MapReduce Job in the terminal:
   
    ![avatar](HW4/GCP_MapReduce1.png)
    ![avatar](HW4/GCP_MapReduce2.png)

4. Output file of results: <br/>
    See: HW4/maximumTemperatureByDay

---

### Homework 5

See: HW5

1. Source code for Spark Application:

```Python
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
```

2. Spark Result: HW5/output.txt

3. Screenshot for Jupyter Notebook and Spark Container:

![avatar](HW5/Jupyter-Notebook.png)
![avatar](HW5/Spark_Container.png)

---

### Project - Checkpoint - Option 1

Part A: Source code for the main terminal application:

```Java
import java.util.Scanner;

public class main {

    public static void main (String[] args) {

        System.out.println("Welcome to Big Data Processing ToolBox");
        System.out.println("Author: Zhengyu HU  Andrew ID: zhengyuh  Course: 14848");
        System.out.println("Please select the application to run");
        System.out.println("1. Jupyter Notebook");
        System.out.println("2. Apache Hadoop");
        System.out.println("3. Apache Spark");
        System.out.println("4. SonarQube & SonarScanner");
        System.out.println("Type the index of application (no-index input to stop)");
        String originalIP = "";

        Scanner scan = new Scanner(System.in);
        boolean isQuit = false;
        while (scan.hasNext()) {
            String str = scan.next();
            switch(str) {
                case "1":
                    String ip1 = originalIP + "";
                    System.out.println("Jupyter Notebook URL is " + ip1);
                    break;
                case "2":
                    String ip2 = originalIP + "";
                    System.out.println("Apache Hadoop URL is " + ip2);
                    break;
                case "3":
                    String ip3 = originalIP + "";
                    System.out.println("Apache Spark URL is " + ip3);
                    break;
                case "4":
                    String ip4 = originalIP + "";
                    System.out.println("SonarQube & SonarScanner URL is " + ip4);
                    break;
                default:
                    System.out.println("Quit the ToolBox");
                    isQuit = true;
                    break;
            }
            if (isQuit)
                break;
        }
        scan.close();
    }
}

```

![avatar](Project/Terminal_Demo_Result.png)

Part B: Docker images of applications:

1. URL for driver: <br/>
   <https://hub.docker.com/repository/docker/hobo965859229/my-driver>

2. URL for Jupyter Notebook: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/my-jupyter-notebook

3. URL for Apache Hadoop: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/my-hadoop

4. URL for Apache Spark: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/my-spark

5. URL for SonarQube & SonarScanner: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/my-sonarqube

Part C: Screenshot for the Kubernetes Engine with the containers running

![avatar](Project/Screenshot/K8s_Result.png)

Part D: Steps to run Docker images on Kubernetes Engine

1. Driver for the application (Local Test): <br/>
   $ javac main.java <br/>
   $ java main <br/>
   Function: Read shell input and map to corresponding URL

2. Docker image Test: Complete Dockerfile for driver and 4 applications

3. Built docker image: <br/>
   $ docker build -t *Image-name*

4. Push docker image to my own dockerhub: <br/>
   $ docker push hobo965859229/ *Image-name*

5. Complete my- *Application-name* -deployment.yaml and my- *Application-name* -serveice.yaml

6. Open local Kubernete clusters then deploy images: <br/>
   $ kubectl apply -f my- *Application-name* -deployment.yaml <br/>
   $ kubectl apply -f my- *Application-name* -service.yaml

---

### Project - Final - Option 1

Part A: Demo Vedio:

Youtube: https://www.youtube.com/watch?v=m1OZ6U2O3ZE

Part B: Screen shot for applications (Please check the consistency of IP)

1. Local Client:

![avatar](Project/Screenshot/Client_Interface.png)

2. Apache Hadoop: 35.202.46.187:6766

![avatar](Project/Screenshot/Hadoop.png)

3. Apache Spark: 35.202.46.187:6866 

![avatar](Project/Screenshot/Spark.png)

4. SonarQube: 35.202.46.187:6966

![avatar](Project/Screenshot/SonarQube.png)

5. Jupyter Notebook: 35.202.46.187:6070 

![avatar](Project/Screenshot/Jupyter.png)

Part C: Screen shot for Google Cloud Platform - Kubernetes

1. Kubernetes - Cluster

![avatar](Project/Screenshot/K8s_Shell.png)

2. Kubernetes - Service & Ingreess

![avatar](Project/Screenshot/ip.png)

3. Kubernetes - Workload

![avatar](Project/Screenshot/K8s_Workloads.png)

4. Cloud Shell

![avatar](Project/Screenshot/K8s_Get.png)

Part D: Source code 

1. Source code for the client terminal applications

```Go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// sendRequest Send requests to driver
func sendRequest(conn net.Conn, text string) {
	conn.Write([]byte(text))
}

// receiveRespond Receive requests from driver
func receiveRespond(conn net.Conn) string{
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
		return ""
	}
	fmt.Println(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))
	return string(buffer[:n])
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c","start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func terminal(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	for {
		fmt.Println("Welcome to Big Data Processing Application")
		fmt.Println("Please type the number that corresponds to which application you would like to run:")
		fmt.Println("1. Apache Hadoop")
		fmt.Println("2. Apache Spark")
		fmt.Println("3. Jupyter Notebook")
		fmt.Println("4. SonarQube and SonarScanner")
		fmt.Println("Type the number here >")

		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		sendRequest(conn,text)
		website := receiveRespond(conn)
		open(website)
	}
}

func main() {
	// Hard code GCP external IPs
	server := "35.202.46.187:6666"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	terminal(conn)

}
```

2. Source code for the server of driver

```Go
package main

import (
    "net/http/httputil"
	"os"
    "fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Create a new proxy
func newProxy(targetHost string) *httputil.ReverseProxy {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil
	}
	// Create a proxy to convert the internal port of K8s to external IPs
	return httputil.NewSingleHostReverseProxy(url)
}

// Establish the proxy server
func establishProxy(){
	// Internal Ports defined by YAML
	JupyterUrl := os.Getenv("JUPYTER_NOTEBOOK") + ":8888"
	SonarQubeUrl  := os.Getenv("SONARQUBE") + ":9000"
	SparkUrl  := os.Getenv("SPARK") + ":8080"
	HadoopUrl := os.Getenv("HADOOP") + ":50070"

	go proxyServer(HadoopUrl,":6766")
	go proxyServer(SparkUrl,":6866")
	go proxyServer(JupyterUrl,":6966")
	go proxyServer(SonarQubeUrl,":6070")

}

// Start service
func proxyServer(proxyUrl string, addr string) {
	_ = http.ListenAndServe(addr, newProxy(proxyUrl))
}

// Send back URL to clients
func sendRequest(conn net.Conn, url string) {
	conn.Write([]byte(url))
	fmt.Printf("Send back URL: %s\n", url)
}

func main() {

	// Establish the socket
	addr,_:=net.ResolveTCPAddr("tcp4", ":6666")

	netListen, err := net.ListenTCP("tcp", addr)
	CheckError(err)
	defer netListen.Close()

	// Establish the proxy server
	establishProxy()
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		// Handlle new connection
		handleConnection(conn)
	}
}

// Hard code of expose GCP internal IPs to external URLs
func handleConnection(conn net.Conn) {
	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		receiveData := string(buffer[:n])
		switch receiveData {
		// Format: "http:// <GCP_IPS> : <PORTS>"
		case "1":
			sendRequest(conn,"http://35.202.46.187:6766")
		case "2":
			sendRequest(conn,"http://35.202.46.187:6866")
		case "3":
			sendRequest(conn,"http://35.202.46.187:6966")
		case "4":
			sendRequest(conn,"http://35.202.46.187:6070")
		}
	}
}
```


Part E: Docker images of applications:

1. URL for driver: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/driver

2. URL for Jupyter Notebook: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/jupyter-notebook

3. URL for Apache Hadoop: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/hadoop

4. URL for Apache Spark: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/spark

5. URL for SonarQube & SonarScanner: <br/>
   https://hub.docker.com/repository/docker/hobo965859229/sonarqube

Part F: Steps to run the whole project

1. Find 4 available applications and pack them to Docker images (See Part E) <br/>
    $ docker build -t *Image-name* <br/>
    $ docker push hobo965859229/ *Image-name* : v1.0

2. Prepare deployment and service YAML files for GCP Kubernetes (See Project/cfg)
    
3. Construct the driver based on a Golang container (See Project/Terminal) <br/>
    $ docker build -t hobo965859229/driver:v1.0 <br/>
    $ docker push hobo965859229/driver:v1.0

![avatar](Project/Screenshot/Driver_Image.png)

4. Now we can see all docker images run locally:
   
![avatar](Project/Screenshot/Local_Image.png)

5. Upload all local files to GitHub: https://github.com/CMUzhengyuh/CMU_14848_Course

6. Afterwards, we can start deploy all the services and deployment to GCP

7. Create a GCP Kubernetes cluster with the following configuration:

![avatar](Project/Screenshot/Final_Create_K8s_Cluster.png)
![avatar](Project/Screenshot/Final_Create_K8s_Cluster_Configuration.png)

8. Open the Cloud shell and go into the Kubernete cluster

![avatar](Project/Screenshot/K8s_Shell.png)

9. Git clone contents from my public repository: <br/>
    $ git clone https://github.com/CMUzhengyuh/CMU_14848_Course.git

![avatar](Project/Screenshot/Git.png)

10. Set the directory to Project/cfg <br/>
    $ cd CMU_14848_Course/Project/cfg

11. Deploy all the YAML file on Kubernetes <br/>
    -deployment.yaml: $ kubectl apply -f *Application-name* -deployment.yaml <br/>
    -service.yaml: $ kubectl apply -f *Application-name* -service.yaml <br/>
    hadoop.yaml: $ kubectl create -f hadoop.yaml

12. After Step 11 GCP shell would be like this:

![avatar](Project/Screenshot/deploy.png)

13. Check the Service & Ingress section of Kubernetes Engine

![avatar](Project/Screenshot/K8s_Service.png)

14. Update the exposed IPs of section-driver-service to local server.go

15. Write hard code of IP:35.202.46.187 to server.go

16. Rebuild the image and update the version: <br/>
    $ docker build -t hobo965859229/driver:v2.0 <br/>
    $ docker push hobo965859229/driver:v2.0

17. Re-deploy the service of driver so that Kubernetes can run the latest driver <br/>
    $ kubectl apply -f driver-service.yaml

18. Keep GCP Kubernetes running and open the client in local shell <br/>
    $ ~/CMU_14848_Course/Project/UI

19. Build the client shell locally <br/>
    $ go build client.go

![avatar](Project/Screenshot/Client_Build.png)

20. Open the client Unix file and start experience the tool

![avatar](Project/Screenshot/Client_Interface.png)
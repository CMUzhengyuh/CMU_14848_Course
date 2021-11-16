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

### Homework 3

See: HW4

1. Mapper source code
``` Python
import sys

for line in sys.stdin:
    line = line.strip()
    temperature = int(line[87:92])
    q = int(line[92])
    if ((temperature != 9999)) and (q == 0 or q == 1 or q == 4 or q == 5 or q == 9) == True:
        print('%s\t%d' % (line[15:23], int(line[87:92])))

```

2. Reducer source code
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

3. Screenshot of Hadoop MapReduce Job in the terminal
   
    ![avatar](HW4/GCP_MapReduce1.png)
    ![avatar](HW4/GCP_MapReduce2.png)

4. Output file of results:
    See: HW4/maximumTemperatureByDay

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

![avatar](Project/K8s_Result.png)

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



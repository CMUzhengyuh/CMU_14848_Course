# CMU_14848_Course
Repo for CMU 14848 Cloud Infrastructure

Student: Zhengyu HU

Andrew ID: zhengyuh

Time: 2021 Fall

---

### Homework 2

Part A:

1. URL for your Docker image: 
https://hub.docker.com/repository/docker/hobo965859229/hello-world/general

2.  Screenshot for the execution of your docker container on GCP:
    
    See: HW2/HW2/HW2_PartA

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

See: HW2/HW2/HW2_PartB

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

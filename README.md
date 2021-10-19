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

### Homework 2

See: HW3/NoSQL
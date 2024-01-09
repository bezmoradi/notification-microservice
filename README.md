# How to Use The Notification Microservice

This project is intended to be used in tandem with the [Knowledge Base Microservice](https://github.com/bezmoradi/knowledge-base-microservice)

## Intentions

I use both this project and the [Knowledge Base Microservice](https://github.com/bezmoradi/knowledge-base-microservice) on a **daily** basis to review tips about the programming languages, frameworks, libraries, and especially insights from books I've previously read. By review, I mean reviewing the notes I have taken, the snippets which I have found useful, and all in all the tips which help me do my day-to-day job faster.

## How Do They Work?

The way that I have created this personal **Knowledge Base** is that I use a GitHub repository like [Tutorials](https://github.com/bezmoradi/tutorials) repo of mine by adding all my notes in a markdown format then it's [Knowledge Base Microservice](https://github.com/bezmoradi/knowledge-base-microservice)'s responsibility to access that repo, fetching one document randomly then pass it to Kafka (Of course the number of documents that can be fetched is configurable).  
Then it's this repository's job to listen to Kafka for new events and as soon as a new one is fired (Let's say every 6 hours or so), that document will be reformatted to HTML then sent to my personal email (The email settings is also configurable via `.env` files).

## Kafka Installation

As these two microservices are loosely-coupled and Event-Driven Architecture (EDA) is used, you need to make sure an up and running instance of Kafka is at your disposal. For Kafka installation, we are going to use Docker. First, create a Zookeeper container:

```text
$ docker run -d --name zookeeper -p 2181:2181 wurstmeister/zookeeper
```

Next, create a Kafka container and link it to Zookeeper:

```text
$ docker run -d --name kafka -p 9092:9092 --link zookeeper:zookeeper -e KAFKA_ADVERTISED_HOST_NAME=localhost -e KAFKA_ADVERTISED_PORT=9092 -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 wurstmeister/kafka
```

Do `docker ps` to make sure both services are running. If so, Kafka is now available via `localhost:9092` URI (This Kafka setup is useful if you want to run the Go binary **directly** without using Docker. For a Dockerized app, keep on reading to the last section on this file).

## Installation

Clone this repo by the following command:

```text
$ git clone git@github.com:bezmoradi/notification-microservice.git
```

Make a copy of the `.env.example` file and call it `.env` which includes the following keys:

```text
MAIL_SERVICE_HOST=
MAIL_SERVICE_USER=
MAIL_SERVICE_PASS=
MAIL_SERVICE_PORT=

KAFKA_BROKER=
```

Usually the default value for `KAFKA_BROKER` is `localhost:9092` and for the mail service, any third-party service can be used; I personally prefer to use Gmail simply because there isn't goring to be any crazy load on the SMTP server to get blocked by Google (Also it's free to use).  
To use Gmail, on the top right corner of your Gmail account inside the browser, click on you avatar then hit "Manage your Google Account". Next, on the left sidebar click on "Security" tab then under "Hot you sing in to Google" section, click on "2-Step Verification". Google will ask you to enter you phone number then sends you a verification code. After entering the code in the field provided, you're all set. Then inside that "2-Step Verification" page, scroll all the way to the bottom of the page to see "App passwords". Click on it and create one. The newly-generated password should look something like `dprh jyip jhrj obzh`. The important point here is to remove all white spaces to be like `dprhjyipjhrjobzh`. The completed `.env` should be something like this one:

```text
MAIL_SERVICE_HOST=smtp.gmail.com
MAIL_SERVICE_USER=<your_handle>@gmail.com
MAIL_SERVICE_PASS=dprhjyipjhrjobzh
MAIL_SERVICE_PORT=465

KAFKA_BROKER=localhost:9092
```

-   Now you can run the app using `go run main.go` command.

## How to Create A Docker Container

A `Dockerfile` is also included in this repo for those (including myself) who like to run the app via Docker. First create an image:

```text
$ docker build -t notification-microservice-image .
```
As the `.env` file is ignored inside the `.dockerignore` file, while creating a new container we have to pass it to docker:  

```text
$ docker run -d --env-file ./.env notification-microservice-image
```

## How to Connect The Dockerized App to Dockerized Kafka 
As said before, if you are going to run this microservice on Docker, the problem you'd face is that the Dockerized app cannot access Dockerized Kafka and the reason being is that by setting `KAFKA_BROKER=localhost:9092`, within our app's container `localhost` will point to the container itself, not the container machine where Kafka is running that's why our app cannot access it and we will get the following error:

```text
we would get the `Connection error: connect ECONNREFUSED 127.0.0.1:9092` error.
```

To resolve this, you need to update the value of `KAFKA_BROKER` in order to connect to the Kafka container's hostname or IP address from within the Docker network. As the first step, use the container name as the hostname:

```text
#KAFKA_BROKER=localhost:9092
KAFKA_BROKER=kafka:9092
```

`kafka` is any name while we create a Dockerized Kafka container. Technically speaking, Docker allows containers within the same network to communicate using container names.  
Next, Ensure all containers are on the same network; in other words, confirm that your Go container, Zookeeper, and the Kafka container are on the same Docker network. You can create a custom Docker network and attach all containers to it to allow them to communicate easily:

```text
$ docker network create <network_name>
```

Let's first delete all containers we created in the previous steps and repeat the process once more. For Zookeeper we have:

```text
$ docker run --network my-network -d --name zookeeper -p 2181:2181 wurstmeister/zookeeper
```

And for Kafka we have:

```text
$ docker run --network my-network -d --name kafka -p 9092:9092 --link zookeeper:zookeeper -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092 -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 wurstmeister/kafka
```
- Finally create a container for our Go app:  

```text
$ docker run -d --network my-network --env-file ./.env notification-microservice-image
```

Or if you want to keep the existing containers, you can add them to a network as follows:

```text
$ docker network connect <network_name> <container_name>
```

To see the list of containers belonging to a network we have:

```text
$ docker network inspect <network_name>
```

When Kafka runs inside a Docker container, it might have its listeners configured to accept connections only from localhost (`localhost:9092`). However, when trying to connect from another container within the same Docker network, Kafka needs to allow connections from external sources, not just `localhost`. That's why we need the below flag:

```text
-e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092
```

The above flag, defines the network addresses that Kafka binds to and allows Kafka to listen on all network interfaces (`0.0.0.0`) for incoming connections on port `9092`.

```text
-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092
```

The above configuration also tells clients how to connect to Kafka. It's the address to which clients should connect, and this is what's advertised to them. Here, `kafka` is the hostname within the Docker network (a.k.a name of the container), and it advertises itself on port `9092`. By configuring Kafka to listen on all interfaces (`0.0.0.0`) and advertise its address correctly within the Docker network (using its hostname `kafka` or its IP address within the network), you make communication from other containers within the same Docker network possible.
If you create a Dockerized version of this repo, [Knowledge Base Microservice](https://github.com/bezmoradi/knowledge-base-microservice), Zookeeper, and also Kafka without any issues, the same setup can be used for a prod environment which I will discuss next.

## CI/CD & Deployment

I use [GitHub Actions](https://docs.github.com/en/actions) for CI/CD purposes and also deploying these microservices to DigitalOcean (No need to say that any other cloud provider can also be used). Inside [.github/workflows/ci-cd.yml](https://github.com/bezmoradi/notification-microservice/blob/master/.github/workflows/ci-cd.yml) file, I have defined a simple workflow which:

1. Clone the repo
2. Login to Docker Hub
3. Build the Go image & push it to Docker Hub
4. SSH to DigitalOcean Droplet
5. Pull down the image and create a container off of it

I feel like it's straight-forward and I need to clarify two things. There are five secrets as follows:

```text
secrets.DOCKERHUB_USERNAME
secrets.DOCKERHUB_PASSWORD
secrets.DIGITALOCEAN_SERVER_IPV4
secrets.DIGITALOCEAN_SERVER_USER
secrets.SSH_PRIVATE_KEY
```
To create these credentials, go to `https://github.com/<your_username>/<repo_name>/settings` then click on `https://github.com/<your_username>/<repo_name>/settings/secrets/actions` and hit the "New repository secret" as many time as needed.  
All these secrets are self-explanatory except `secrets.SSH_PRIVATE_KEY`. I have used the `appleboy/ssh-action@master` package for Actions and the way it works is that you need to create a private SSH key on you DigitalOcean Droplet and pass that key to it. 

Colegas
=======

Colegas is a µservices/nannoservices proof of concept


## Goals


## Description


## Requirements

  - Linux
  - **Go 1.6+**, to build de µservices and *cli* commands
  - **Docker** and **docker compose**, to start and stop the middleware infrastructure services 


## How to run

### Build

Build the person µservice/nannoservice and the CLI

    $ cd persond && ./build && cd ..
    $ cd person && ./build && cd ..


Build the person application

    $ cd personapp/
    $ mvn clean package


### Start the system

From **project root** First start the middleware services needed (the *command bus* and the *event bus*). After that we can start each µservice

    $ # Start middleware infrastructure
    $ docker-compose up -d
    
    $ # Start person µservice in verbose mode and send output to persond.log file
    $ nohup ./persond/persond -verbose > ./persond.log 2>&1 &

    $ # Start person application
    $ nohup java -jar personapp/target/personservice.jar  > ./personapp.log 2>&1 &


At this point we can send request to the person service via the CLI...

    $ ./person/person jomoespe


... or use the application.

    $ curl http://localhost:8080/person/jomoespe



>  Note: At this moment the ''person'' service is a dummy which only manages the person ids: *illescas*, *juergas*, *jbbarquero*, *dparra* and *jomoespe*, and only return the person name.


### Stop the system

    $ docker-compose kill

// TODO document how to stop the persond service

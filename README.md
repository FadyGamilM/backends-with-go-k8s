# backends-with-go-k8s

# Microservices we have : 
- Counter Microservice 
    - keep produces numbers, sends the number to another microservice which stores it in redis 

- Server Microservice 
    - receives number from Counter microservice and stores it in redis 

- Poller Microservice 
    - reads numbers from server 

# To run the microservices : 
1- First we need to build the images then run instance of each image 
> for counter microservice 
```cmd
docker build -f ./counter/Dockerfile -t counter-service .
```
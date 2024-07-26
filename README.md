# backends-with-go-k8s

# Microservices we have : 
- Counter Microservice 
    - keep produces numbers, sends the number to another microservice which stores it in redis 

- Server Microservice 
    - receives number from Counter microservice and stores it in redis 

- Poller Microservice 
    - reads numbers from server 
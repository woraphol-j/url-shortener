[![CircleCI](https://circleci.com/gh/woraphol-j/url-shortener/tree/master.svg?style=svg)](https://circleci.com/gh/woraphol-j/url-shortener/tree/master)
## URL Shortener

### About the project
This project is neatly designed and built to demonstrate the use of Golang and its ecosystem for building a simple yet functional URL Shortener service.

### Design decision
- Instead of building things from scratch, I adopt go-kit which is a microservice framework that facilitates implementing microservice with various aspects such as logging or instrumentation being already provided. Besides, Go-kit enforces a solid convention that makes the code consistent and highly testable. The rest of the code follows this same convention even if they are not part of gokit components.
- It uses go-fmt for linting and formating.
- Even though Go claims it has a fully functional testing framework and some people do not even recommend using any framework other than the built-in one, I still decided to choose ginkgo as the testing framework as among other things it is the BDD framework that I am very familir with from my Node.js background. In addition, I also select `Gomega` for assertion and `gomock` to do mocking.
- To ensure highest code quality, I also set up the project so that every commit pushed into the repository will be tested in CircleCI. For steps of execution, please have a look.
at `./circleci/config.yml`.
The CI page is here https://circleci.com/gh/woraphol-j/url-shortener/tree/master

## Getting Started
### Prerequisite
1. Make sure you have docker and docker-compose installed.
2. Clone the project.
3. Run the following command to prepare the project.
```bash
make prepare
```
4. Run the following command at the top level of the project in order to get the service and its dependency up and running:
```bash
docker-compose up -d
```
5. You can test the service by using the following cURL script:
```bash
curl -d "{\"url\": \"http://www.medium.com\"}" -X POST http://localhost:8080/shorturls
```
this will return the shortend url in JSON format. To convert it back and redirect to the original url you can put the shortUrl in the
browser address bar and enter, or use the returned short url in the following curl script to perform the test:
```bash
curl -v http://localhost:8080/T5yrYS3kX
```
Note that the above url is just an example. The one you use must be the one returned in the first step.

6. You can inspect the mongo database by going to the following url
```
http://localhost:8081/
```
7. You can monitor the service by going to the following prometheus dashboard
```
http://localhost:9090/
```

### Test
Run the following command to execute both unit and integration test:
```bash
make test
```
Note that the integration test needs Mongo to be up and running. which is already provided in the Makefile script.

### Room for improvement
Due to the time constraint, this project is quickly built to demonstrate the use of Golang and its ecosystem to develop a url shortener service. Although the project is production-ready and fully functional, there are still some parts that can be improved as follows:
#### 1. Better short url code generation
  - Currently, the system uses `shortid` (https://github.com/teris-io/shortid) for short url code generation. It has some pros and cons but for the sake of simplicity and for this demo, it was eventually chosen. My research reveals that there are a few other options to achieve it as well such as base62, hasid etc. Perhaps, in the future, I will try these alternatives out. Despite that, please note that the current code is designed to make that library switching relatively easy as I code it against the interface so that when change actually occurs, the business logic in the service layer does not need to be aware of it at all.

#### 2. Use centralized logging system
  - The service is designed to follow the `XI. Logs` item of 12 factor app guideline in that log should be written to the console and it is up to the deployment to decide where the log should be collected and processed. In the future, I may consider using tools such as firebeat to redirect log from the container to ELK stack, for example.

#### 3. Add more test coverage
  - Currently there are 2 test files where are `service_test.go` for unit test and `transport_test.go` for integration test. With more time, I would add
  more test coverage.



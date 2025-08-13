# OpenTelemetry + zipkin, jeager, prometheus, grafana

- service A -> service B -> service C

- to run applications, service A: 8080 and service B: 8081

  ```bash
    docker compose build --build-arg GITHUB_TOKEN=<token_to_private_repos>
    docker compose up -d
  ```

- to run request: POST to <http://locahost:8080>

 ```bash
  curl  'http://localhost:8080/hello'
 ```

 ou

 ```bash
  go run test/main.go
 ```

- to see zpkin:
  - <http://localhost:9411/>
- to see jeager:
  - <http://localhost:16686/>
- to see prometheus:
  - <http://localhost:9090/>
- to see grafana
  - <http://localhost:3000>

## todo

- logs

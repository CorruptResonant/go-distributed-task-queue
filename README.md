GDTQ: High-Performance Distributed Task Queue
GDTQ is an infrastructure-layer system built in Go designed to handle asynchronous task processing using a decoupled architecture. It utilizes Redis as a persistent message broker and features a concurrent worker pool with graceful shutdown capabilities.

Core Engineering Features
Decoupled Architecture: Separates the Task Producer (REST API) from the Task Consumers (Worker Pool) using Redis. This allows for independent scaling of each layer.

Concurrency and Scalability: Implements a Goroutine-based worker pool that can be horizontally scaled to process multiple concurrent asynchronous tasks.

Persistence and Reliability: Uses Redis LPUSH and BRPOP operations to ensure that tasks are persisted and survive application-layer crashes.

Graceful Shutdown: Orchestrates system exits using sync.WaitGroup and os.Signal handling. The system ensures in-flight tasks are completed before workers exit to prevent data loss.

Containerization: Fully dockerized environment using multi-stage builds and Docker Compose for one-click deployment and environment consistency.

Real-time Observability: Features a custom metrics endpoint and a dashboard to monitor queue depth and system status in real-time.

Technical Stack
Language: Go 1.25 (Standard Library, Context, Sync)

Broker: Redis (via go-redis/v9)

Infrastructure: Docker, Docker Compose

API: RESTful JSON API

Frontend: Vanilla JavaScript and HTML5

Architecture Diagram
Code snippet

graph LR
    User((User)) -->|POST /enqueue| API[Go API Server]
    API -->|LPUSH| Redis[(Redis Broker)]
    Redis -->|BRPOP| W1[Worker 1]
    Redis -->|BRPOP| W2[Worker 2]
    Redis -->|BRPOP| W3[Worker 3]
    W1 -->|Log| Console
    W2 -->|Log| Console
    W3 -->|Log| Console

    
Getting Started
Prerequisites
Docker and Docker Compose

Go 1.25 or higher (for local development)

Redis (for local development)

Installation and Execution (Docker Compose)
Clone the repository

Bash

git clone https://github.com/CorruptResonant/go-distributed-task-queue.git
cd go-distributed-task-queue
Launch the system

Bash

docker-compose up --build
Open the dashboard Open the frontend/index.html file in a web browser.

Testing
Run the unit test suite to verify task serialization and model integrity:

Bash

go test ./internal/models/...
System Demonstration
1. Concurrency Test
Use the Load Test button on the dashboard to inject 50 tasks. Observe the logs to see Workers 1, 2, and 3 processing tasks in parallel. The dashboard queue depth will decrement as workers finish their processing blocks.

2. Graceful Shutdown Test
While tasks are being processed, trigger a Ctrl+C in the terminal. The API stops accepting new tasks immediately, but the workers finish their current processing block before logging a stopping message and exiting the process. This demonstrates the use of Go Context and WaitGroups to maintain data integrity.

Future Roadmap
At-Least-Once Delivery: Implementation of the RPOPLPUSH pattern for task acknowledgment.

Dead Letter Queue: Automatic handling of failed tasks after a set number of retries.

Dynamic Worker Scaling: Adjusting worker pool size based on real-time queue depth.

Created by Sujal Kapoor
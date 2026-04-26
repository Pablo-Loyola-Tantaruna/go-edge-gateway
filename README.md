# GopherGuard Edge Gateway

**GopherGuard** is a high-performance Layer 7 API Gateway and Load Balancer built in **Go**. Designed with distributed systems principles, it provides an ultra-lightweight solution for traffic management, observability, and resilience in cloud-native architectures.

## 🚀 Key Features

* **Atomic Round Robin Balancing:** High-concurrency load balancing implementation using atomic operations (`sync/atomic`) to eliminate lock contention and ensure sub-millisecond overhead.
* **Active Health Monitoring:** A background monitoring engine that performs active TCP probing to proactively detect node failures and eject them from the traffic pool before they impact the end user.
* **Cloud-Native Observability:** Native **Prometheus** integration. It exposes critical performance metrics, including request throughput and p99 latency histograms, for real-time monitoring via Grafana.
* **Resilient Retries:** Intelligent retry logic. If a node fails during a request, the Gateway automatically reroutes the traffic to a healthy node, maximizing application availability.
* **Structured Telemetry:** Integrated logging middleware providing full traceability: Request IDs, methods, paths, and precise processing duration.
* **YAML-Driven Configuration:** Decoupled architecture allowing for seamless backend and server management without binary recompilation.

## 🛠️ Tech Stack

* **Language:** Go (Golang) 1.21+
* **Metrics:** Prometheus Client Library
* **Configuration:** YAML v3
* **Deployment:** Docker & Kubernetes Ready

## 📋 Quick Start

### 1. Requirements
* Go 1.21 or higher.
* Prometheus (optional, for metrics visualization).

### 2. Configuration
Define your microservices in the `config.yaml` file:

```yaml
server:
  port: 8080
  health_check_interval: "20s"

backends:
  - url: "http://localhost:8081"
  - url: "http://localhost:8082"
```

### 3. Running the Gateway
Bash
#### Install dependencies
```
go mod tidy
```
### Run the server
```
go run cmd/server/main.go
```
## 📊 Monitoring & Observability
The gateway exposes industry-standard metrics at:

**Endpoint:** http://localhost:8080/metrics

You can plug this directly into a Prometheus instance to visualize throughput and error rates.

## 🧠 Architecture Decision Records (ADR)
**Why Go?** Chosen for its native concurrency model (Goroutines) and its ability to handle thousands of persistent connections with a negligible memory footprint compared to JVM-based solutions.

**Zero-Dependency Core:** The proxy engine leverages Go’s standard library to ensure maximum compatibility, security, and minimal "bloat," avoiding heavy third-party frameworks.

**Thread Safety: All state** management within the ServerPool utilizes sync.RWMutex and sync/atomic to guarantee data integrity under high-concurrency scenarios.

### Developed by Pablo Loyola 
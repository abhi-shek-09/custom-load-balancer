# Go Load Balancer Project
A modern, flexible load balancer supporting both Layer 4 (TCP) and Layer 7 (HTTP/HTTPS) modes, with multiple routing strategies, health checks, TLS termination, retry logic, logging, IP-filtering and rate limiting.

## Features
#### Layer 4 TCP Proxying (Round-Robin)
Proxies raw TCP connections to backend servers using a simple round-robin algorithm.

#### Layer 7 HTTP Reverse Proxying
Handles full HTTP requests, forwards them to appropriate backend servers, and streams the response.

#### Health Checks
Continuously pings backend servers to ensure traffic is only routed to healthy instances.

#### Multiple Load Balancing Strategies
Supports Round-Robin, Random selection, and Least Connections algorithms for smarter traffic routing.

#### TLS Termination
Terminates HTTPS connections at the load balancer using configurable TLS certificates.

#### Retry Logic
Automatically retries failed HTTP requests (on 5xx responses) to alternate healthy backends.

#### Timeouts (Server and Backend)
Enforces request read/write and backend connection timeouts to protect against slow or stuck clients and servers.

#### IP Filtering (Whitelist/Blacklist via CIDR)
Blocks or allows requests based on configurable IP ranges using CIDR notation.

#### Rate Limiting (Per-IP)
Throttles requests from individual clients using a token bucket mechanism to prevent abuse.

#### Blocked Path Filtering
Instantly blocks access to sensitive or internal endpoints by matching request paths.

#### Automatic HTTP→HTTPS Redirection
Redirects all plain HTTP traffic to HTTPS using a parallel redirect server.

#### Custom Logging
Logs every request with detailed info: method, backend, status, retry attempts, TLS details, and client IP.

#### Configurable via config.yaml
All backends, modes, strategies, limits, CIDRs, and certificates can be defined in a central configuration file.



## How It Works
- Initialization
  - Loads configuration parameters (backends, strategy, limits, TLS, etc.) from config.yaml.
  - Initializes the backend pool with health states and prepares the strategy manager.
  - Sets the desired load balancing strategy (e.g., round-robin, random, least connections) at startup.
- Health Monitoring
  - Continuously checks backend availability every 5 seconds using TCP probes.
  - Automatically marks unhealthy backends and excludes them from routing decisions.
  - Logs health status transitions (UP → DOWN, DOWN → UP) for debugging and visibility.
- Routing
  - L4 Mode: Proxies raw TCP connections to healthy backend servers using round-robin logic.
  - L7 Mode: Handles full HTTP traffic — parses requests, applies middleware (rate limit, IP filtering, etc.), and forwards them via reverse proxy.
  - Supports automatic retry on backend failure, and enforces read/write timeouts.
- Security
  - Blocks unauthorized IPs using CIDR-based allowlist defined in the config.
  - Implements per-client IP rate limiting using a token bucket algorithm.
  - Optionally redirects HTTP traffic to HTTPS with TLS termination at the load balancer.
- Logging and Error Handling
  - Logs every incoming request with method, backend used, status, retry attempts, client IP, and TLS details (if applicable).
  - Logs internal errors such as backend failures, configuration issues, and panic recoveries.
  - Provides visibility into load balancing behavior and backend utilization patterns.

## Setup
#### 1. Clone the Repo
```
git clone https://github.com/abhi-shek-09/load-balancer.git
cd load-balancer
```
#### 2. Install Dependencies
This is a pure Go project. Just make sure Go modules are enabled:
```
go mod tidy
```
#### 3. Generate TLS Cert (Development Only)
```
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```
This generates cert.pem and key.pem for TLS termination.

#### 4. Start the Load Balancer
``` 
go run main.go --config=config.yaml 
```
Or build and run the binary.

#### 5. Start all the backends

L4
```
nc -lk 9001
```
and repeat for the rest

L7
``` 
curl -k https://localhost:8443
```

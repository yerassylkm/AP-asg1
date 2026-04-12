Project Structure & Repository Links
The project is split into three separate repositories to follow microservices best practices:

: Contains the logic for Order and Payment services, and the Docker infrastructure.

: The "Source of Truth" for all gRPC definitions.

: Contains the pre-compiled Go code (stubs) generated from the protos, imported by both services as a Go module.

```mermaid
graph TD
    Client[Postman / Client] -- HTTP/JSON (Port 8080) --> OrderService[Order Service]
    
    subgraph "Internal Network"
        OrderService -- gRPC Call (Port 50051) --> PaymentService[Payment Service]
    end

    OrderService -- SQL --> OrderDB[(PostgreSQL: order_db)]
    PaymentService -- SQL --> PaymentDB[(PostgreSQL: payment_db)]
```

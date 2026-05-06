# Assignment 3: Event-Driven Architecture

This project implements an event-driven architecture using RabbitMQ as the message broker. The system consists of three microservices: Order Service, Payment Service, and Notification Service.

## Architecture Overview

```
Order Service -> Payment Service (gRPC/HTTP) -> RabbitMQ -> Notification Service
```

- **Order Service**: Creates orders and initiates payments.
- **Payment Service**: Processes payments and publishes events to RabbitMQ.
- **Notification Service**: Consumes payment events and sends notifications.

## Services

### Order Service
- **Port**: 8080
- **Endpoints**:
  - `POST /orders` - Create order
  - `GET /orders/:id` - Get order by ID
  - `GET /orders?customer_id=...` - Get orders by customer
  - `PATCH /orders/:id/cancel` - Cancel order

### Payment Service
- **Ports**: 8081 (HTTP), 50051 (gRPC)
- **Endpoints**:
  - `POST /payments` - Process payment
  - `GET /payments/:order_id` - Get payment by order ID

### Notification Service
- **No public ports** (internal consumer)
- Listens to `payment.completed` queue

## Message Broker

- **RabbitMQ**: Used for event publishing/subscribing
- **Queue**: `payment.completed` (durable)
- **Message Format**: JSON
- **Payload**:
  ```json
  {
    "order_id": "uuid",
    "amount": 10000,
    "customer_email": "user@example.com",
    "status": "Authorized"
  }
  ```

## Key Features

### Manual ACKs
Messages are acknowledged only after successful processing in Notification Service.

### Durability
Queues are declared as durable to survive broker restarts.

### Idempotency
Notification Service uses PostgreSQL to track processed message IDs, preventing duplicate notifications.

### Graceful Shutdown
All services handle SIGINT/SIGTERM for clean shutdown of DB and broker connections.

## Running the System

```bash
docker-compose up --build
```

## Testing

1. Create an order via Order Service:
   ```bash
   curl -X POST http://localhost:8080/orders \
     -H "Content-Type: application/json" \
     -d '{"customer_id": "123", "customer_email": "user@example.com", "item_name": "Test Item", "amount": 5000}'
   ```

2. Check payment via Payment Service:
   ```bash
   curl http://localhost:8081/payments/{order_id}
   ```

3. Observe notification logs in Notification Service container.

## Dead Letter Queue (Bonus)

The system implements DLQ for handling failed message processing:

- **Dead Letter Exchange**: `payment.dlx`
- **Dead Letter Queue**: `payment.dlq`
- **Routing Key**: `failed`

Messages that fail processing multiple times are automatically moved to DLQ. To demonstrate:

1. Modify Notification Service to fail on specific order IDs
2. Send a payment event for that order
3. Observe the message in `payment.dlq` via RabbitMQ Management UI (http://localhost:15672)

The main queue `payment.completed` has `x-dead-letter-exchange` set to `payment.dlx`.
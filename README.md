# Assignment

## **Overview**

This project is a **modular web application** built in **Go**, designed to handle API requests, manage product data, and ensure **transactional consistency** across various services. The project leverages **Redis** for caching, **PostgreSQL** for data storage, and **background task processing** for async jobs (e.g., image processing).

## **Technologies Used**

- **Go**: Backend API development
- **PostgreSQL**: Relational database for data persistence
- **Redis**: In-memory data store for caching
- **RabbitMQ/Redis Queue**: For handling asynchronous tasks
- **Logrus/Zap**: For logging

## **Architecture**

The project is divided into multiple **modules** to ensure maintainability and scalability:

1. **API Module**: Handles incoming HTTP requests, routes, and communicates with other services like the database and cache.
2. **Database Module**: Manages database interactions (CRUD operations, transactions).
3. **Cache Module**: Uses Redis for storing frequently accessed data and reducing database load.
4. **Asynchronous Tasks**: Background jobs are managed via a message queue (RabbitMQ/Redis Queue) to offload heavy tasks such as image processing.
5. **Logging**: Centralized logging is implemented to capture detailed logs for easier debugging and monitoring.

## **Setup Instructions**

### **1. Clone the Repository**
```bash
git clone https://github.com/abxnxsh/assgn.git
cd assgn

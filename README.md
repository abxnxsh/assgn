# Assignment

### Overview

This project is a **scalable, modular** web application designed to handle high API loads, manage transactional consistency, and integrate **asynchronous background tasks** with **caching** and **logging** mechanisms. The goal of the project is to build a robust system for managing product data, processing images, and ensuring data consistency across various services.

Key features:
- **Modular architecture**: The application is broken down into separate modules for API handling, database management, asynchronous tasks, caching, and logging.
- **Scalable**: Designed to handle increased load through **distributed caching** (Redis) and **background task processing** (using message queues).
- **Transactional Consistency**: Ensures data consistency with retries, compensating transactions, and a robust caching strategy.


The system is divided into several key modules, each with a clear responsibility:

- **API Module**: Handles all HTTP requests and routes, with each API route interacting with controllers that manage the business logic.
- **Database Module**: Manages database interactions, including connections, queries, and migrations. A relational database (e.g., PostgreSQL) is used for persistence.
- **Cache Module**: Caching layer built with **Redis** to speed up frequently accessed data.
- **Asynchronous Task Module**: For handling background tasks like image processing, which are offloaded to a message queue (e.g., RabbitMQ or Redis Queue).
- **Logging Module**: Centralized logging system using tools like **Logrus** or **Zap** to maintain structured logs of all activities and errors.

#### **Core Technologies**
- **Go**: Programming language used for the backend services.
- **PostgreSQL**: Relational database used for storing product data.
- **Redis**: Distributed cache for speeding up queries and storing session data.
- **RabbitMQ** or **Redis Queue**: For handling background tasks and queues.
- **Docker**: Containerization of services (API, Redis, PostgreSQL, etc.).

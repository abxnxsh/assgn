# Assignment

### Overview

This project is a **scalable, modular** web application designed to handle high API loads, manage transactional consistency, and integrate **asynchronous background tasks** with **caching** and **logging** mechanisms. The goal of the project is to build a robust system for managing product data, processing images, and ensuring data consistency across various services.

Key features:
- **Modular architecture**: The application is broken down into separate modules for API handling, database management, asynchronous tasks, caching, and logging.
- **Scalable**: Designed to handle increased load through **distributed caching** (Redis) and **background task processing** (using message queues).
- **Transactional Consistency**: Ensures data consistency with retries, compensating transactions, and a robust caching strategy.

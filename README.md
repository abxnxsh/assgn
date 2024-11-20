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
```
## **Tasks and Implementation**

### 1. **API Design**

The design of the API focuses on providing clear and functional endpoints for interacting with product data. Each API endpoint corresponds to specific actions, ensuring ease of use and maintainability:

- **Product Endpoints**:
  - **GET /products/{id}**: Retrieves a product by its ID from the database. If not found in the cache, it fetches from the database and stores the result in Redis.
  - **POST /products**: Inserts a new product into the database. It accepts a JSON body containing product details like `name`, `description`, `images`, and `price`.
 
![image](https://github.com/user-attachments/assets/60ee57a4-5aa1-471a-9c5b-de0f814453d6)


### 2. **Database Integration**

I connected the Go application to a PostgreSQL database by:

- Using the `github.com/lib/pq` driver for PostgreSQL.
- Writing functions to insert and retrieve product data, ensuring smooth interaction with the database.
- Implementing transactional consistency, so if one operation fails, it rolls back the entire transaction, keeping the data consistent.

  ![image](https://github.com/user-attachments/assets/9f95de3e-640f-4aad-8f36-ee5b5f783a71)


### 3. **Caching with Redis**

To reduce database load, Redis was integrated for caching:

- Whenever a product is requested, the system first checks if it exists in Redis.
- If not, it fetches from the database and stores the data in Redis for future requests.
- Redis helps in improving response time by storing frequently accessed data in-memory.

  ![image](https://github.com/user-attachments/assets/d4f54953-b674-45b4-b46b-b83b131eb7bc)


### 4. **Handling Asynchronous Tasks**

For tasks like image processing, I used Redis Queue to handle background jobs:

- When an image needs processing, it's added to a Redis queue.
- A worker retrieves the task from the queue and processes it asynchronously, freeing the main API to handle other requests.

  ![image](https://github.com/user-attachments/assets/129278e5-d6e6-4b71-9ac2-5d3f27254788)

  
### 5. **Logging and Error Handling**

To track activities and errors, I implemented logging:

- **Logrus** or **Zap** was used for centralized logging.
- Logs provide important details like when a product is fetched, inserted, or if an error occurs, making it easier to debug issues.

### 6. **Scalability**

The application is designed with scalability in mind by:

- Caching data in Redis to handle a high volume of read requests.
- Offloading long-running tasks to Redis Queue, allowing the API to respond quickly to user requests.

### 7. **Testing**

To ensure reliability, I wrote unit and integration tests for all major functions:

- **Unit tests** were created for each function, such as database operations, caching, and API endpoints.
- **Integration tests**  for different modules (API, database, and cache) to work together.



# 📊 Globant Data

Welcome to the **Globant Data** repository! This project follows a **microservices architecture** and includes multiple services and scripts written in **Golang** and **Python** for data processing, authorization, and visualization.

---

## 🚀 Features

- 🛠 **Golang Microservices** - Three independent services handling data ingestion, processing, and API functionalities.
- 🐍 **Python Script** - Manages data processing and backup from the core microservices.
- 📊 **Streamlit Dashboard** - Interactive web-based visualization of analytical data.
- 🔗 **Modular Architecture** - Each component operates independently and can be deployed separately, except for microservices requiring a database and the Streamlit dashboard, which depends on the core microservice.

---

## 🏗️ Project Structure

```
📁 globant-data
│── 📂 globant-api         # API Gateway microservice
│   │── 📄 Dockerfile      # Docker configuration for deployment
├── 📂 globant-auth-ms     # Authentication microservice
│   │── 📄 Dockerfile      # Docker configuration for deployment
├── 📂 globant-ms          # Data processing microservice
│   │── 📄 Dockerfile      # Docker configuration for deployment
│── 📂 python-scripts      # Python scripts for data processing
│   │── 📄 Dockerfile      # Docker configuration for deployment
│   │── 📂 task            # Backup and data processing scripts
│── 📄 requirements.txt    # Python dependencies
│── 📂 globant-streamlit   # Streamlit dashboard
│   │── 📄 Dockerfile      # Docker configuration for deployment
│   ├── streamlit_app.py   # Streamlit script
│── 📄 docker-compose.yml  # Docker Compose configuration
│── 📄 README.md           # Project documentation
│── 📄 .gitignore          # Files to ignore in version control
```

---

## 🏗 Microservices Overview

### 🛠 **globant-api (API Gateway)**
- Acts as the central entry point for data ingestion, authentication, and metric retrieval.
- Provides an authorization process to validate user identification.
- Offers API endpoints to trigger data ingestion and retrieve metrics.

### 🛠 **globant-auth-ms (Authentication Service)**
- Manages authentication and authorization.
- Validates tokens and user credentials.
- Encrypts and securely stores user credentials using hashing and salting mechanisms.

### 🛠 **globant-ms (Data Processing Service)**
- Handles ingestion of employee, job, and department data into dynamically generated database schemas.
- Exposes REST API endpoints for accessing processed data.
- Provides employee hiring metrics.

Each microservice follows **clean architecture principles** and can be deployed independently in both on-premise and cloud environments using **Docker**.

---

## 🛠️ Deployment with Docker Compose

To deploy all services using **Docker Compose**, ensure Docker and Docker Compose are installed, then run:

```bash
docker-compose up --build
```

This will start:
- All three **Golang microservices**
- The **Python data processing script** (Scheduled execution)
- The **Streamlit dashboard**
- **PostgreSQL Database**

To stop all services:
```bash
docker-compose down
```

---

## ⚙️ Usage

Use the **Postman collection** as a reference to interact with the services.

### **1️⃣ Create a User**
Create a user using the *CreateUser* endpoint. This step is required since the API gateway mandates user authentication.
```bash
curl --location 'localhost:8083/globant-auth-ms/v1/user/user' \
--header 'Content-Type: application/json' \
--data '{
    "user_name": "test"
}'
```

### **2️⃣ Upload Data Files**
Set up **Authorization** and **X-User** headers to access API endpoints and upload the required files (**Note**: Upload *jobs* and *departments* before *employees* to maintain dependencies.)

- **jobs**
- **departments**
- **employees**

```bash
curl --location 'http://localhost:8080/globant-api/v1/upload' \
--header 'Authorization: 1dz9_Sj1SIwk_FxpwjBIRX2HtyghGHQYWCVYk_gZ2KU=' \
--header 'X-user: 7c28a5de-a135-4f62-b0df-607f4ac651db' \
--form 'file=@"/home/user/folder/hired_employees.csv"'
```

### **3️⃣ Retrieve Insights**
Fetch quarterly hiring metrics using filters:
```bash
curl --location 'localhost:8082/globant-ms/v1/quarter_metrics?year=2021&department_name=Marketing&job_name=Accountant%20I'
```
Retrieve overall hired employee metrics:
```bash
curl --location 'localhost:8082/globant-ms/v1/hired_metrics?year=2021&department_name=Marketing'
```

---

## 📌 Pending Tasks

### **Microservices**
- Expand unit tests to cover all endpoints and error conditions.
- Implement full CRUD operations where applicable.
- Conduct thorough QA testing.

### **Scripts**
- Improve the backup process to allow both cron-based and manual execution.
- Develop a restore backup script.
- Implement a scalable data processing script for handling large files outside the microservices.

### **Dashboard**
- Retrieve data from API endpoints instead of internal microservice endpoints.
- Improve chart visualizations.
- Enhance filtering options.

### **Must-Have Features**
- Implement **Airflow** for pipeline orchestration.
- Develop a **file ingestion traceability** microservice (logging user activity, timestamps, and row counts for ingested data).
- Introduce **Kafka** or **Spark** processing to handle individual row processing and reduce microservice load.
- Enable cloud deployment via **CI/CD** and **Infrastructure as Code (IaC)**.

---

## 📧 Contact
For any inquiries or suggestions, feel free to reach out:
 
🐙 GitHub: [@carloscaceres2093](https://github.com/carloscaceres2093)


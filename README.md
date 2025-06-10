<a name="readme-top"></a>

[![GitHub contributors](https://img.shields.io/github/contributors/adirhmn/payslip-generation)](https://github.com/adirhmn/payslip-generation/graphs/contributors)
[![GitHub forks](https://img.shields.io/github/forks/adirhmn/payslip-generation)](https://github.com/adirhmn/payslip-generation/network)
[![GitHub stars](https://img.shields.io/github/stars/adirhmn/payslip-generation)](https://github.com/adirhmn/payslip-generation/stargazers)

<br />
<div align="center">
  <h3 align="center">Payslip Generation</h3>

  <p align="center">
    Payslip Generation is a REST API designed to automate employee payroll processing. Through this API, administrators can manage attendance periods, track employee attendance, overtime, and reimbursements, and generate accurate payslips based on processed data.
    <br />
    <a href="https://github.com/adirhmn/payslip-generation/issues">Report Bug</a>
    ·
    <a href="https://github.com/adirhmn/payslip-generation/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#software-architecture">Software Architecture</a></li>
        <li><a href="#key-features">Key Features</a></li>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#how-to-running-project">How to Running Project</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <ul>
        <li><a href="#api-endpoints"> API Endpoints</a></li>
        <li><a href="#postman-collection">Postman Collection</a></li>
      </ul>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

Payslip Generation is a backend RESTful service that provides functionality for managing employee payroll operations. It enables administrators to define attendance periods, log daily attendance, record overtime and reimbursements, and process payroll for each period, ensuring a reliable and auditable payslip generation system.

### Software Architecture

The **Payslip Generation** system is designed with a modular and layered architecture to ensure separation of concerns, maintainability, and scalability. The application follows a typical **MVC (Model-View-Controller)** and **Service-Oriented Architecture** with clear responsibilities for each component.

Architecture Overview

```text
┌────────────────────────────────────────────┐
│            Presentation Layer              │
│  (HTTP Handlers / Controllers - REST API)  │
└────────────────────────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│                Service Layer               │
│ (Business logic: payroll, attendance, etc) │
└────────────────────────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│              Repository Layer              │
│      (Database access, query abstraction)  │
└────────────────────────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│                Database Layer              │
│       (PostgreSQL - relational schema)     │
└────────────────────────────────────────────┘
```

This architecture allows smooth integration of features like locking data after payroll processing, linking logs to user requests, and clean separation of business logic.

### Key Features

<b>1. User Management</b>

The system includes 100 pre-filled employees and 1 admin user. Each user has a unique username, password (hashed), full name, and salary. Admins have elevated access rights.

<b>2. Authentication and Authorization</b>

Secure login using JWT-based authentication. Differentiates between Admin and Employee roles. Only authenticated users can access the API, with admin-only privileges for sensitive endpoints (e.g., running payroll, creating attendance periods).

<b>3. Attendance Period Management</b>

Admins can define attendance periods by setting start and end dates. Each period can only be processed once.

<b>4. Attendance Recording</b>

Employees’ presence can be tracked daily. The system ensures uniqueness of attendance entries per user per date.

<b>5. Overtime Submission</b>

Overtime entries can be added with strict validation (1–3 hours per day). Linked to specific attendance periods.

<b>6. Reimbursement Submission</b>

Reimbursements can be logged for each user within a specific attendance period, along with a description and amount.

<b>7. Payroll Processing</b>

Admins can run payroll for a defined period. Once processed, all related attendance, overtime, and reimbursement records become read-only and cannot affect payslip results. Only one payroll can be processed per period.

<b>8. Payslip Generation</b>

Automatically generates payslips based on base salary, attendance, overtime, and reimbursements. Final take-home pay is calculated and stored.

<b>9. Audit Logging</b>

Every change to core records (attendances, payslips, etc.) is logged for traceability with a link to the request_id.

<b>10. Request Logging</b>

Tracks metadata of each request such as URL, IP address, and user performing the action for better traceability and debugging.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

- [![Go][Go]][Go-url]
- [![PostgreSQL][PostgreSQL]][PostgreSQL-url]
- [![Docker][Docker]][Docker-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

Welcome to "Payslip Generation" API! This section will guide you through the process of getting started with payslip generation services.

Get a local copy up and running follow these simple example steps.

### Prerequisites

1. Git Installation:

Make sure Git is installed on your computer. If not, you can download and install it from [here](https://git-scm.com/).

2. Docker Installation:

Install Docker by following the instructions for your operating system from [here](https://docs.docker.com/desktop/install/windows-install/).

### How To Running Project

#### Clone Repository to Local

1. Copy Repository URL

   On the repository page, look for the `Code` or `Clone` button located at the top right. Click on the button and copy the displayed URL (usually in HTTPS or SSH format).

2. Open Terminal

   Open the terminal or command prompt on your computer.

3. Navigate to Destination Directory

   Use the `cd` command to navigate to the directory where you want to store the project. For example:

   ```bash
   cd path/to/destination/directory
   ```

4. Clone Repository

   Type the following command to clone the repository:

   ```bash
   git clone [Repository URL]
   ```

   Replace `[Repository URL]` with the URL

   ```bash
   git clone https://github.com/adirhmn/payslip-generation
   ```

5. Done!

   The project has now been successfully cloned to your computer. You can start working or exploring the code of the project.

#### Running Project with Docker Compose

1. Navigate to Project Directory

   After cloning the repository, navigate to the project directory:

   ```bash
   cd repository-directory
   ```

2. Create an .env file

   You can create it by changing the `.env-example` file to `.env`

3. Install Go Dependencies

   Before running the app, install the necessary Go dependencies:

   ```bash
   go mod tidy
   ```

4. Running DB Migration

   After the containers are running, apply the database migrations using the migrate tool:

   ```bash
   migrate -verbose -path migrations -database "postgres://user:password@127.0.0.1:6932/payslipdb?sslmode=disable" up
   ```

   Ensure the database URL, user, and password match those in your .env file.

5. Start Docker Compose

   Make sure Docker is installed and running on your system.
   Run the following command to start the application using Docker Compose:

   ```bash
   docker-compose up --build
   ```

   This command will build the Docker images and start the containers in detached mode.
   Congratulations your app has been running at `http://localhost:8080` and you can checking status app by access `http://localhost:8080/v1/ping` and you will get response

   ```
   {
    "success": true,
    "error": "",
    "data": {
        "server_says": "pong"
      }
   }
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

#### Running Unit Testing

To run unit tests and view the coverage report, follow these steps:

1. Navigate to the Project Root

   Make sure you're in the root directory of the project:

   ```bash
   cd repository-directory
   ```

2. Run Unit Tests with Coverage

   Use the following command to run the tests and generate a coverage report:

   ```bash
   go test ./... -coverprofile=cover.out
   ```

3. View the Coverage Report

   After running the tests, use the Go cover tool to open the HTML coverage report:

   ```bash
   go tool cover -html=cover.out
   ```

   This will open a browser window displaying the test coverage across your codebase.

<!-- USAGE EXAMPLES -->

## Usage

#### API Endpoints

Explore and interact with the Payslip Generation API using the following endpoints. You can use tools like [Postman](https://www.postman.com/) for a convenient API testing experience.

#### Postman Collection

To simplify API testing, you can use the provided Postman collection.

1. **Download Postman:**
   If you don't have Postman installed, you can download it [here](https://www.postman.com/downloads/).

2. **Import Collection:**
   Download Payslip Generation API Postman collection [here](https://github.com/adirhmn/payslip-generation/blob/main/payslip_generation.postman_collection), and import it into Postman.

3. **Set Environment Variables:**

   - Create a new environment in Postman.
   - Set the following variables:
     - `localhost`: `http://localhost:8080/api/v1/`
     - `token_employee`: `[BEARER_TOKEN_EMPLOYEE]`
     - `token_admin`: `[BEARER_TOKEN_ADMIN]`

4. **Get Bearer Token for Authentication**  
   To obtain the `BEARER_TOKEN_EMPLOYEE` or `BEARER_TOKEN_ADMIN`, you must first log in using one of the default accounts:

   **Admin Account:**

   - Username: `admin`
   - Password: `admin123`

   **Employee Accounts:**

   - Username: `employee_1`, `employee_2`, ..., `employee_100`
   - Password: `employee123`

   **Login Endpoint:**  
   `POST http://localhost:8080/v1/login`

   **Request Body:**

   ```json
   {
     "username": "admin",
     "password": "admin123"
   }
   ```

   **Example Response:**

   ```json
   {
     "success": true,
     "error": "",
     "data": {
       "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
     }
   }
   ```

   Copy the token from the response and use it to replace `BEARER_TOKEN_ADMIN` or `BEARER_TOKEN_EMPLOYEE` in your Postman environment. Token only available for 24 hours.

5. **Explore and Test:**

   - Browse the available requests in the Pyaslip Generation API collection.
   - Update variables like token regularly when the validity period expires.

6. **Execute Requests:**
   - Execute requests to add attendance period, submit attendance, and more.

_For screenshot, please look up to the [Documentation](https://github.com/adirhmn/payslip-generation/tree/main/documentation)_

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[PostgreSQL]: https://img.shields.io/badge/PostgreSQL-20232A?style=for-the-badge&logo=postgresql&logoColor=61DAFB
[PostgreSQL-url]: https://www.postgresql.org/
[Go]: https://img.shields.io/badge/Go-4A4A55?style=for-the-badge&logo=go&logoColor=61DAFB
[Go-url]: https://go.dev/
[Docker]: https://img.shields.io/badge/Docker-0769AD?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/

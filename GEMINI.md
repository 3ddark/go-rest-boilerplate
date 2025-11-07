# Backend Project Context

This document provides an overview of the backend project for the ERP system.

## 1. Project Goal

The main goal of this project is to build a multi-language supported web ERP project. The backend is responsible for handling business logic, data processing, and providing APIs for the frontend.

## 2. Architecture

The project follows a clean architecture pattern, separating concerns into different layers:

- **`cmd`**: Contains the main entry points for the application (API server and worker).
- **`internal`**: Contains the core business logic of the application.
  - **`domain`**: Defines the core data models and interfaces.
  - **`repository`**: Implements data access logic, interacting with the database.
  - **`service`**: Contains the business logic, orchestrating repositories and other services.
  - **`handler`**: Handles incoming requests (HTTP and GraphQL) and calls the appropriate services.
  - **`platform`**: Provides platform-level functionalities like database connection, authentication, logging, and message queuing.
  - **`dto`**: Data Transfer Objects used for transferring data between layers.

## 3. Core Technologies

- **Go**: The primary programming language.
- **Fiber**: A Go web framework for building high-performance APIs.
- **GORM**: An ORM library for Go to interact with the database.
- **PostgreSQL**: The primary database.
- **RabbitMQ**: A message broker for asynchronous tasks.
- **GraphQL**: A query language for APIs.
- **JWT**: For securing the APIs.
- **Zerolog**: For structured logging.
- **Prometheus**: For collecting metrics.

## 4. Key Features

### 4.1. API Server (`cmd/api`)

- Provides both RESTful and GraphQL APIs.
- Handles user authentication and authorization.
- Exposes endpoints for managing users and reports.
- Includes middleware for logging, recovery, CORS, and Prometheus metrics.

### 4.2. Worker (`cmd/worker`)

- A separate process for handling asynchronous tasks.
- Consumes messages from RabbitMQ.
- Currently, it processes report generation and sends welcome emails.

### 4.3. User Management

- **Authentication**: Users can log in via the `/api/v1/login` endpoint or the `login` GraphQL mutation.
- **Authorization**: Role-based access control is implemented using a permission service.
- **CRUD Operations**: The API provides endpoints for creating, reading, updating, and deleting users.

### 4.4. Asynchronous Report Generation

- Users can request a report via the `/api/v1/reports` endpoint.
- The request is published to a RabbitMQ queue.
- A worker process consumes the message, generates the report, and updates its status in the database.
- The status of the report can be checked via the `/api/v1/reports/:id` endpoint.

## 5. Database Schema

The database schema is defined by the GORM models in the `internal/domain` directory.

- **`users`**: Stores user information.
- **`user_permissions`**: Stores user permissions for different resources.
- **`reports`**: Stores information about generated reports.

## 6. Configuration

The application is configured using environment variables, loaded by the `godotenv` library. The configuration is defined in the `internal/config/config.go` file.

## 7. Internationalization (i18n)

The project has a placeholder for i18n, but it is not yet fully implemented.

## 8. Next Steps

- Implement the i18n functionality.
- Add more features to the ERP system.
- Write more unit and integration tests.
- Enhance the monitoring and alerting capabilities.

# UI Project Context

This document provides an overview of the UI project for the ERP system.

## 1. Project Goal

The main goal of this project is to build a modern, responsive, and feature-rich user interface for the ERP system. It will interact with the backend APIs to display and manage data.

## 2. Architecture

The project is a Next.js application with TypeScript. It follows a component-based architecture, with a clear separation of concerns:

- **`src/app`**: Contains the main pages and layouts of the application.
- **`src/components`**: Contains reusable components, such as the `DataTable`.
- **`src/api`**: Contains functions for interacting with the backend APIs.
- **`src/utils`**: Contains utility functions for tasks like data exporting and local storage management.

## 3. Core Technologies

- **Next.js**: A React framework for building server-side rendered and static web applications.
- **React**: A JavaScript library for building user interfaces.
- **TypeScript**: A typed superset of JavaScript that compiles to plain JavaScript.
- **Tailwind CSS**: A utility-first CSS framework for rapid UI development.
- **@tanstack/react-table**: A headless UI library for building powerful data tables.
- **ESLint**: A tool for identifying and reporting on patterns found in ECMAScript/JavaScript code.

## 4. Key Features

### 4.1. Data Table

The core of the application is a powerful and reusable `DataTable` component with a wide range of features:

- **Server-side Processing**: The table is designed to work with server-side data, handling pagination, sorting, and filtering efficiently.
- **Customizable Columns**: The columns can be easily configured, with options for sticky columns and custom cell rendering.
- **Row Selection**: Users can select single or multiple rows.
- **Filtering and Sorting**: The table supports global and column-specific filtering, as well as multi-column sorting.
- **Column Resizing and Visibility**: Users can resize and toggle the visibility of columns.
- **Data Export**: The table data can be exported to CSV and Excel formats.
- **User Preferences**: The user's preferences for the table (column visibility, filters, sorting, etc.) are saved to local storage.

### 4.2. API Interaction

The `src/api/userService.ts` file currently uses mock data to simulate API calls. This will be replaced with actual API calls to the backend to fetch and manage user data.

### 4.3. Styling

The application is styled using Tailwind CSS, which allows for rapid and consistent styling. The main styles are defined in `src/app/globals.css`.

## 5. Next Steps

- Replace the mock API service with actual API calls to the backend.
- Implement authentication and authorization, so that only logged-in users can access the application.
- Add more pages and features to the ERP system, such as forms for creating and editing data.
- Implement internationalization (i18n) to support multiple languages.
- Write unit and integration tests for the components and pages.

# 🎶 Music Service (Go)

## 📖 Description

**Music Service** is a Go-based project that integrates with the **Spotify API** to retrieve music track data. It includes a structured architecture with repository, service, and handler layers. The project uses **Docker** and **Docker Compose** for containerized deployment and supports automation with **Taskfile**.

## 🚀 Features

- 🎧 **Spotify API Integration**: Retrieve tracks.
- 🗄 **MySQL Database**: Stores playlists, songs ,users and relation between songs and playlists.
- 🏗 **Layered Architecture**: Separation of logic between repositories, services, and handlers.
- 🧪 **Comprehensive Testing**: Unit tests for handlers and repositories.
- 🐳 **Docker & Docker Compose**: For container management.
- 🔄 **Database schema**: `migrate-create` functionality to creating schemas
- 🔄 **Database Migrations**: `migration-up/down` functionality to manage schema changes.
- 🛠 **Taskfile Integration**: Automate tasks like migrations, container startup, and testing.

## 🛠 Project Setup
### Prerequisites

Ensure you have the following installed:

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Taskfile](https://taskfile.dev/installation/)

### 🏗 Local Development Setup

1. **Clone the repository**:
    ```bash
    git clone https://github.com/cloud9cloud9/music-service-go.git
    ```

2. **Navigate to the project directory**:
    ```bash
    cd music-service
    ```

3. **Run Docker containers**:
    ```bash
    task db-up
    ```

4. **Create new migration files**:
   Before running the migrations, create new migration files for `up` and `down` by using the following command:
    ```bash
    task migrate-create
    ```

5. **Run database migrations**:
   After creating the migration files, run the migrations:
    ```bash
    task migration-up
    ```

6. **Run tests**:
   Finally, you can run the tests:
    ```bash
    task test
    ```


# E-Commerce API Core

A robust, logic-heavy backend service for an e-commerce platform built with Go and PostgreSQL. This API handles everything from secure user authentication and role-based access control to persistent shopping carts and a custom mock payment gateway integration.

Built based on the architectural requirements of backend learning roadmaps to demonstrate complex data modeling and internal service communication.

---

## Tech Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL (with `lib/pq` driver)
* **Authentication:** JSON Web Tokens (JWT) & bcrypt password hashing
* **Payment Gateway:** Custom Mock Stripe Service (Internal HTTP communication)
* **Infrastructure:** Docker-ready database connections

---

## Key Features

* **Secure Authentication:** User registration and login utilizing Bcrypt for password hashing and JWT for stateless session management.
* **Role-Based Access Control (RBAC):** Middleware protects administrative routes. Only users with the `admin` role can add, update, or delete products.
* **Smart Shopping Cart:** * Cart items are tied to the authenticated user.
  * Uses PostgreSQL `ON CONFLICT` constraints to intelligently increment product quantities if an item already exists in the cart.
* **Mock Payment Gateway (Stripe Simulation):** Calculates the exact cart total from the secure database (preventing client-side price manipulation) and communicates via HTTP with an internal mock service to generate a secure checkout session URL.
* **Automatic Migrations:** The application automatically verifies and initializes all required database tables (`users`, `products`, `cart_items`) on startup.

---

## 📖 API Reference

### 🔐 Authentication
* `POST /register` - Create a new user account.
* `POST /login` - Authenticate and receive a JWT token.

### 📦 Products Management
* `GET /products` - Retrieve a list of all products (Public).
* `GET /products/{id}` - Retrieve a specific product by ID (Public).
* `POST /products` - Add a new product (**Admin only**).
* `PUT /products/{id}` - Update product details (**Admin only**).
* `DELETE /products/{id}` - Delete a product (**Admin only**).

### 🛒 Shopping Cart (Requires JWT)
* `GET /cart` - View the current user's shopping cart and total items.
* `POST /cart` - Add a product to the cart or increase its quantity.

### 💳 Checkout (Requires JWT)
* `POST /checkout` - Calculates the total basket value and requests a payment session from the internal mock gateway.

---

## 🏗️ Getting Started (Local Development)

### Prerequisites
* Go 1.22+ installed.
* A running PostgreSQL instance (can be spun up via Docker).

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/MrGulo/E-CommerceAPI.git
   cd E-CommerceAPI
   ```

2. **Configure Database:**
   Ensure PostgreSQL is running on port `5433` (or update the connection string in `database/db.go`) with the user `postgres` and password `secret`.

3. **Install Dependencies:**
   ```bash
   go mod download
   ```

4. **Run the API:**
   ```bash
   go run main.go
   ```
   *The server will start on `http://localhost:8080` and automatically create the required database schema.*

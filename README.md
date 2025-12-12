# Kora-Excahnge - Borderless Payment Platform
---

## Features

###  Multi-Currency Wallet System
- Support for 6 stablecoins pegged to African and international currencies:
  - **cNGN** - Nigerian Naira Stablecoin →  `NGN`
  - **cXAF** - CFA Franc Stablecoin → `XAF`
  - **USDx** - USD Stablecoin → `USD`
  - **EURx** - EUR Stablecoin → `EUR`
  - **cGHS** - Ghanaian Cedi Stablecoin → `GHS`
  - **cKES** - Kenyan Shilling Stablecoin → `KES`

### Core Transaction Features
1. **Deposit** - Instantly credit your wallet with any supported currency
2. **Swap** - Exchange between different currencies with real-time FX rates
3. **Transfer** - Send funds to other users in same or different currencies with automatic conversion
4. **Real-time Exchange Rates** - Powered by FastForex API with 24-hour caching

### Dashboard Features
- **Portfolio Overview** - Total balance in USD with live conversion
- **Wallet Distribution Chart** - Visual pie chart showing currency allocation
- **Exchange Rates Card** - Live rates for all supported currencies with refresh option
- **Transaction History** - Complete audit trail of all operations
- **Wallet Address** - Unique address for receiving funds with one-click copy


---

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Chi Router v5
- **Database**: PostgreSQL with pgxpool
- **Authentication**: JWT (golang-jwt)
- **API Integration**: FastForex API for real-time exchange rates

### Frontend
- **Framework**: React 19
- **Styling**: TailwindCSS
- **State Management**: TanStack Query (React Query)
- **Routing**: React Router DOM v7
- **Icons**: Lucide React

### DevOps & Tools
- **Version Control**: Git
- **Package Managers**: Go modules, npm
- **Database Migration**: SQL scripts

---


## Setup Instructions

### 1. Clone the Repository
```bash
git clone <repository-url>
cd Interstellar
```

### 2. Database Setup

**Create PostgreSQL Database:**
```bash
psql -U postgres
CREATE DATABASE interstellar;
\q
```

**Run Database Migrations:**
```bash
cd backend
psql -U postgres -d interstellar -f sql/db.sql
```

### 3. Backend Setup

**Navigate to backend directory:**
```bash
cd backend
```

**Create `.env` file:**
```env
PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/interstellar?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
EXCHANGERATE_API_KEY=your-fastforex-api-key
```

**Install dependencies:**
```bash
go mod download
```

**Run the server:**
```bash
go run cmd/server/main.go
```

The backend server will start on `http://localhost:8080`

### 4. Frontend Setup

**Navigate to frontend directory:**
```bash
cd ../frontend
```

**Install dependencies:**
```bash
npm install
```

**Create `.env` file (optional):**
```env
REACT_APP_API_URL=http://localhost:8080
```

**Start the development server:**
```bash
npm start
```

The frontend will start on `http://localhost:3000`

---

## API Keys

### FastForex API Key
1. Sign up at [FastForex.io](https://www.fastforex.io/)
2. Get your free API key
3. Add it to your backend `.env` file as `EXCHANGERATE_API_KEY`

---

## Demo Accounts

Since this is a sandbox prototype, you can create your own account through the registration page. There are no pre-configured admin credentials.

**To create an account:**
1. Navigate to `http://localhost:3000/register`
2. Fill in your details:
   - Name
   - Email
   - Password (min 6 characters)
3. Click "Create Account"
4. You'll be automatically logged in and redirected to the dashboard

**To test transfers between accounts:**
1. Create a second account in a different browser/incognito window
2. Copy the wallet address from the dashboard banner
3. Use the "Send" feature from your first account to transfer funds

---

## 🎯 Feature Walkthrough

### 1. Registration & Login
- **Register**: Create a new account with email and password
- **Login**: Access your wallet with credentials
- **Auto Wallet Creation**: A unique wallet with address is created upon registration

### 2. Deposit Funds
1. Click **"Deposit"** button on dashboard
2. Select currency (e.g., cNGN, USDx, EURx)
3. Enter amount
4. Confirm - funds are instantly credited

### 3. Swap Currencies
1. Click **"Swap"** button
2. Select "From" currency and amount
3. Select "To" currency
4. View real-time conversion rate and estimated amount
5. Confirm swap - balances update instantly

### 4. Send Money (Transfer)
1. Click **"Send"** button
2. Enter recipient's wallet address
3. Select currency and amount
4. **Optional**: Enable "Convert to a different currency"
   - Select recipient's preferred currency
   - System automatically converts using real-time FX rates
5. Confirm - funds are transferred instantly

### 5. View Exchange Rates
- Live rates displayed on dashboard
- Click refresh icon to force update
- Shows rates for NGN, XAF, EUR, GHS, KES against USD
- Last updated timestamp displayed

### 6. Portfolio Dashboard
- **Total Portfolio Value**: All balances converted to USD
- **Wallet Distribution**: Pie chart showing currency allocation
- **Individual Balances**: Card view for each currency
- **Wallet Address**: Easily copy your address to share with others

---

## Project Structure

```
Interstellar/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       ├── main.go          # Application entry point
│   │       └── api.go           # Route definitions
│   ├── internal/
│   │   ├── fxrates/            # FX rate service (FastForex integration)
│   │   ├── middleware/         # JWT auth & CORS
│   │   ├── transactions/       # Deposit, swap, transfer logic
│   │   ├── users/              # User management
│   │   ├── wallets/            # Wallet operations
│   │   └── utils/              # JWT utilities
│   ├── pkg/
│   │   └── response/           # Standardized API responses
│   ├── sql/
│   │   └── db.sql              # Database schema
│   ├── go.mod
│   ├── go.sum
│   └── .env
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/
│   │   │   ├── DepositModal.js
│   │   │   ├── SwapModal.js
│   │   │   ├── SendModal.js
│   │   │   ├── WalletCard.js
│   │   │   ├── WalletPieChart.js
│   │   │   ├── ExchangeRatesCard.js
│   │   │   ├── Layout.js
│   │   │   └── Navbar.js
│   │   ├── hooks/
│   │   │   └── useWallet.js    # React Query hooks
│   │   ├── pages/
│   │   │   ├── Login.js
│   │   │   ├── Register.js
│   │   │   └── Dashboard.js
│   │   ├── services/
│   │   │   └── api.js          # API client
│   │   ├── App.js
│   │   └── index.js
│   ├── package.json
│   └── .env
│
└── README.md
```

---

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Wallets (Protected)
- `GET /api/wallets` - Get user's wallet
- `GET /api/wallets/balances` - Get all balances
- `GET /api/wallets/balance/{currency}` - Get specific currency balance

### Transactions (Protected)
- `POST /api/transactions/deposit` - Deposit funds
- `POST /api/transactions/swap` - Swap currencies
- `POST /api/transactions/transfer` - Transfer to another wallet
- `GET /api/transactions` - Get transaction history
- `GET /api/transactions/{id}` - Get specific transaction

### FX Rates (Public)
- `GET /api/fx-rates?base=USD` - Get all exchange rates
- `POST /api/fx-rates/convert` - Convert between currencies
- `POST /api/fx-rates/refresh` - Force refresh cache

---

### Database Schema
- **users**: User accounts with credentials
- **wallets**: One wallet per user with JSONB balances
- **transactions**: Comprehensive transaction log with support for all transaction types

### Caching Strategy
- FX rates are cached for 24 hours to reduce API calls
- Cache invalidation via manual refresh endpoint



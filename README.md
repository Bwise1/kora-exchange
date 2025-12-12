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
- **Security Audit Logs** - Password-protected access to view all account activities including IP addresses, request methods, and timestamps


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
git clone https://github.com/Bwise1/kora-exchange.git
cd kora-exchange
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

**Setup env :**
```bash
cp .env.example .env
```
- replace env variables
  
**Install dependencies:**
```bash
go mod download
```

**Run the server:**
```bash
go run cmd/server
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

**Audit Logs Access:**
- Default audit password: `admin123`
- Can be changed in backend `.env` file via `AUDIT_PASSWORD` variable

---

## Feature Walkthrough

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

### 7. Security Audit Logs
1. Click **"Audit Logs"** in the side menu
2. Enter the audit access password (default: `admin123`)
3. View comprehensive security logs including:
   - Operation type (LOGIN, REGISTER, DEPOSIT, SWAP, TRANSFER, etc.)
   - Timestamp of each action
   - Client IP address (with Cloudflare support)
   - HTTP request method and path
4. Track all account activities for security and compliance

---

## Project Structure

```
Interstellar/
├── backend/
├── frontend/
└── README.md
```

---

### Deployment
- backend is deployed on a vps and auto deployment is handled with github actions https://kora.benjys.me/
- frontend is deployed on vercel https://kora-exchange.vercel.app/dashboard
---

### Database Schema
- **users**: User accounts with credentials
- **wallets**: One wallet per user with JSONB balances
- **transactions**: Comprehensive transaction log with support for all transaction types
- **audit_logs**: Security audit trail with user_id, IP addresses, operations, and request metadata

### Caching Strategy
- FX rates are cached for 24 hours to reduce API calls
- Cache invalidation via manual refresh endpoint



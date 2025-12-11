# Operation Borderless - Implementation Action Plan

## Project Overview
Build a multi-currency stablecoin payment sandbox supporting African and international stablecoins (cNGN, cXAF, USDx, EURx) with real-time FX swaps, cross-border transfers, and compliance logging.

---

## Database Schema Design

### 1. Users Table (Already Exists ✓)
```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
```

### 2. Wallets Table
```sql
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_address VARCHAR(100) UNIQUE NOT NULL, -- Unique wallet identifier
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_address ON wallets(wallet_address);
```

### 3. Wallet Balances Table
Multi-currency balances for each wallet (one row per currency per wallet)
```sql
CREATE TABLE IF NOT EXISTS wallet_balances (
    id UUID PRIMARY KEY,
    wallet_id UUID NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    currency_code VARCHAR(10) NOT NULL, -- cNGN, cXAF, USDx, EURx, etc.
    balance NUMERIC(20, 8) NOT NULL DEFAULT 0, -- High precision for crypto
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_wallet FOREIGN KEY (wallet_id) REFERENCES wallets(id),
    CONSTRAINT unique_wallet_currency UNIQUE (wallet_id, currency_code),
    CONSTRAINT positive_balance CHECK (balance >= 0)
);

CREATE INDEX idx_wallet_balances_wallet_id ON wallet_balances(wallet_id);
CREATE INDEX idx_wallet_balances_currency ON wallet_balances(currency_code);
```

### 4. Currencies Table (Reference Data)
```sql
CREATE TABLE IF NOT EXISTS currencies (
    code VARCHAR(10) PRIMARY KEY, -- cNGN, cXAF, USDx, EURx
    name VARCHAR(100) NOT NULL, -- "Nigerian Naira Stablecoin"
    symbol VARCHAR(10), -- ₦, $, €
    decimals INT DEFAULT 8, -- Precision for display
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Insert supported currencies
INSERT INTO currencies (code, name, symbol, decimals) VALUES
('cNGN', 'Nigerian Naira Stablecoin', '₦', 2),
('cXAF', 'Central African Franc Stablecoin', 'FCFA', 2),
('USDx', 'USD Stablecoin', '$', 2),
('EURx', 'EUR Stablecoin', '€', 2),
('cGHS', 'Ghanaian Cedi Stablecoin', '₵', 2),
('cKES', 'Kenyan Shilling Stablecoin', 'KSh', 2)
ON CONFLICT (code) DO NOTHING;
```

### 5. Transactions Table
Records all transaction types: DEPOSIT, SWAP, TRANSFER
```sql
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    transaction_type VARCHAR(20) NOT NULL, -- DEPOSIT, SWAP, TRANSFER
    status VARCHAR(20) NOT NULL DEFAULT 'COMPLETED', -- PENDING, COMPLETED, FAILED

    -- Source wallet info
    from_wallet_id UUID REFERENCES wallets(id),
    from_currency VARCHAR(10) NOT NULL,
    from_amount NUMERIC(20, 8) NOT NULL,

    -- Destination wallet info
    to_wallet_id UUID REFERENCES wallets(id),
    to_currency VARCHAR(10) NOT NULL,
    to_amount NUMERIC(20, 8) NOT NULL,

    -- FX rate information (for swaps and cross-currency transfers)
    exchange_rate NUMERIC(20, 8), -- Rate used for conversion
    exchange_rate_id UUID REFERENCES fx_rates(id), -- Reference to FX rate used

    -- Metadata
    description TEXT,
    transaction_hash VARCHAR(100) UNIQUE, -- Unique transaction identifier

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT check_transaction_type CHECK (transaction_type IN ('DEPOSIT', 'SWAP', 'TRANSFER')),
    CONSTRAINT check_status CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED'))
);

CREATE INDEX idx_transactions_from_wallet ON transactions(from_wallet_id);
CREATE INDEX idx_transactions_to_wallet ON transactions(to_wallet_id);
CREATE INDEX idx_transactions_type ON transactions(transaction_type);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_hash ON transactions(transaction_hash);
```

### 6. FX Rates Table
Store real-time and historical exchange rates
```sql
CREATE TABLE IF NOT EXISTS fx_rates (
    id UUID PRIMARY KEY,
    from_currency VARCHAR(10) NOT NULL,
    to_currency VARCHAR(10) NOT NULL,
    rate NUMERIC(20, 8) NOT NULL, -- Exchange rate
    source VARCHAR(50), -- 'API', 'MANUAL', 'MOCK'
    valid_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    valid_until TIMESTAMP, -- NULL means currently valid
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    CONSTRAINT unique_currency_pair_time UNIQUE (from_currency, to_currency, valid_from)
);

CREATE INDEX idx_fx_rates_pair ON fx_rates(from_currency, to_currency);
CREATE INDEX idx_fx_rates_active ON fx_rates(is_active, valid_from DESC);
CREATE INDEX idx_fx_rates_valid_range ON fx_rates(from_currency, to_currency, valid_from, valid_until);
```

### 7. Audit Logs Table (Compliance Mode)
Track all user actions with IP, device, country, browser
```sql
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL, -- 'LOGIN', 'DEPOSIT', 'SWAP', 'TRANSFER', etc.
    resource_type VARCHAR(50), -- 'WALLET', 'TRANSACTION', 'USER'
    resource_id UUID, -- ID of the resource being acted upon

    -- Request metadata
    ip_address INET NOT NULL,
    user_agent TEXT,
    device_type VARCHAR(50), -- 'mobile', 'desktop', 'tablet'
    browser VARCHAR(50),
    country_code VARCHAR(3), -- ISO country code

    -- Additional context
    request_method VARCHAR(10), -- GET, POST, PUT, DELETE
    request_path TEXT,
    status_code INT,

    -- Data
    old_value JSONB, -- Previous state (for updates)
    new_value JSONB, -- New state

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_audit_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_action ON audit_logs(action);
CREATE INDEX idx_audit_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_ip ON audit_logs(ip_address);
CREATE INDEX idx_audit_resource ON audit_logs(resource_type, resource_id);
```

### 8. AI Assistant Queries Table (Bonus Feature)
Log AI assistant interactions
```sql
CREATE TABLE IF NOT EXISTS ai_queries (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    query_text TEXT NOT NULL,
    response_text TEXT,
    query_type VARCHAR(50), -- 'FX_CONVERSION', 'RATE_INQUIRY', 'TREND_ANALYSIS'
    currency_pair VARCHAR(20), -- e.g., 'cNGN/USDx'
    processing_time_ms INT,
    model_used VARCHAR(50), -- 'gpt-4', 'gemini-pro', etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_ai_queries_user ON ai_queries(user_id);
CREATE INDEX idx_ai_queries_created_at ON ai_queries(created_at DESC);
```

---

## Implementation Phases

### Phase 1: Foundation (Database & Core Wallet System)
**Goal:** Set up database, wallet creation, and balance management

#### Tasks:
1. **Create migration files**
   - Split db.sql into versioned migrations
   - Create `sql/migrations/001_create_tables.sql`

2. **Implement Wallets Module** (`internal/wallets/`)
   - `models.go` - Wallet, WalletBalance, Currency structs
   - `repository.go` - CRUD operations for wallets and balances
   - `service.go` - Business logic (create wallet, get balances, calculate USD equivalent)
   - `handler.go` - HTTP handlers for wallet endpoints

3. **API Endpoints:**
   - `POST /api/wallets` - Create wallet (auto-generate wallet address)
   - `GET /api/wallets` - Get user's wallets
   - `GET /api/wallets/:id/balances` - Get all currency balances for a wallet
   - `GET /api/wallets/:id/balance/usd` - Get total USD equivalent

4. **Helper Package** (`pkg/utils/`)
   - Wallet address generator (unique identifiers)
   - Currency converter (rate * amount calculations)

---

### Phase 2: Deposits & Transactions
**Goal:** Enable simulated deposits and transaction recording

#### Tasks:
1. **Implement Transactions Module** (`internal/transactions/`)
   - `models.go` - Transaction, DepositRequest structs
   - `repository.go` - Save, query transactions
   - `service.go` - Process deposits, validate balances
   - `handler.go` - HTTP handlers

2. **API Endpoints:**
   - `POST /api/wallets/:id/deposit` - Simulate deposit
     - Body: `{"currency": "cNGN", "amount": 1000}`
   - `GET /api/wallets/:id/transactions` - Get transaction history
   - `GET /api/transactions/:id` - Get single transaction details

3. **Business Logic:**
   - Validate currency is supported
   - Update wallet balance atomically
   - Create transaction record with unique hash
   - Return updated balance

---

### Phase 3: FX Rates & Currency Swaps
**Goal:** Integrate FX rates and enable currency swaps

#### Tasks:
1. **Implement FX Rates Module** (`internal/fxrates/`)
   - `models.go` - FXRate, RateRequest structs
   - `repository.go` - Store and retrieve rates
   - `service.go` - Fetch live rates, cache rates, calculate conversions
   - `handler.go` - HTTP handlers for rate queries
   - `fetcher.go` - Integration with FX rate API (e.g., Exchange Rate API, CurrencyAPI)

2. **Rate Provider Options:**
   - **Live Mode:** Use free API like exchangerate-api.com or currencyapi.com
   - **Mock Mode:** Use fixed rates stored in database (via env var `USE_MOCK_RATES=true`)

3. **API Endpoints:**
   - `GET /api/fx-rates` - Get all current rates
   - `GET /api/fx-rates/:from/:to` - Get specific pair rate
   - `POST /api/fx-rates/refresh` - Manually refresh rates (admin only)

4. **Implement Swap Functionality** (`internal/swaps/`)
   - `service.go` - Swap logic with validations
   - `handler.go` - Swap endpoint

5. **Swap API Endpoint:**
   - `POST /api/wallets/:id/swap`
     - Body: `{"from_currency": "cNGN", "to_currency": "USDx", "amount": 5000}`

6. **Swap Business Logic:**
   - Get current FX rate
   - Validate sufficient balance in from_currency
   - Calculate to_amount = amount * rate
   - Deduct from_currency balance
   - Add to_currency balance
   - Create SWAP transaction record
   - Return transaction details

---

### Phase 4: Cross-Border Transfers
**Goal:** Enable sending funds between wallets with auto-conversion

#### Tasks:
1. **Implement Transfers Module** (`internal/transfers/`)
   - `service.go` - Transfer logic with optional currency conversion
   - `handler.go` - Transfer endpoint

2. **API Endpoint:**
   - `POST /api/wallets/:id/transfer`
     - Body:
     ```json
     {
       "to_wallet_address": "WLT-ABC123",
       "currency": "cNGN",
       "amount": 10000,
       "target_currency": "USDx" // Optional: if different, auto-swap
     }
     ```

3. **Transfer Business Logic:**
   - Validate sender has sufficient balance
   - Validate recipient wallet exists
   - If same currency: simple transfer
   - If different currency:
     - Get FX rate
     - Calculate converted amount
     - Deduct from sender in from_currency
     - Add to recipient in to_currency
   - Create TRANSFER transaction record
   - Return transaction details

---

### Phase 5: Middleware & Security
**Goal:** Add CORS, audit logging, and compliance features

#### Tasks:
1. **CORS Middleware** (`internal/middleware/cors.go`)
   - Support both production and local frontend URLs
   - Configure allowed methods, headers

2. **Audit Logging Middleware** (`internal/middleware/audit.go`)
   - Extract IP address (handle proxy headers)
   - Parse User-Agent for browser/device info
   - Extract country from IP (use GeoIP library or service)
   - Log every request to audit_logs table
   - Enrich with user_id from JWT context

3. **Rate Limiting Middleware** (Optional but recommended)
   - Prevent abuse of swap/transfer endpoints

4. **Update api.go:**
   - Add CORS middleware globally
   - Add audit middleware to protected routes
   - Add rate limiting to sensitive endpoints

---

### Phase 6: Transaction Explorer & Analytics
**Goal:** Build explorer API and wallet analytics

#### Tasks:
1. **Implement Explorer Module** (`internal/explorer/`)
   - `service.go` - Query recent transactions, top currencies, volume stats
   - `handler.go` - Explorer endpoints

2. **API Endpoints:**
   - `GET /api/explorer/transactions` - Recent transactions (paginated)
     - Query params: `?page=1&limit=20&type=TRANSFER`
   - `GET /api/explorer/stats` - Platform statistics
     - Total volume, transaction count, active wallets
   - `GET /api/explorer/transactions/:hash` - Transaction by hash

3. **Wallet Analytics:**
   - `GET /api/wallets/:id/analytics` - Wallet-specific stats
     - Balance distribution (for pie chart)
     - Transaction volume by currency
     - Recent activity summary

---

### Phase 7: AI Assistant Integration (Bonus)
**Goal:** Add LLM-powered FX assistant

#### Tasks:
1. **Implement AI Module** (`internal/ai/`)
   - `client.go` - OpenAI/Gemini client wrapper
   - `service.go` - Process queries, extract intent, format responses
   - `handler.go` - Chat endpoint

2. **API Endpoint:**
   - `POST /api/ai/assistant`
     - Body: `{"query": "Convert 500 cNGN to USDx"}`
     - Response:
     ```json
     {
       "response": "500 cNGN = 0.32 USDx at current rate (1550 NGN/USD)",
       "data": {
         "from_currency": "cNGN",
         "to_currency": "USDx",
         "amount": 500,
         "result": 0.32,
         "rate": 1550
       }
     }
     ```

3. **AI Capabilities:**
   - Currency conversion calculations
   - Rate comparisons
   - Historical trend queries (requires storing FX rate history)
   - Best time to swap recommendations

4. **Implementation Options:**
   - **OpenAI GPT-4:** Best quality, requires API key
   - **Google Gemini:** Good alternative, free tier available
   - **Local LLM:** Use Ollama with llama2/mistral for privacy

---

## API Endpoints Summary

### Authentication
- `POST /api/auth/register` - User registration ✓
- `POST /api/auth/login` - User login ✓

### Wallets
- `POST /api/wallets` - Create wallet
- `GET /api/wallets` - Get user's wallets
- `GET /api/wallets/:id/balances` - Get wallet balances
- `GET /api/wallets/:id/balance/usd` - Get USD equivalent

### Transactions
- `POST /api/wallets/:id/deposit` - Deposit funds
- `POST /api/wallets/:id/swap` - Swap currencies
- `POST /api/wallets/:id/transfer` - Transfer funds
- `GET /api/wallets/:id/transactions` - Transaction history
- `GET /api/transactions/:id` - Get transaction

### FX Rates
- `GET /api/fx-rates` - Get all rates
- `GET /api/fx-rates/:from/:to` - Get specific rate
- `POST /api/fx-rates/refresh` - Refresh rates (admin)

### Explorer
- `GET /api/explorer/transactions` - Recent transactions
- `GET /api/explorer/stats` - Platform stats
- `GET /api/explorer/transactions/:hash` - Get by hash

### Analytics
- `GET /api/wallets/:id/analytics` - Wallet analytics

### AI Assistant (Bonus)
- `POST /api/ai/assistant` - Chat with AI

---

## Environment Variables

```env
# Database
DATABASE_URL=postgresql://user:password@localhost:5432/borderless

# JWT
JWT_SECRET=your-secret-key-here

# FX Rates
USE_MOCK_RATES=false
FX_RATE_API_KEY=your-api-key
FX_RATE_API_URL=https://api.exchangerate-api.com/v4/latest

# AI Assistant (Optional)
OPENAI_API_KEY=sk-...
AI_MODEL=gpt-4
ENABLE_AI_ASSISTANT=true

# CORS
FRONTEND_URL=http://localhost:3000
ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com

# Server
PORT=8080
```

---

## Recommended Implementation Order

1. ✅ **Users & Auth** (Already done)
2. **Database Schema** - Create all tables
3. **Wallets** - Basic wallet CRUD + balances
4. **Deposits** - Simulated deposits
5. **FX Rates** - Rate fetching and storage
6. **Swaps** - Currency swapping
7. **Transfers** - Cross-border transfers
8. **Middleware** - CORS + Audit logging
9. **Explorer** - Transaction explorer
10. **Analytics** - Balance pie charts, stats
11. **AI Assistant** - Bonus feature

---

## Testing Strategy

### Unit Tests
- Service layer logic
- Currency conversion calculations
- Balance validation
- FX rate calculations

### Integration Tests
- API endpoint tests
- Database transactions
- Wallet operations end-to-end

### Test Scenarios
- Deposit → Check balance updated
- Swap → Verify both balances change correctly
- Transfer same currency → Direct transfer
- Transfer different currency → Auto-conversion works
- Insufficient balance → Returns error
- Invalid currency → Returns error

---

## Next Steps

Ready to start? Let's begin with:
1. Creating the database migration file with all tables
2. Implementing the wallets module
3. Setting up CORS middleware

Which would you like to tackle first?

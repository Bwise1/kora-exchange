
## API Endpoints

---
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

### Audit Logs (Protected)
- `POST /api/users/verify-password` - Verify audit access password
- `GET /api/audit-logs?limit=100&offset=0` - Get user's audit logs
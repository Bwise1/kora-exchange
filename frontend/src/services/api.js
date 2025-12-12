const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Generic fetch wrapper with auth
async function fetchWithAuth(endpoint, options = {}) {
  const token = localStorage.getItem('token');

  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  // Get response text first
  const text = await response.text();

  // Try to parse as JSON
  let data;
  try {
    data = text ? JSON.parse(text) : {};
  } catch (error) {
    console.error('Failed to parse JSON response:', text);
    throw new Error('Invalid response from server');
  }

  if (!response.ok) {
    throw new Error(data.message || 'Something went wrong');
  }

  return data;
}

// Auth API
export const authAPI = {
  login: async (email, password) => {
    return fetchWithAuth('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  },
  
  register: async (name, email, password) => {
    return fetchWithAuth('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify({ name, email, password }),
    });
  },
};

// Wallet API
export const walletAPI = {
  getWallet: async () => {
    return fetchWithAuth('/api/wallets');
  },

  getBalances: async () => {
    return fetchWithAuth('/api/wallets/balances');
  },

  getBalance: async (currency) => {
    return fetchWithAuth(`/api/wallets/balance/${currency}`);
  },
};

// Transaction API
export const transactionAPI = {
  deposit: async (currency, amount) => {
    return fetchWithAuth('/api/transactions/deposit', {
      method: 'POST',
      body: JSON.stringify({ currency, amount }),
    });
  },

  getTransactions: async (limit = 50, offset = 0) => {
    return fetchWithAuth(`/api/transactions?limit=${limit}&offset=${offset}`);
  },

  getTransaction: async (id) => {
    return fetchWithAuth(`/api/transactions/${id}`);
  },
};

// FX Rates API
export const fxRatesAPI = {
  getRates: async (baseCurrency = 'USD') => {
    return fetchWithAuth(`/api/fx-rates?base=${baseCurrency}`);
  },

  getRatesForCurrency: async (currency) => {
    return fetchWithAuth(`/api/fx-rates/${currency}`);
  },

  convert: async (from, to, amount) => {
    return fetchWithAuth('/api/fx-rates/convert', {
      method: 'POST',
      body: JSON.stringify({ from, to, amount }),
    });
  },

  refreshRates: async (baseCurrency = 'USD') => {
    return fetchWithAuth(`/api/fx-rates/refresh?base=${baseCurrency}`, {
      method: 'POST',
    });
  },
};

// Mock data for features not yet implemented in backend
export const mockData = {
  transactions: [
    {
      id: '1',
      type: 'DEPOSIT',
      fromCurrency: 'cNGN',
      fromAmount: 50000,
      toCurrency: 'cNGN',
      toAmount: 50000,
      status: 'COMPLETED',
      createdAt: new Date(Date.now() - 3600000).toISOString(),
    },
    {
      id: '2',
      type: 'SWAP',
      fromCurrency: 'cNGN',
      fromAmount: 10000,
      toCurrency: 'USDx',
      toAmount: 6.45,
      status: 'COMPLETED',
      createdAt: new Date(Date.now() - 7200000).toISOString(),
    },
    {
      id: '3',
      type: 'TRANSFER',
      fromCurrency: 'USDx',
      fromAmount: 25,
      toCurrency: 'USDx',
      toAmount: 25,
      status: 'COMPLETED',
      createdAt: new Date(Date.now() - 86400000).toISOString(),
    },
  ],
  
  fxRates: {
    'cNGN/USDx': 0.000645,
    'USDx/cNGN': 1550,
    'cNGN/EURx': 0.00059,
    'EURx/cNGN': 1695,
    'cXAF/USDx': 0.00165,
    'USDx/cXAF': 606,
  },
  
  currencies: [
    { code: 'cNGN', name: 'Nigerian Naira Stablecoin', symbol: '₦' },
    { code: 'cXAF', name: 'CFA Franc Stablecoin', symbol: 'FCFA' },
    { code: 'USDx', name: 'USD Stablecoin', symbol: '$' },
    { code: 'EURx', name: 'EUR Stablecoin', symbol: '€' },
    { code: 'cGHS', name: 'Ghanaian Cedi Stablecoin', symbol: '₵' },
    { code: 'cKES', name: 'Kenyan Shilling Stablecoin', symbol: 'KSh' },
  ],
};

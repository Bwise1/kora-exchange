import React from 'react';

const currencyInfo = {
  cNGN: {
    name: 'cNGN',
    fullName: 'Naira Stablecoin',
    icon: 'NG',
    color: 'bg-green-500',
  },
  cXAF: {
    name: 'cXAF',
    fullName: 'CFA Franc Stablecoin',
    icon: 'XAF',
    color: 'bg-amber-500',
  },
  USDx: {
    name: 'USDx',
    fullName: 'US Dollar Stablecoin',
    icon: '$',
    color: 'bg-blue-600',
  },
  EURx: {
    name: 'EURx',
    fullName: 'Euro Stablecoin',
    icon: '€',
    color: 'bg-blue-400',
  },
  cGHS: {
    name: 'cGHS',
    fullName: 'Cedi Stablecoin',
    icon: 'GHS',
    color: 'bg-yellow-500',
  },
  cKES: {
    name: 'cKES',
    fullName: 'Shilling Stablecoin',
    icon: 'KES',
    color: 'bg-red-500',
  },
};

export default function WalletCard({ currency, balance, usdValue, isLoading }) {
  const info = currencyInfo[currency] || {
    name: currency,
    fullName: 'Stablecoin',
    icon: currency.slice(0, 2),
    color: 'bg-gray-500',
  };

  if (isLoading) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 p-5 animate-pulse">
        <div className="flex items-center gap-3 mb-4">
          <div className="w-12 h-12 bg-gray-200 rounded-xl" />
          <div className="flex-1">
            <div className="h-4 w-16 bg-gray-200 rounded mb-2" />
            <div className="h-3 w-24 bg-gray-200 rounded" />
          </div>
        </div>
        <div className="h-8 w-32 bg-gray-200 rounded mb-1" />
        <div className="h-4 w-24 bg-gray-200 rounded" />
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl border border-gray-200 p-5 hover:border-gray-300 hover:shadow-sm transition-all">
      <div className="flex items-center gap-3 mb-4">
        <div className={`w-12 h-12 ${info.color} rounded-xl flex items-center justify-center shadow-sm`}>
          <span className="text-white font-bold text-sm">{info.icon}</span>
        </div>
        <div>
          <p className="font-semibold text-gray-900">{info.name}</p>
          <p className="text-xs text-gray-500">{info.fullName}</p>
        </div>
      </div>

      <div>
        <p className="text-2xl font-bold text-gray-900 mb-1">
          {typeof balance === 'number'
            ? balance.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
            : '0.00'
          }
        </p>
        <p className="text-sm text-gray-500">
          ≈ ${typeof usdValue === 'number'
            ? usdValue.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
            : '0.00'
          } USD
        </p>
      </div>
    </div>
  );
}

import React from 'react';
import { TrendingUp, RefreshCw } from 'lucide-react';

const currencyInfo = {
  NGN: { name: 'Nigerian Naira', symbol: '₦', flag: '🇳🇬' },
  XAF: { name: 'CFA Franc', symbol: 'FCFA', flag: '🇨🇲' },
  EUR: { name: 'Euro', symbol: '€', flag: '🇪🇺' },
  GHS: { name: 'Ghanaian Cedi', symbol: '₵', flag: '🇬🇭' },
  KES: { name: 'Kenyan Shilling', symbol: 'KSh', flag: '🇰🇪' },
};

export default function ExchangeRatesCard({ rates, lastUpdated, isLoading, onRefresh }) {
  const displayCurrencies = ['NGN', 'XAF', 'EUR', 'GHS', 'KES'];

  const formatRate = (rate) => {
    if (!rate) return '0.00';
    return rate.toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 4
    });
  };

  const formatTime = (timestamp) => {
    if (!timestamp) return 'Never';
    const date = new Date(timestamp);
    const now = new Date();
    const diffMinutes = Math.floor((now - date) / 60000);

    if (diffMinutes < 1) return 'Just now';
    if (diffMinutes < 60) return `${diffMinutes}m ago`;

    const diffHours = Math.floor(diffMinutes / 60);
    if (diffHours < 24) return `${diffHours}h ago`;

    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (isLoading && !rates) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 p-6">
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-2">
            <TrendingUp className="w-5 h-5 text-blue-600" />
            <h3 className="text-lg font-bold text-gray-900">Exchange Rates</h3>
          </div>
        </div>
        <div className="space-y-3 animate-pulse">
          {[1, 2, 3, 4, 5].map((i) => (
            <div key={i} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-gray-200 rounded-full" />
                <div>
                  <div className="h-4 w-24 bg-gray-200 rounded mb-1" />
                  <div className="h-3 w-16 bg-gray-200 rounded" />
                </div>
              </div>
              <div className="h-5 w-20 bg-gray-200 rounded" />
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl border border-gray-200 p-6">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-2">
          <TrendingUp className="w-5 h-5 text-blue-600" />
          <h3 className="text-lg font-bold text-gray-900">Exchange Rates</h3>
        </div>
        <div className="flex items-center gap-3">
          {lastUpdated && (
            <span className="text-xs text-gray-500">
              Updated {formatTime(lastUpdated)}
            </span>
          )}
          {onRefresh && (
            <button
              onClick={onRefresh}
              className="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
              title="Refresh rates"
            >
              <RefreshCw className="w-4 h-4" />
            </button>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <div className="text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">
          1 USD equals
        </div>
        {displayCurrencies.map((currency) => {
          const info = currencyInfo[currency];
          const rate = rates?.[currency];

          return (
            <div
              key={currency}
              className="flex items-center justify-between p-3 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors"
            >
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-white rounded-full flex items-center justify-center shadow-sm border border-gray-200">
                  <span className="text-lg">{info.flag}</span>
                </div>
                <div>
                  <p className="text-sm font-semibold text-gray-900">{currency}</p>
                  <p className="text-xs text-gray-500">{info.name}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-sm font-bold text-gray-900">
                  {info.symbol} {formatRate(rate)}
                </p>
              </div>
            </div>
          );
        })}
      </div>

      <div className="mt-4 pt-4 border-t border-gray-200">
        <p className="text-xs text-gray-500 text-center">
          Powered by FastForex • Updates every 24 hours
        </p>
      </div>
    </div>
  );
}

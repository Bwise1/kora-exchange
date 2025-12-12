import React, { useState, useEffect } from 'react';
import { X, ArrowDownUp } from 'lucide-react';
import { transactionAPI, fxRatesAPI } from '../services/api';

const currencies = [
  { code: 'cNGN', name: 'Nigerian Naira', symbol: '₦' },
  { code: 'cXAF', name: 'CFA Franc', symbol: 'FCFA' },
  { code: 'USDx', name: 'USD Stablecoin', symbol: '$' },
  { code: 'EURx', name: 'EUR Stablecoin', symbol: '€' },
  { code: 'cGHS', name: 'Ghanaian Cedi', symbol: '₵' },
  { code: 'cKES', name: 'Kenyan Shilling', symbol: 'KSh' },
];

// Map frontend currency codes to backend codes
const mapToBackendCurrency = (currency) => {
  const mapping = {
    'cNGN': 'NGN',
    'cXAF': 'XAF',
    'USDx': 'USD',
    'EURx': 'EUR',
    'cGHS': 'GHS',
    'cKES': 'KES',
  };
  return mapping[currency] || currency;
};

export default function SwapModal({ isOpen, onClose, onSuccess, balances }) {
  const [fromCurrency, setFromCurrency] = useState('cNGN');
  const [toCurrency, setToCurrency] = useState('USDx');
  const [amount, setAmount] = useState('');
  const [estimatedAmount, setEstimatedAmount] = useState(0);
  const [exchangeRate, setExchangeRate] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState('');

  // Fetch exchange rate when currencies or amount change
  useEffect(() => {
    const fetchRate = async () => {
      if (!amount || parseFloat(amount) <= 0 || fromCurrency === toCurrency) {
        setEstimatedAmount(0);
        setExchangeRate(0);
        return;
      }

      setIsCalculating(true);
      try {
        const backendFrom = mapToBackendCurrency(fromCurrency);
        const backendTo = mapToBackendCurrency(toCurrency);

        const result = await fxRatesAPI.convert(backendFrom, backendTo, parseFloat(amount));
        setEstimatedAmount(result.data.result);
        setExchangeRate(result.data.rate);
      } catch (err) {
        console.error('Failed to fetch exchange rate:', err);
        setEstimatedAmount(0);
        setExchangeRate(0);
      } finally {
        setIsCalculating(false);
      }
    };

    const timeoutId = setTimeout(fetchRate, 300); // Debounce
    return () => clearTimeout(timeoutId);
  }, [fromCurrency, toCurrency, amount]);

  const handleSwapCurrencies = () => {
    const temp = fromCurrency;
    setFromCurrency(toCurrency);
    setToCurrency(temp);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    // Validation
    if (!amount || parseFloat(amount) <= 0) {
      setError('Please enter a valid amount');
      return;
    }

    if (fromCurrency === toCurrency) {
      setError('Cannot swap the same currency');
      return;
    }

    // Check if user has sufficient balance
    const fromBalance = balances?.[fromCurrency] || 0;
    if (parseFloat(amount) > fromBalance) {
      setError(`Insufficient ${fromCurrency} balance. Available: ${fromBalance.toFixed(2)}`);
      return;
    }

    setIsLoading(true);

    try {
      await transactionAPI.swap(fromCurrency, toCurrency, parseFloat(amount));

      // Reset form
      setAmount('');
      setError('');
      setEstimatedAmount(0);
      setExchangeRate(0);

      // Call success callback
      if (onSuccess) {
        onSuccess();
      }

      // Close modal
      onClose();
    } catch (err) {
      setError(err.message || 'Failed to process swap');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-md mx-4 overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 className="text-2xl font-bold text-gray-900">Swap Currencies</h2>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {/* From Currency */}
          <div>
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              From
            </label>
            <select
              value={fromCurrency}
              onChange={(e) => setFromCurrency(e.target.value)}
              className="input-field"
              disabled={isLoading}
            >
              {currencies.map((curr) => (
                <option key={curr.code} value={curr.code}>
                  {curr.code} - {curr.name}
                </option>
              ))}
            </select>
            <div className="mt-1 text-xs text-gray-500">
              Available: {(balances?.[fromCurrency] || 0).toFixed(2)} {fromCurrency}
            </div>
          </div>

          {/* Amount Input */}
          <div>
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              Amount
            </label>
            <div className="relative">
              <input
                type="number"
                step="0.01"
                min="0.01"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="0.00"
                className="input-field pr-16"
                disabled={isLoading}
                required
              />
              <span className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-500 font-medium">
                {fromCurrency}
              </span>
            </div>
          </div>

          {/* Swap Direction Button */}
          <div className="flex justify-center">
            <button
              type="button"
              onClick={handleSwapCurrencies}
              className="p-3 text-blue-600 hover:bg-blue-50 rounded-full transition-colors"
              disabled={isLoading}
            >
              <ArrowDownUp className="w-5 h-5" />
            </button>
          </div>

          {/* To Currency */}
          <div>
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              To
            </label>
            <select
              value={toCurrency}
              onChange={(e) => setToCurrency(e.target.value)}
              className="input-field"
              disabled={isLoading}
            >
              {currencies.map((curr) => (
                <option key={curr.code} value={curr.code}>
                  {curr.code} - {curr.name}
                </option>
              ))}
            </select>
          </div>

          {/* Estimated Amount */}
          {amount && parseFloat(amount) > 0 && fromCurrency !== toCurrency && (
            <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-700">You'll receive:</span>
                <span className="text-lg font-bold text-gray-900">
                  {isCalculating ? '...' : estimatedAmount.toFixed(4)} {toCurrency}
                </span>
              </div>
              {exchangeRate > 0 && (
                <div className="text-xs text-gray-600">
                  Exchange rate: 1 {fromCurrency} = {exchangeRate.toFixed(6)} {toCurrency}
                </div>
              )}
            </div>
          )}

          {/* Error Message */}
          {error && (
            <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex gap-3">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 btn-secondary"
              disabled={isLoading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 btn-primary"
              disabled={isLoading || isCalculating || fromCurrency === toCurrency}
            >
              {isLoading ? 'Processing...' : 'Swap'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

import React, { useState, useEffect } from 'react';
import { X, Send } from 'lucide-react';
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

export default function SendModal({ isOpen, onClose, onSuccess, balances }) {
  const [recipientAddress, setRecipientAddress] = useState('');
  const [fromCurrency, setFromCurrency] = useState('cNGN');
  const [toCurrency, setToCurrency] = useState('');
  const [amount, setAmount] = useState('');
  const [estimatedAmount, setEstimatedAmount] = useState(0);
  const [exchangeRate, setExchangeRate] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [isCalculating, setIsCalculating] = useState(false);
  const [error, setError] = useState('');
  const [enableConversion, setEnableConversion] = useState(false);

  // Reset toCurrency when conversion is toggled
  useEffect(() => {
    if (!enableConversion) {
      setToCurrency('');
      setEstimatedAmount(0);
      setExchangeRate(0);
    } else {
      // Set to a different currency by default
      const defaultTo = fromCurrency === 'cNGN' ? 'USDx' : 'cNGN';
      setToCurrency(defaultTo);
    }
  }, [enableConversion, fromCurrency]);

  // Fetch exchange rate when currencies or amount change
  useEffect(() => {
    const fetchRate = async () => {
      if (!enableConversion || !toCurrency || !amount || parseFloat(amount) <= 0) {
        setEstimatedAmount(0);
        setExchangeRate(0);
        return;
      }

      if (fromCurrency === toCurrency) {
        setEstimatedAmount(parseFloat(amount));
        setExchangeRate(1);
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
  }, [fromCurrency, toCurrency, amount, enableConversion]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    // Validation
    if (!recipientAddress || recipientAddress.trim() === '') {
      setError('Please enter recipient wallet address');
      return;
    }

    if (!amount || parseFloat(amount) <= 0) {
      setError('Please enter a valid amount');
      return;
    }

    // Check if user has sufficient balance
    const fromBalance = balances?.[fromCurrency] || 0;
    if (parseFloat(amount) > fromBalance) {
      setError(`Insufficient ${fromCurrency} balance. Available: ${fromBalance.toFixed(2)}`);
      return;
    }

    if (enableConversion && fromCurrency === toCurrency) {
      setError('Cannot convert to the same currency');
      return;
    }

    setIsLoading(true);

    try {
      const targetCurrency = enableConversion && toCurrency ? toCurrency : undefined;
      await transactionAPI.transfer(
        recipientAddress.trim(),
        fromCurrency,
        parseFloat(amount),
        targetCurrency
      );

      // Reset form
      setRecipientAddress('');
      setAmount('');
      setError('');
      setEstimatedAmount(0);
      setExchangeRate(0);
      setEnableConversion(false);

      // Call success callback
      if (onSuccess) {
        onSuccess();
      }

      // Close modal
      onClose();
    } catch (err) {
      setError(err.message || 'Failed to process transfer');
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
          <div className="flex items-center gap-2">
            <Send className="w-6 h-6 text-blue-600" />
            <h2 className="text-2xl font-bold text-gray-900">Send Funds</h2>
          </div>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {/* Recipient Address */}
          <div>
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              Recipient Wallet Address
            </label>
            <input
              type="text"
              value={recipientAddress}
              onChange={(e) => setRecipientAddress(e.target.value)}
              placeholder="Enter wallet address"
              className="input-field"
              disabled={isLoading}
              required
            />
            <p className="mt-1 text-xs text-gray-500">
              The unique wallet address of the recipient
            </p>
          </div>

          {/* From Currency */}
          <div>
            <label className="block text-sm font-semibold text-gray-900 mb-2">
              From Currency
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

          {/* Enable Currency Conversion */}
          <div className="flex items-center gap-3 p-4 bg-gray-50 rounded-lg">
            <input
              type="checkbox"
              id="enableConversion"
              checked={enableConversion}
              onChange={(e) => setEnableConversion(e.target.checked)}
              className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              disabled={isLoading}
            />
            <label htmlFor="enableConversion" className="text-sm font-medium text-gray-900 cursor-pointer">
              Convert to a different currency
            </label>
          </div>

          {/* To Currency (only if conversion enabled) */}
          {enableConversion && (
            <div>
              <label className="block text-sm font-semibold text-gray-900 mb-2">
                Recipient Receives
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
          )}

          {/* Estimated Amount */}
          {enableConversion && amount && parseFloat(amount) > 0 && toCurrency && fromCurrency !== toCurrency && (
            <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-700">Recipient will receive:</span>
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

          {/* Same Currency Info */}
          {!enableConversion && amount && parseFloat(amount) > 0 && (
            <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-700">Recipient will receive:</span>
                <span className="text-lg font-bold text-gray-900">
                  {parseFloat(amount).toFixed(2)} {fromCurrency}
                </span>
              </div>
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
              disabled={isLoading || isCalculating}
            >
              {isLoading ? 'Sending...' : 'Send'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

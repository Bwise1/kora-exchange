import React, { useState, useMemo } from 'react';
import { useCurrencies, useFxRates } from '../hooks/useWallet';

export default function SwapWidget({ balances = {} }) {
  const [fromCurrency, setFromCurrency] = useState('cNGN');
  const [toCurrency, setToCurrency] = useState('USDx');
  const [fromAmount, setFromAmount] = useState('');
  const [isSwapping, setIsSwapping] = useState(false);

  const { data: currenciesData } = useCurrencies();
  const { data: ratesData } = useFxRates();

  const currencies = currenciesData?.data || [];
  const rates = ratesData?.data || {};

  const toAmount = useMemo(() => {
    if (!fromAmount || !rates) return '';
    const rateKey = `${fromCurrency}/${toCurrency}`;
    const rate = rates[rateKey];
    if (!rate) return 'Rate unavailable';
    return (parseFloat(fromAmount) * rate).toFixed(4);
  }, [fromAmount, fromCurrency, toCurrency, rates]);

  const handleSwapCurrencies = () => {
    setFromCurrency(toCurrency);
    setToCurrency(fromCurrency);
    setFromAmount('');
  };

  const handleSwap = async () => {
    setIsSwapping(true);
    // Simulate swap - in production this would call the API
    await new Promise(resolve => setTimeout(resolve, 1500));
    setIsSwapping(false);
    setFromAmount('');
    alert('Swap functionality coming soon! Backend endpoint not yet implemented.');
  };

  const availableBalance = balances[fromCurrency] || 0;

  return (
    <div className="card-glass p-6">
      <h3 className="text-lg font-semibold text-white mb-6">Quick Swap</h3>

      {/* From */}
      <div className="space-y-2 mb-4">
        <label className="text-sm text-gray-400">From</label>
        <div className="flex gap-3">
          <select
            value={fromCurrency}
            onChange={(e) => setFromCurrency(e.target.value)}
            className="input-field w-32"
          >
            {currencies.map((c) => (
              <option key={c.code} value={c.code}>{c.code}</option>
            ))}
          </select>
          <div className="flex-1 relative">
            <input
              type="number"
              value={fromAmount}
              onChange={(e) => setFromAmount(e.target.value)}
              placeholder="0.00"
              className="input-field w-full pr-16"
            />
            <button
              onClick={() => setFromAmount(availableBalance.toString())}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-primary-400 hover:text-primary-300 font-medium"
            >
              MAX
            </button>
          </div>
        </div>
        <p className="text-xs text-gray-500">
          Available: {availableBalance.toLocaleString()} {fromCurrency}
        </p>
      </div>

      {/* Swap button */}
      <div className="flex justify-center my-4">
        <button
          onClick={handleSwapCurrencies}
          className="p-3 bg-dark-600 hover:bg-dark-500 rounded-xl text-gray-400 hover:text-white transition-all hover:rotate-180 duration-300"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
          </svg>
        </button>
      </div>

      {/* To */}
      <div className="space-y-2 mb-6">
        <label className="text-sm text-gray-400">To</label>
        <div className="flex gap-3">
          <select
            value={toCurrency}
            onChange={(e) => setToCurrency(e.target.value)}
            className="input-field w-32"
          >
            {currencies.map((c) => (
              <option key={c.code} value={c.code}>{c.code}</option>
            ))}
          </select>
          <input
            type="text"
            value={toAmount}
            readOnly
            placeholder="0.00"
            className="input-field flex-1 bg-dark-800"
          />
        </div>
        {rates[`${fromCurrency}/${toCurrency}`] && (
          <p className="text-xs text-gray-500">
            Rate: 1 {fromCurrency} = {rates[`${fromCurrency}/${toCurrency}`]} {toCurrency}
          </p>
        )}
      </div>

      {/* Swap action button */}
      <button
        onClick={handleSwap}
        disabled={!fromAmount || isSwapping || fromCurrency === toCurrency}
        className="btn-primary w-full flex items-center justify-center gap-2"
      >
        {isSwapping ? (
          <>
            <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            Swapping...
          </>
        ) : (
          'Swap'
        )}
      </button>
    </div>
  );
}

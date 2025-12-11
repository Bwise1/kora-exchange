import React, { useState } from 'react';
import Layout from '../components/Layout';
import WalletCard from '../components/WalletCard';
import TransactionList from '../components/TransactionList';
import DepositModal from '../components/DepositModal';
import { useWallet, useBalances } from '../hooks/useWallet';
import { Eye, Plus, ArrowLeftRight, Send, Download } from 'lucide-react';

export default function Dashboard() {
  const [isDepositModalOpen, setIsDepositModalOpen] = useState(false);
  const { data: walletData, isLoading: walletLoading, refetch: refetchWallet } = useWallet();
  const { data: balancesData, isLoading: balancesLoading, refetch: refetchBalances } = useBalances();

  const balances = balancesData?.data || walletData?.data?.balances || {};
  const isLoading = walletLoading || balancesLoading;

  // Calculate total portfolio value in USD (mock conversion for demo)
  const exchangeRates = {
    cNGN: 0.00065, // 1 NGN ≈ 0.00065 USD
    cXAF: 0.0016,  // 1 XAF ≈ 0.0016 USD
    USDx: 1,       // 1 USDx = 1 USD
    EURx: 1.08,    // 1 EUR ≈ 1.08 USD
    cGHS: 0.082,   // 1 GHS ≈ 0.082 USD
    cKES: 0.0077,  // 1 KES ≈ 0.0077 USD
  };

  const totalUSD = Object.entries(balances).reduce((sum, [currency, balance]) => {
    const rate = exchangeRates[currency] || 0;
    return sum + (balance * rate);
  }, 0);

  const walletCount = Object.keys(balances).length || 4;
  const percentageChange = 2.4; // Mock percentage

  const handleDepositSuccess = () => {
    // Refetch wallet and balances after successful deposit
    refetchWallet();
    refetchBalances();
  };

  return (
    <Layout>
      <div className="max-w-7xl mx-auto">
        {/* Portfolio Value Section */}
        <div className="bg-white rounded-2xl border border-gray-200 p-8 mb-6">
          <div className="flex items-start justify-between">
            <div>
              <div className="flex items-center gap-2 mb-2">
                <p className="text-sm text-gray-600">Total Portfolio Value</p>
                <Eye className="w-4 h-4 text-gray-400" />
              </div>
              <div className="flex items-end gap-3 mb-6">
                <h1 className="text-4xl font-bold text-gray-900">
                  ${isLoading ? '0.00' : totalUSD.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </h1>
                <span className="text-sm text-gray-600 mb-1.5">USD</span>
                <span className="text-sm text-green-600 mb-1.5 flex items-center gap-1">
                  <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z" clipRule="evenodd" />
                  </svg>
                  +{percentageChange}%
                </span>
              </div>

              {/* Action Buttons */}
              <div className="flex flex-wrap gap-3">
                <button
                  onClick={() => setIsDepositModalOpen(true)}
                  className="inline-flex items-center gap-2 px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
                >
                  <Plus className="w-4 h-4" />
                  Deposit
                </button>
                <button className="inline-flex items-center gap-2 px-5 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-900 rounded-lg font-medium transition-colors">
                  <ArrowLeftRight className="w-4 h-4" />
                  Swap
                </button>
                <button className="inline-flex items-center gap-2 px-5 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-900 rounded-lg font-medium transition-colors">
                  <Send className="w-4 h-4" />
                  Send
                </button>
                <button className="inline-flex items-center gap-2 px-5 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-900 rounded-lg font-medium transition-colors">
                  <Download className="w-4 h-4" />
                  Withdraw
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Your Assets Section */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-900">Your Assets</h2>
            <span className="text-sm text-gray-600">{walletCount} Wallets</span>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {Object.entries(balances).length > 0 ? (
              Object.entries(balances).map(([currency, balance]) => (
                <WalletCard
                  key={currency}
                  currency={currency}
                  balance={balance}
                  usdValue={balance * (exchangeRates[currency] || 0)}
                  isLoading={isLoading}
                />
              ))
            ) : (
              // Default currencies
              ['cNGN', 'cXAF', 'USDx', 'EURx'].map((currency) => (
                <WalletCard
                  key={currency}
                  currency={currency}
                  balance={0}
                  usdValue={0}
                  isLoading={isLoading}
                />
              ))
            )}
          </div>
        </div>

        {/* Recent Activity Section */}
        <div className="mb-6">
          <TransactionList />
        </div>

        {/* Deposit Modal */}
        <DepositModal
          isOpen={isDepositModalOpen}
          onClose={() => setIsDepositModalOpen(false)}
          onSuccess={handleDepositSuccess}
        />
      </div>
    </Layout>
  );
}

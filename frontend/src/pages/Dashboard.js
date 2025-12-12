import React, { useState } from "react";
import Layout from "../components/Layout";
import WalletCard from "../components/WalletCard";
import WalletPieChart from "../components/WalletPieChart";
import DepositModal from "../components/DepositModal";
import SwapModal from "../components/SwapModal";
import SendModal from "../components/SendModal";
import ExchangeRatesCard from "../components/ExchangeRatesCard";
import { useWallet, useBalances, useFxRates } from "../hooks/useWallet";
import { fxRatesAPI } from "../services/api";
import { Eye, Plus, ArrowLeftRight, Send, Download } from "lucide-react";

export default function Dashboard() {
  const [isDepositModalOpen, setIsDepositModalOpen] = useState(false);
  const [isSwapModalOpen, setIsSwapModalOpen] = useState(false);
  const [isSendModalOpen, setIsSendModalOpen] = useState(false);
  const [isRefreshingRates, setIsRefreshingRates] = useState(false);
  const {
    data: walletData,
    isLoading: walletLoading,
    refetch: refetchWallet,
  } = useWallet();
  const {
    data: balancesData,
    isLoading: balancesLoading,
    refetch: refetchBalances,
  } = useBalances();
  const {
    data: fxData,
    isLoading: fxLoading,
    refetch: refetchFxRates,
  } = useFxRates("USD");

  const balances = balancesData?.data || walletData?.data?.balances || {};
  const walletAddress = walletData?.data?.wallet_address || "";
  const isLoading = walletLoading || balancesLoading || fxLoading;

  // Get real-time exchange rates from API with fallback
  const apiRates = fxData?.data?.rates || {};

  // Fallback mock rates if API fails (approximate values)
  const mockRates = {
    NGN: 1550, // 1 USD = 1550 NGN
    XAF: 606, // 1 USD = 606 XAF
    EUR: 0.92, // 1 USD = 0.92 EUR
    GHS: 15.2, // 1 USD = 15.20 GHS
    KES: 129.5, // 1 USD = 129.50 KES
  };

  // Use API rates if available, otherwise fallback to mock
  const fxRates = Object.keys(apiRates).length > 0 ? apiRates : mockRates;

  // Map stablecoin codes to their real currency codes for FX API
  const currencyMapping = {
    cNGN: "NGN",
    cXAF: "XAF",
    USDx: "USD",
    EURx: "EUR",
    cGHS: "GHS",
    cKES: "KES",
  };

  // Calculate total portfolio value in USD using real-time rates
  const totalUSD = Object.entries(balances).reduce(
    (sum, [currency, balance]) => {
      const realCurrency = currencyMapping[currency] || currency;
      // If currency is USD, rate is 1, otherwise get rate from API
      const rate = realCurrency === "USD" ? 1 : fxRates[realCurrency] || 0;
      // Convert to USD (rates are from USD base, so we divide for non-USD currencies)
      // Only calculate if rate is valid (not 0)
      const usdValue =
        realCurrency === "USD" ? balance : rate > 0 ? balance / rate : 0;
      return sum + usdValue;
    },
    0
  );

  const walletCount = Object.keys(balances).length || 4;

  const handleDepositSuccess = () => {
    // Refetch wallet and balances after successful deposit
    refetchWallet();
    refetchBalances();
  };

  const handleSwapSuccess = () => {
    // Refetch wallet and balances after successful swap
    refetchWallet();
    refetchBalances();
  };

  const handleSendSuccess = () => {
    // Refetch wallet and balances after successful transfer
    refetchWallet();
    refetchBalances();
  };

  const handleRefreshRates = async () => {
    setIsRefreshingRates(true);
    try {
      await fxRatesAPI.refreshRates("USD");
      await refetchFxRates();
    } catch (error) {
      console.error("Failed to refresh rates:", error);
    } finally {
      setIsRefreshingRates(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-7xl mx-auto">
        {/* Wallet Address Banner */}
        {walletAddress && (
          <div className="bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-2xl p-4 mb-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-medium text-gray-600 mb-1">
                  Your Wallet Address
                </p>
                <p className="text-sm font-mono font-semibold text-gray-900">
                  {walletAddress}
                </p>
              </div>
              <button
                onClick={() => {
                  navigator.clipboard.writeText(walletAddress);
                  // You could add a toast notification here
                }}
                className="px-4 py-2 bg-white hover:bg-gray-50 text-blue-600 text-sm font-medium rounded-lg border border-blue-200 transition-colors"
              >
                Copy
              </button>
            </div>
          </div>
        )}

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
                  $
                  {isLoading
                    ? "0.00"
                    : totalUSD.toLocaleString("en-US", {
                        minimumFractionDigits: 2,
                        maximumFractionDigits: 2,
                      })}
                </h1>
                <span className="text-sm text-gray-600 mb-1.5">USD</span>
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
                <button
                  onClick={() => setIsSwapModalOpen(true)}
                  className="inline-flex items-center gap-2 px-5 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-900 rounded-lg font-medium transition-colors"
                >
                  <ArrowLeftRight className="w-4 h-4" />
                  Swap
                </button>
                <button
                  onClick={() => setIsSendModalOpen(true)}
                  className="inline-flex items-center gap-2 px-5 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-900 rounded-lg font-medium transition-colors"
                >
                  <Send className="w-4 h-4" />
                  Send
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
            {Object.entries(balances).length > 0
              ? Object.entries(balances).map(([currency, balance]) => {
                  const realCurrency = currencyMapping[currency] || currency;
                  const rate =
                    realCurrency === "USD" ? 1 : fxRates[realCurrency] || 0;
                  const usdValue =
                    realCurrency === "USD"
                      ? balance
                      : rate > 0
                      ? balance / rate
                      : 0;

                  return (
                    <WalletCard
                      key={currency}
                      currency={currency}
                      balance={balance}
                      usdValue={usdValue}
                      isLoading={isLoading}
                    />
                  );
                })
              : // Default currencies
                ["cNGN", "cXAF", "USDx", "EURx"].map((currency) => (
                  <WalletCard
                    key={currency}
                    currency={currency}
                    balance={0}
                    usdValue={0}
                    isLoading={isLoading}
                  />
                ))}
          </div>
        </div>

        {/* Wallet Distribution and Exchange Rates Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <div>
            <WalletPieChart balances={balances} isLoading={isLoading} />
          </div>
          <div>
            <ExchangeRatesCard
              rates={fxRates}
              lastUpdated={fxData?.data?.last_updated}
              isLoading={fxLoading || isRefreshingRates}
              onRefresh={handleRefreshRates}
            />
          </div>
        </div>

        {/* Deposit Modal */}
        <DepositModal
          isOpen={isDepositModalOpen}
          onClose={() => setIsDepositModalOpen(false)}
          onSuccess={handleDepositSuccess}
        />

        {/* Swap Modal */}
        <SwapModal
          isOpen={isSwapModalOpen}
          onClose={() => setIsSwapModalOpen(false)}
          onSuccess={handleSwapSuccess}
          balances={balances}
        />

        {/* Send Modal */}
        <SendModal
          isOpen={isSendModalOpen}
          onClose={() => setIsSendModalOpen(false)}
          onSuccess={handleSendSuccess}
          balances={balances}
        />
      </div>
    </Layout>
  );
}

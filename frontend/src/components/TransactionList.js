import React from 'react';
import { useTransactions } from '../hooks/useWallet';
import { ArrowLeftRight, TrendingDown, Send, TrendingUp } from 'lucide-react';

const typeConfig = {
  Swap: {
    icon: ArrowLeftRight,
    bgColor: 'bg-blue-100',
    iconColor: 'text-blue-600',
  },
  Deposit: {
    icon: TrendingDown,
    bgColor: 'bg-green-100',
    iconColor: 'text-green-600',
  },
  Send: {
    icon: Send,
    bgColor: 'bg-orange-100',
    iconColor: 'text-orange-600',
  },
  Withdraw: {
    icon: TrendingUp,
    bgColor: 'bg-purple-100',
    iconColor: 'text-purple-600',
  },
};

const statusConfig = {
  Completed: {
    color: 'text-green-600',
    bgColor: 'bg-green-50',
    dot: 'bg-green-500',
  },
  Pending: {
    color: 'text-orange-600',
    bgColor: 'bg-orange-50',
    dot: 'bg-orange-500',
  },
  Failed: {
    color: 'text-red-600',
    bgColor: 'bg-red-50',
    dot: 'bg-red-500',
  },
};

function formatDate(dateString) {
  const date = new Date(dateString);
  const options = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  };
  return date.toLocaleDateString('en-US', options);
}

// Mock transactions for demo (remove when API is ready)
const mockTransactions = [
  {
    id: '1',
    type: 'Swap',
    description: 'USDx to cNGN',
    details: 'Rate: 1 USDx = 1,560 cNGN',
    date: 'Oct 24, 2023 at 10:46 AM',
    status: 'Completed',
    amount: '-500.00 USDx',
    amountSecondary: '+170,000.00 cNGN',
    isNegative: true,
  },
  {
    id: '2',
    type: 'Deposit',
    description: 'Received EURx from External Wallet',
    details: '',
    date: 'Oct 23, 2023 at 4:20 PM',
    status: 'Completed',
    amount: '+200.00 EURx',
    isNegative: false,
  },
  {
    id: '3',
    type: 'Send',
    description: 'Sent to @SarahDesign',
    details: '',
    date: 'Oct 22, 2023 at 9:15 AM',
    status: 'Pending',
    amount: '-1,200.00 cNGN',
    isNegative: true,
  },
  {
    id: '4',
    type: 'Withdraw',
    description: 'Withdrawal to Bank Account ****4582',
    details: '',
    date: 'Oct 20, 2023 at 11:30 AM',
    status: 'Completed',
    amount: '-500.00 USDx',
    isNegative: true,
  },
];

export default function TransactionList() {
  const { data, isLoading, isError } = useTransactions();

  // Use mock data for now, replace with API data when available
  const transactions = data?.data || mockTransactions;

  if (isLoading) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 p-6">
        <h3 className="text-xl font-bold text-gray-900 mb-6">Recent Activity</h3>
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex items-center gap-4 p-4 animate-pulse">
              <div className="w-10 h-10 bg-gray-200 rounded-xl" />
              <div className="flex-1">
                <div className="h-4 w-32 bg-gray-200 rounded mb-2" />
                <div className="h-3 w-24 bg-gray-200 rounded" />
              </div>
              <div className="h-4 w-20 bg-gray-200 rounded" />
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 p-6">
        <h3 className="text-xl font-bold text-gray-900 mb-4">Recent Activity</h3>
        <p className="text-gray-500">Failed to load transactions</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden">
      <div className="flex items-center justify-between p-6 pb-4">
        <h3 className="text-xl font-bold text-gray-900">Recent Activity</h3>
        <button className="text-sm text-blue-600 hover:text-blue-700 font-medium">
          View All
        </button>
      </div>

      {transactions.length === 0 ? (
        <div className="text-center py-12 px-6">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <p className="text-gray-500">No transactions yet</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Type
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Description
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Date
                </th>
                <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">
                  Amount
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {transactions.map((tx) => {
                const config = typeConfig[tx.type] || typeConfig.Swap;
                const Icon = config.icon;
                const statusStyle = statusConfig[tx.status] || statusConfig.Pending;

                return (
                  <tr key={tx.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center gap-3">
                        <div className={`w-10 h-10 ${config.bgColor} rounded-xl flex items-center justify-center`}>
                          <Icon className={`w-5 h-5 ${config.iconColor}`} />
                        </div>
                        <span className="font-medium text-gray-900">{tx.type}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <div>
                        <p className="text-sm font-medium text-gray-900">{tx.description}</p>
                        {tx.details && (
                          <p className="text-xs text-gray-500 mt-0.5">{tx.details}</p>
                        )}
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className="text-sm text-gray-600">{tx.date}</span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center gap-1.5">
                        <span className={`w-2 h-2 rounded-full ${statusStyle.dot}`}></span>
                        <span className={`text-sm font-medium ${statusStyle.color}`}>
                          {tx.status}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right">
                      <div>
                        <p className={`text-sm font-semibold ${tx.isNegative ? 'text-gray-900' : 'text-green-600'}`}>
                          {tx.amount}
                        </p>
                        {tx.amountSecondary && (
                          <p className="text-xs text-green-600 mt-0.5">
                            {tx.amountSecondary}
                          </p>
                        )}
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

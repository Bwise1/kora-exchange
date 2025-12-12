import React from 'react';
import { PieChart } from 'lucide-react';

const currencyInfo = {
  cNGN: { name: 'Nigerian Naira', symbol: '₦', color: '#3b82f6' },
  cXAF: { name: 'CFA Franc', symbol: 'FCFA', color: '#10b981' },
  USDx: { name: 'USD Stablecoin', symbol: '$', color: '#8b5cf6' },
  EURx: { name: 'EUR Stablecoin', symbol: '€', color: '#f59e0b' },
  cGHS: { name: 'Ghanaian Cedi', symbol: '₵', color: '#ef4444' },
  cKES: { name: 'Kenyan Shilling', symbol: 'KSh', color: '#ec4899' },
};

export default function WalletPieChart({ balances, isLoading }) {
  // Calculate total and percentages
  const totalValue = Object.values(balances).reduce((sum, val) => sum + val, 0);

  const chartData = Object.entries(balances)
    .filter(([_, balance]) => balance > 0)
    .map(([currency, balance]) => ({
      currency,
      balance,
      percentage: totalValue > 0 ? (balance / totalValue) * 100 : 0,
      info: currencyInfo[currency] || { name: currency, symbol: currency, color: '#6b7280' },
    }))
    .sort((a, b) => b.balance - a.balance);

  // Generate pie chart slices using SVG paths
  const generatePieSlices = () => {
    if (chartData.length === 0) return [];

    const size = 200;
    const center = size / 2;
    const radius = size / 2;

    let cumulativePercent = 0;

    return chartData.map((item) => {
      const startAngle = cumulativePercent * 360;
      const endAngle = (cumulativePercent + item.percentage / 100) * 360;

      // Convert angles to radians
      const startRad = (startAngle - 90) * (Math.PI / 180);
      const endRad = (endAngle - 90) * (Math.PI / 180);

      // Calculate coordinates
      const x1 = center + radius * Math.cos(startRad);
      const y1 = center + radius * Math.sin(startRad);
      const x2 = center + radius * Math.cos(endRad);
      const y2 = center + radius * Math.sin(endRad);

      const largeArcFlag = item.percentage > 50 ? 1 : 0;

      const pathData = [
        `M ${center} ${center}`,
        `L ${x1} ${y1}`,
        `A ${radius} ${radius} 0 ${largeArcFlag} 1 ${x2} ${y2}`,
        'Z',
      ].join(' ');

      cumulativePercent += item.percentage / 100;

      return {
        ...item,
        path: pathData,
      };
    });
  };

  const pieSlices = generatePieSlices();

  if (isLoading) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 p-6">
        <div className="flex items-center gap-2 mb-6">
          <PieChart className="w-5 h-5 text-blue-600" />
          <h3 className="text-lg font-bold text-gray-900">Wallet Distribution</h3>
        </div>
        <div className="animate-pulse">
          <div className="w-64 h-64 bg-gray-200 rounded-full mx-auto mb-6" />
          <div className="space-y-2">
            {[1, 2, 3].map((i) => (
              <div key={i} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-3">
                  <div className="w-4 h-4 bg-gray-200 rounded" />
                  <div className="h-4 w-24 bg-gray-200 rounded" />
                </div>
                <div className="h-4 w-16 bg-gray-200 rounded" />
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (chartData.length === 0 || totalValue === 0) {
    return (
      <div className="bg-white rounded-2xl border border-gray-200 p-6">
        <div className="flex items-center gap-2 mb-6">
          <PieChart className="w-5 h-5 text-blue-600" />
          <h3 className="text-lg font-bold text-gray-900">Wallet Distribution</h3>
        </div>
        <div className="flex flex-col items-center justify-center py-12">
          <div className="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
            <PieChart className="w-12 h-12 text-gray-400" />
          </div>
          <p className="text-gray-500 text-center">No balance to display</p>
          <p className="text-sm text-gray-400 text-center mt-1">
            Deposit funds to see your wallet distribution
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl border border-gray-200 p-6">
      <div className="flex items-center gap-2 mb-6">
        <PieChart className="w-5 h-5 text-blue-600" />
        <h3 className="text-lg font-bold text-gray-900">Wallet Distribution</h3>
      </div>

      {/* Pie Chart */}
      <div className="flex justify-center mb-6">
        <svg width="200" height="200" viewBox="0 0 200 200" className="transform rotate-0">
          {pieSlices.map((slice, index) => (
            <path
              key={slice.currency}
              d={slice.path}
              fill={slice.info.color}
              stroke="white"
              strokeWidth="3"
              className="transition-all hover:opacity-80 cursor-pointer"
              style={{
                filter: 'drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1))',
              }}
            />
          ))}
        </svg>
      </div>

      {/* Legend */}
      <div className="space-y-2">
        {chartData.map((item) => (
          <div
            key={item.currency}
            className="flex items-center justify-between p-3 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <div className="flex items-center gap-3">
              <div
                className="w-4 h-4 rounded"
                style={{ backgroundColor: item.info.color }}
              />
              <div>
                <p className="text-sm font-semibold text-gray-900">{item.currency}</p>
                <p className="text-xs text-gray-500">{item.info.name}</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-sm font-bold text-gray-900">
                {item.balance.toLocaleString('en-US', {
                  minimumFractionDigits: 2,
                  maximumFractionDigits: 2,
                })}
              </p>
              <p className="text-xs text-gray-500">
                {item.percentage.toFixed(1)}%
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

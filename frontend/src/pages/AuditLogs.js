import React, { useState } from 'react';
import { Shield, Lock, Eye, EyeOff, Clock, MapPin, Monitor } from 'lucide-react';
import Layout from '../components/Layout';
import { userAPI, auditLogsAPI } from '../services/api';

export default function AuditLogs() {
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isVerifying, setIsVerifying] = useState(false);
  const [isVerified, setIsVerified] = useState(false);
  const [error, setError] = useState('');
  const [auditLogs, setAuditLogs] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  const handleVerify = async (e) => {
    e.preventDefault();
    setError('');
    setIsVerifying(true);

    try {
      await userAPI.verifyPassword(password);
      setIsVerified(true);
      setPassword('');

      // Fetch audit logs after verification
      await fetchAuditLogs();
    } catch (err) {
      setError(err.message || 'Invalid password');
    } finally {
      setIsVerifying(false);
    }
  };

  const fetchAuditLogs = async () => {
    setIsLoading(true);
    try {
      const response = await auditLogsAPI.getAuditLogs(100, 0);
      setAuditLogs(response.data || []);
    } catch (err) {
      setError('Failed to fetch audit logs');
    } finally {
      setIsLoading(false);
    }
  };

  const formatTimestamp = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  const getOperationColor = (operation) => {
    const colors = {
      LOGIN: 'bg-green-100 text-green-800',
      REGISTER: 'bg-blue-100 text-blue-800',
      DEPOSIT: 'bg-purple-100 text-purple-800',
      SWAP: 'bg-orange-100 text-orange-800',
      TRANSFER: 'bg-pink-100 text-pink-800',
      VIEW_WALLET: 'bg-gray-100 text-gray-800',
      VIEW_TRANSACTIONS: 'bg-gray-100 text-gray-800',
    };
    return colors[operation] || 'bg-gray-100 text-gray-800';
  };

  if (!isVerified) {
    return (
      <Layout>
        <div className="max-w-md mx-auto mt-20">
          <div className="bg-white rounded-2xl border border-gray-200 p-8 shadow-sm">
            <div className="flex flex-col items-center mb-6">
              <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-4">
                <Shield className="w-8 h-8 text-blue-600" />
              </div>
              <h1 className="text-2xl font-bold text-gray-900 mb-2">Audit Logs</h1>
              <p className="text-sm text-gray-600 text-center">
                Enter your password to view your security audit logs
              </p>
            </div>

            <form onSubmit={handleVerify} className="space-y-4">
              <div>
                <label className="block text-sm font-semibold text-gray-900 mb-2">
                  Password
                </label>
                <div className="relative">
                  <div className="absolute left-3 top-1/2 -translate-y-1/2">
                    <Lock className="w-5 h-5 text-gray-400" />
                  </div>
                  <input
                    type={showPassword ? 'text' : 'password'}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Enter your password"
                    className="input-field pl-10 pr-12"
                    disabled={isVerifying}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
                  >
                    {showPassword ? (
                      <EyeOff className="w-5 h-5" />
                    ) : (
                      <Eye className="w-5 h-5" />
                    )}
                  </button>
                </div>
              </div>

              {error && (
                <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-sm text-red-600">{error}</p>
                </div>
              )}

              <button
                type="submit"
                className="w-full btn-primary"
                disabled={isVerifying}
              >
                {isVerifying ? 'Verifying...' : 'Verify & View Logs'}
              </button>
            </form>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-2xl border border-gray-200 p-6 mb-6">
          <div className="flex items-center gap-3 mb-2">
            <Shield className="w-6 h-6 text-blue-600" />
            <h1 className="text-2xl font-bold text-gray-900">Security Audit Logs</h1>
          </div>
          <p className="text-sm text-gray-600">
            Track all activities performed on your account for security and compliance
          </p>
        </div>

        {isLoading ? (
          <div className="bg-white rounded-2xl border border-gray-200 p-8">
            <div className="animate-pulse space-y-4">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="h-20 bg-gray-100 rounded-lg" />
              ))}
            </div>
          </div>
        ) : auditLogs.length === 0 ? (
          <div className="bg-white rounded-2xl border border-gray-200 p-12">
            <div className="text-center">
              <Shield className="w-16 h-16 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                No audit logs yet
              </h3>
              <p className="text-gray-600">
                Your activity logs will appear here as you use the platform
              </p>
            </div>
          </div>
        ) : (
          <div className="bg-white rounded-2xl border border-gray-200 overflow-hidden">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Operation
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Time
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      IP Address
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Method
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                      Path
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {auditLogs.map((log) => (
                    <tr key={log.id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-semibold ${getOperationColor(
                            log.operation
                          )}`}
                        >
                          {log.operation}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center gap-2 text-sm text-gray-900">
                          <Clock className="w-4 h-4 text-gray-400" />
                          {formatTimestamp(log.timestamp)}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center gap-2 text-sm text-gray-900 font-mono">
                          <MapPin className="w-4 h-4 text-gray-400" />
                          {log.client_ip}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="text-sm font-medium text-gray-700">
                          {log.request_method}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-2">
                          <Monitor className="w-4 h-4 text-gray-400 flex-shrink-0" />
                          <span className="text-sm text-gray-600 truncate max-w-md">
                            {log.request_path}
                          </span>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
}

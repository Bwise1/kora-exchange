import { useQuery } from '@tanstack/react-query';
import { walletAPI, mockData } from '../services/api';

export function useWallet() {
  return useQuery({
    queryKey: ['wallet'],
    queryFn: walletAPI.getWallet,
    retry: 1,
    staleTime: 30000, // 30 seconds
  });
}

export function useBalances() {
  return useQuery({
    queryKey: ['balances'],
    queryFn: walletAPI.getBalances,
    retry: 1,
    staleTime: 30000,
  });
}

export function useTransactions() {
  // Using mock data until backend endpoint is ready
  return useQuery({
    queryKey: ['transactions'],
    queryFn: async () => {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500));
      return { data: mockData.transactions };
    },
    staleTime: 60000,
  });
}

export function useFxRates() {
  // Using mock data until backend endpoint is ready
  return useQuery({
    queryKey: ['fxRates'],
    queryFn: async () => {
      await new Promise(resolve => setTimeout(resolve, 300));
      return { data: mockData.fxRates };
    },
    staleTime: 60000,
  });
}

export function useCurrencies() {
  return useQuery({
    queryKey: ['currencies'],
    queryFn: async () => {
      await new Promise(resolve => setTimeout(resolve, 200));
      return { data: mockData.currencies };
    },
    staleTime: 300000, // 5 minutes - currencies don't change often
  });
}

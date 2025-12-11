import { useMutation } from '@tanstack/react-query';
import { authAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';

export function useLogin() {
  const { login } = useAuth();
  
  return useMutation({
    mutationFn: ({ email, password }) => authAPI.login(email, password),
    onSuccess: (data) => {
      login(data.data.user, data.data.token);
    },
  });
}

export function useRegister() {
  const { login } = useAuth();
  
  return useMutation({
    mutationFn: ({ name, email, password }) => authAPI.register(name, email, password),
    onSuccess: (data) => {
      // After registration, auto-login if token is returned
      if (data.data?.token) {
        login(data.data.user, data.data.token);
      }
    },
  });
}

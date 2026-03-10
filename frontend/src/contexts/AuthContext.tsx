import React, { createContext, useContext, type ReactNode } from 'react';
import { useAuthStore } from '../stores/useAuthStore';

interface AuthContextType {
  isAuthenticated: boolean;
  token: string | null;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const { isAuthenticated, token, setToken, logout: storeLogout } = useAuthStore();

  const login = (newToken: string) => {
    setToken(newToken);
  };

  const logout = () => {
    storeLogout();
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

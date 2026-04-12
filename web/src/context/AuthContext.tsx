import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from "react";
import { getMe } from "@/lib/api";

interface User {
  subject: string;
  email: string;
  name: string;
  isAdmin: boolean;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  login: (credential: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const validateToken = useCallback(async (t: string) => {
    try {
      const data = await getMe(t);
      setUser({ subject: data.subject, email: data.email, name: data.name, isAdmin: data.is_admin });
      setToken(t);
      localStorage.setItem("auth_token", t);
    } catch {
      setUser(null);
      setToken(null);
      localStorage.removeItem("auth_token");
    }
  }, []);

  useEffect(() => {
    const stored = localStorage.getItem("auth_token");
    if (stored) {
      validateToken(stored).finally(() => setIsLoading(false));
    } else {
      setIsLoading(false);
    }
  }, [validateToken]);

  const login = async (credential: string) => {
    setIsLoading(true);
    await validateToken(credential);
    setIsLoading(false);
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    localStorage.removeItem("auth_token");
  };

  return (
    <AuthContext.Provider value={{ user, token, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}

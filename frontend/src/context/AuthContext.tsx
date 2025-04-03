import { createContext, useContext, useEffect, useState } from "react";
import { authApi } from "../api/auth";
import { User } from "../types/auth"

type AuthContextType = {
    user: User | null;
    isLoading: boolean;
    isAuthenticated: boolean;
    login: (email: string, password: string) => Promise<void>;
    logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType>(null!);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const isAuthenticated = !!user

    useEffect(() => {
        authApi.getProfile()
            .then(setUser)
            .catch(() => setUser(null))
            .finally(() => setIsLoading(false));
    }, []);

    const login = async (email: string, password: string) => {
        setIsLoading(true);
        try {
            const user = await authApi.login({ email, password });
            setUser(user);
        } finally {
            setIsLoading(false);
        }
    };

    const logout = async () => {
        await authApi.logout();
        setUser(null);
    };

    const value = { 
        user,
        isLoading,
        isAuthenticated,
        login,
        logout
    }

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
}

export const useAuth = () => useContext(AuthContext)
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { JSX } from "react";

export function PrivateRoute({ children }: { children: JSX.Element}) {
    const { user, isLoading } = useAuth();
    const location = useLocation();

    if (isLoading) {
        return <div>Загрузка...</div>
    }

    if (!user) {
        return <Navigate to='/login' state={{ from: location }} replace />;
    }

    return children;
}
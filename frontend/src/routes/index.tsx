import { Routes, Route } from "react-router-dom";
import { PrivateRoute } from "./PrivateRoute";
import HomePage from "../pages/HomePage";
import LoginPage from "../pages/auth/LoginPage";
import ExercisesPage from "../pages/exercises";

export function AppRoutes() {
    return (
        <Routes>
            <Route path='/' element={<HomePage />} />
            <Route path='/login' element={<LoginPage />} />

            <Route
                path='/exercises'
                element={
                    <PrivateRoute>
                        <ExercisesPage />
                    </PrivateRoute>
                }
            />
        </Routes>
    );
}

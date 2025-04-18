import { useState, useEffect } from "react";
import { useAuth } from "./useAuth";
import { Workout } from "../models/workouts";

export function useWorkouts() {
    const { user } = useAuth();
    const [workouts, setWorkouts] = useState<Workout[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchWorkouts = async () => {
        try {
            setLoading(true);
            setError(null);

            const response = await fetch("http://localhost:8080/api/v1/workouts", {
                method: "GET",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Ошибка при загрузке тренировок")
            }

            const data = await response.json();
            setWorkouts(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Неизвестная ошибка");
        } finally {
            setLoading(false);
        }
    }

    const updateWorkout = async (id: number, data: { date: string; notes: string }) => {
        try {
            const response = await fetch(`http://localhost:8080/api/v1/workouts/${id}`, {
                method: "PUT",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    date: data.date ? `${data.date}T00:00:00.000Z` : null,
                    notes: data.notes,
                }),
            });

            if (!response.ok) {
                throw new Error("Ошибка при обновлении");
            }

            return true
        } catch (err) {
            setError(err instanceof Error ? err.message : "Неизвестная ошибка");
            return false
        }
    };

    const deleteWorkout = async (id: number) => {
        try {
            const response = await fetch(`http://localhost:8080/api/v1/workouts/${id}`, {
                method: "DELETE",
                credentials: "include",
            });

            if (!response.ok) {
                throw new Error("Ошибка при удалении тренировки");
            }

            fetchWorkouts();
        } catch (err) {
            setError(err instanceof Error ? err.message : "Неизвестная ошибка");
        }
    };

    useEffect(() => {
        if (user) {
            fetchWorkouts();
        }
    }, [user]);

    return {
        workouts,
        loading,
        error,
        fetchWorkouts,
        updateWorkout,
        deleteWorkout,
    }

}

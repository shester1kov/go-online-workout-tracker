import { useState, useEffect } from "react";
import { WorkoutExercise } from "../models/workoutExercise";
import { Workout } from "../models/workouts";

export function useWorkoutDetails(workoutID: number) {
    const [workout, setWorkout] = useState<Workout | null>(null);
    const [exercises, setExercises] = useState<WorkoutExercise[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchExercises = async () => {
        try {
            const response = await fetch(
                `http://localhost:8080/api/v1/workouts/${workoutID}/exercises`, {
                    method: 'GET',
                    credentials: 'include',
                }
            );

            if (response.status === 404) {
                setExercises([]);
                return;
            }
    
            if (!response.ok) throw new Error("Ошибка загрузки упражнений");
            const data = await response.json();
            setExercises(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Неизвестная ошибка");
        }
    };

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);

                const workoutResponse = await fetch(
                    `http://localhost:8080/api/v1/workouts/${workoutID}`, {
                        method: "GET",
                        credentials: 'include'
                    }
                );
                if (!workoutResponse.ok) throw new Error("Ошибка загрузки тренировки");
                
                const workoutData = await workoutResponse.json();
                setWorkout(workoutData);

                await fetchExercises();
            } catch (err) {
                setError(err instanceof Error ? err.message : "Неизвестная ошибка");
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [workoutID]);

    return { workout, exercises, loading, error, fetchExercises }
}
import { useState, useEffect } from "react";
import { useAuth } from "./useAuth";
import { ExerciseList } from "../models/exercise";

export function useExercises() {
    const { user } = useAuth()
    const [exercises, setExercises] = useState<ExerciseList>()
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    const fetchExercises = async () => {
        try {
            setLoading(true)
            setError(null)

            const response = await fetch('http://localhost:8080/api/v1/exercises', {
                method: "GET",
                credentials: 'include'
            })

            if (!response.ok) {
                throw new Error('Ошибка при загрузке упражнений')
            }

            const data = await response.json()
            setExercises(data)
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Неизвестная ошибка')
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        if (user) {
            fetchExercises();
        }
    }, [user])

    return {
        exercises,
        loading,
        error,
        fetchExercises,
    }
}
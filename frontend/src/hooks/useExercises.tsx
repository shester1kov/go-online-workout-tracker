import { useState, useEffect } from "react";
import { useAuth } from "./useAuth";
import { Exercise } from "../models/exercise";

export function useExercises() {
    const { user } = useAuth();
    const [exercises, setExercises] = useState<Exercise[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [pagination, setPagination] = useState({
        page: 1,
        limit: 10,
        total: 0
    });
    const [filters, setFilters] = useState({
        name: '',
        category: '',
        order: 'asc' as 'asc' | 'desc'
    })

    const fetchExercises = async () => {
        try {
            const params = new URLSearchParams({
                page: pagination.page.toString(),
                limit: pagination.limit.toString(),
                name: filters.name,
                category: filters.category,
                order: filters.order
            }); 

            setLoading(true)
            setError(null)

            const response = await fetch(`http://localhost:8080/api/v1/exercises?${params}`, {
                method: "GET",
                credentials: 'include'
            })

            if (!response.ok) {
                throw new Error('Ошибка при загрузке упражнений')
            }

            const data = await response.json()
            setExercises(data.exercises)
            setPagination(prev => ({
                ...prev,
                total: data.total
            }));
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
    }, [user, pagination.page, filters])

    return {
        exercises,
        loading,
        error,
        pagination,
        filters,
        setPage: (page: number) => setPagination(prev => ({ ...prev, page })),
        setFilters,
        refetch: fetchExercises
    }
}
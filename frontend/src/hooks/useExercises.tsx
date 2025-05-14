import { useState, useEffect, useCallback } from "react";
import { useAuth } from "./useAuth";
import { Exercise } from "../models/exercise";
import { API_URL } from "../config";
import { IExerciseFilters } from "../models/exerciseFilters";

export function useExercises() {
  const { user } = useAuth();
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
  });
  const [filters, setFilters] = useState<IExerciseFilters>({
    name: "",
    category: null,
    order: "asc" as "asc" | "desc",
  });

  const debounce = <T extends (...args: never[]) => unknown>(
    fn: T,
    delay: number
  ) => {
    let timeoutId: number;
    return (...args: Parameters<T>) => {
      window.clearTimeout(timeoutId);
      timeoutId = window.setTimeout(() => fn(...args), delay);
    };
  };

  const fetchExercises = useCallback(
    debounce(async (signal?: AbortSignal) => {
      try {
        const params = new URLSearchParams({
          page: pagination.page.toString(),
          limit: pagination.limit.toString(),
          search: filters.name,
          sort_order: filters.order,
        });

        if (filters.category) {
          params.append("category_id", filters.category);
        }

        setLoading(true);
        setError(null);

        const response = await fetch(
          `${API_URL}/exercises?${params.toString()}`,
          {
            method: "GET",
            credentials: "include",
            signal,
          }
        );

        if (response.status === 404) {
          setExercises([]);
          return;
        }

        if (!response.ok) {
          throw new Error("Ошибка при загрузке упражнений");
        }

        const data = await response.json();
        setExercises(data.exercises);
        setPagination((prev) => ({
          ...prev,
          total: data.total,
        }));
      } catch (err) {
        setError(err instanceof Error ? err.message : "Неизвестная ошибка");
      } finally {
        setLoading(false);
      }
    }, 500),
    [
      pagination.page,
      pagination.limit,
      filters.name,
      filters.category,
      filters.order,
    ]
  );

  useEffect(() => {
    if (!user) return;

    const abortController = new AbortController();
    fetchExercises(abortController.signal);

    return () => {
      abortController.abort();
    };
  }, [user, fetchExercises]);

  const totalPages = Math.ceil(pagination.total / pagination.limit);

  return {
    exercises,
    loading,
    error,
    pagination,
    filters,
    setPage: (page: number) => setPagination((prev) => ({ ...prev, page })),
    setFilters,
    refetch: fetchExercises,
    page: pagination.page,
    limit: pagination.limit,
    setLimit: (limit: number) =>
      setPagination((prev) => ({ ...prev, limit, page: 1 })),
    totalPages,
  };
}

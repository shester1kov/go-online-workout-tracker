import { useState, useEffect } from "react";
import { API_URL } from "../config";
import { Category } from "../models/category";
import { useCallback } from "react";

export function useCategories() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchCategories = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${API_URL}/categories`, {
        method: "GET",
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Ошибка при загрузке категорий");
      }

      const data: Category[] = await response.json();
      setCategories(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Неизвестная ошибка");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  const createCategory = async (name: string, description: string) => {
    const response = await fetch(`${API_URL}/categories`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name, description }),
    });
    if (!response.ok) {
      throw new Error("Ошибка при создании категории");
    }
    await fetchCategories();
  };

  const updateCategory = async (
    id: number,
    name: string,
    description: string
  ) => {
    const response = await fetch(`${API_URL}/categories/${id}`, {
      method: "PUT",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name, description }),
    });
    if (!response.ok) {
      throw new Error("Ошибка при обновлении категории");
    }
    await fetchCategories();
  };

  const deleteCategory = async (id: number) => {
    const response = await fetch(`${API_URL}/categories/${id}`, {
      method: "DELETE",
      credentials: "include",
    });
    if (!response.ok) {
      throw new Error("Ошибка при удалении категории");
    }
    await fetchCategories();
  };

  return {
    categories,
    loading,
    error,
    fetchCategories,
    createCategory,
    updateCategory,
    deleteCategory,
  };
}

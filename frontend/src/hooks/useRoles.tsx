import { useState, useEffect } from "react";
import { Role } from "../models/roles";
import { API_URL } from "../config";

export function useRoles(userID: number) {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRoles = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(`${API_URL}/users/${userID}/roles`, {
        method: "GET",
        credentials: "include",
      });

      if (!response.ok) throw new Error("Ошибка загрузки ролей");

      const data = await response.json();
      setRoles(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Неизвестная ошибка");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, [userID]);

  return { roles, loading, error };
}

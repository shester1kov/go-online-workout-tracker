import { useState, useEffect } from 'react';
import { Entry } from '../models/nutrition';
import { API_URL } from '../config';

export function useNutritions(date: string) {
    const [data, setData] = useState<Entry[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const fetchNutritions = async () => {
        try {
            const response = await fetch(`${API_URL}/nutritions/${date}`, {
            method: 'GET',
            credentials: 'include',
        });
        if (!response.ok) throw new Error('Ошибка при загрузке питания');
        const data = await response.json();
        setData(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
        } 
    };

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);

                await fetchNutritions();
            } catch (err) {
                setError(err instanceof Error ? err.message : "Неизвестная ошибка");
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [date]);


    return { data, loading, error };
}
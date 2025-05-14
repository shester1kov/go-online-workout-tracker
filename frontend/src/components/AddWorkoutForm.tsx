import { useState } from "react";
import { API_URL } from "../config";
//import { Workout } from "../models/workouts";

export function AddWorkoutForm({ onSuccess }: { onSuccess: () => void }) {
    const [date, setDate] = useState("");
    const [notes, setNotes] = useState("");

    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            const dateISO = date ? `${date}T00:00:00.000Z` : null

            const response = await fetch(`${API_URL}/workouts`, {
                method: "POST",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    date: dateISO,
                    notes,
                }),
            });

            if (!response.ok) {
                throw new Error("Ошибка при создании тренировки");
            }

            setDate("");
            setNotes("");
            onSuccess();

        } catch (err) {
            setError(err instanceof Error ? err.message : "Неизвестная ошибка");
        }
    };

    return (
        <form onSubmit={handleSubmit} className="mb-6 p-4 bg-gray-100 rounded-lg">
            <h2 className="text-lg font-semibold mb-3">
                Добавить тренировку
            </h2>

            {error && <div className="text-red-500 mb-2">{error}</div>}

            <div className="space-y-3">
                <input
                    type="date"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                    required
                    className="w-full p-2 border rounded"
                 />
                 <input
                    type="text"
                    value={notes}
                    placeholder="Заметки"
                    onChange={(e) => setNotes(e.target.value)}
                    className="w-full p-2 border rounded"
                    required
                 />

                 <button
                    type="submit"
                    className="w-full bg-blue-500 text-white py-2 rounded hover:bg-blue-600"
                 >
                    Добавить
                 </button>
            </div>
        </form>
    );
}
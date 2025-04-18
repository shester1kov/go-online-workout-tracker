import { useState } from "react";
import { Workout } from "../models/workouts";

type Props = {
    workout: Workout;
    onClose: () => void;
    onSave: (id: number, data: { date: string; notes: string }) => Promise<boolean>;
};

export function EditWorkoutForm({ workout, onClose, onSave }: Props) {
    const [date, setDate] = useState(workout.date.split('T')[0]);
    const [notes, setNotes] = useState(workout.notes);
    const [error, setError] = useState<string | null>(null);
    const [isSaving, setIsSaving] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsSaving(true);
        setError(null);

        try {
            const success = await onSave(workout.id, { date, notes });
            if (!success) {
                throw new Error("Не удалось сохранить изменения");
            }
            onClose();
        } catch (err) {
            setError(err instanceof Error ? err.message : "Неизвестная ошибка");
        } finally {
            setIsSaving(false);
        }

    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
                <h2 className="text-xl font-bold mb-4">
                    Редактировать
                </h2>

                {error && <div className="text-red-500 mb-4">{error}</div>}

                <form onSubmit={handleSubmit} className="space-y-4">
                    <input
                        type="date"
                        value={date}
                        onChange={(e) => setDate(e.target.value)}
                        className="w-full p-2 border rounded"
                        required
                    />

                    <input
                        value={notes}
                        onChange={(e) => setNotes(e.target.value)}
                        className="w-full p-2 border rounded"
                    />

                    <div className="flex justify-end space-x-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 border rounded"
                            disabled={isSaving}
                        >
                            Отмена
                        </button>
                        <button
                            type="submit"
                            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
                            disabled={isSaving}
                        >
                            {isSaving ? "Сохранение..." : "Сохранить"}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
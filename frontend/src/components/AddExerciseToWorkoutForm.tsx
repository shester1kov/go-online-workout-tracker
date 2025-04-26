import { useState } from "react";
import { Exercise } from "../models/exercise";

type Props = {
    //workoutID: number;
    allExercises: Exercise[];
    onAdd: (data: {
        sets: number;
        reps: number;
        weight: number;
        notes: string;
        exercise_id: number;
    }) => Promise<boolean>;
    onClose: () => void;
};

export function AddExerciseToWorkoutForm({
    //workoutID,
    allExercises,
    onAdd,
    onClose,
}: Props) {
    const [sets, setSets] = useState(0);
    const [reps, setReps] = useState(0);
    const [weight, setWeight] = useState(0);
    const [notes, setNotes] = useState("");
    const [selectedExerciseID, setSelectedExerciseID] = useState(allExercises[0]?.id || 0);
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        try {
            await onAdd({
                sets,
                reps,
                weight,
                notes,
                exercise_id: selectedExerciseID,
            });
            onClose();
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
                <h2 className="text-xl font-bold mb-4">
                    Добавить упражнение
                </h2>

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block mb-1">
                            Упражнение
                        </label>
                        <select
                            value={selectedExerciseID}
                            onChange={(e) => setSelectedExerciseID(Number(e.target.value))}
                            className="w-full p-2 border rounded"
                            disabled={isLoading}
                        >
                            {allExercises.map((exercise) => (
                                <option key={exercise.id} value={exercise.id}>
                                    {exercise.name}
                                </option>
                            ))}
                        </select>
                    </div>

                    <div className="grid grid-cols-3 gap-4">
                        <div>
                            <label className="block mb-1">
                                Подходы
                            </label>
                            <input
                                type="number"
                                value={sets}
                                onChange={(e) => setSets(Number(e.target.value))}
                                className="w-full p-2 border rounded"
                                min={0}
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block mb-1">
                                Повторения
                            </label>
                            <input
                                type="number"
                                value={reps}
                                onChange={(e) => setReps(Number(e.target.value))}
                                className="w-full p-2 border rounded"
                                min={0}
                                disabled={isLoading}
                            />
                        </div>

                        <div>
                            <label className="block mb-1">
                                Вес (кг)
                            </label>
                            <input
                                type="number"
                                value={weight}
                                onChange={(e) => setWeight(parseFloat(e.target.value))}
                                className="w-full p-2 border rounded"
                                min={0}
                                disabled={isLoading}
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block mb-1">Заметки</label>
                        <input
                            value={notes}
                            onChange={(e) => setNotes(e.target.value)}
                            className="w-full p-2 border rounded"
                            disabled={isLoading}
                        />
                    </div>

                    <div className="flex justify-end space-x-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 border rounded"
                            disabled={isLoading}
                        >
                            Отмена
                        </button>
                        
                        <button
                            type="submit"
                            disabled={isLoading}
                            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
                        >
                            {
                                isLoading ? "Добавление..." : "Добавить упражнение"
                            } 
                        </button> 
                    </div>
 
                </form>
            </div>
        </div>
    )
}
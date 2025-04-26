import { useState } from "react";
import { WorkoutExercise } from "../models/workoutExercise";
import { Exercise } from "../models/exercise";

type Props = {
    exercise: WorkoutExercise;
    onSave: (data: {
        sets: number;
        reps: number;
        weight: number;
        notes: string;
        exercise_id: number;
    }) => Promise<boolean>;
    onClose: () => void;
    allExercises: Exercise[];
};

export function EditWorkoutExerciseForm({
    exercise,
    onSave,
    onClose,
    allExercises
}: Props) {
    const [sets, setSets] = useState(exercise.sets);
    const [reps, setReps] = useState(exercise.reps);
    const [weight, setWeight] = useState(exercise.weight);
    const [notes, setNotes] = useState(exercise.notes);
    const [selectedExerciseID, setSelectedExerciseID] = useState(exercise.exercise_id);
    const [isLoading, setIsLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        try {
            await onSave({
                sets,
                reps,
                weight,
                notes,
                exercise_id: selectedExerciseID
            });
            onClose();
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4">
            <div className="bg-white rounded-lg p-6 w-full max-w-md">
                <h2 className="text-xl font-bold mb-4">
                    Редактировать упражнение
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
                                onChange={(e) => setWeight(Number(e.target.value))}
                                className="w-full p-2 border rounded"
                                disabled={isLoading}
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block mb-1">
                            Заметки
                        </label>
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
                            type='submit'
                            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
                            disabled={isLoading}
                        >
                            {isLoading ? 'Сохранение...' : 'Сохранить'}
                        </button>
                    </div>

                </form>
            </div>
        </div>
    )
}
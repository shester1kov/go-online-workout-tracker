import { useState } from "react";
import { useParams, Link } from "react-router-dom";
import { useWorkoutDetails } from "../hooks/useWorkoutDetails";
import { useExercises } from "../hooks/useExercises";
import { WorkoutExercise } from "../models/workoutExercise";
import { EditWorkoutExerciseForm } from "../components/EditWorkoutExerciseForm";
import { AddExerciseToWorkoutForm } from "../components/AddExerciseToWorkoutForm";

export default function WorkoutDetails() {
    const { id } = useParams<{ id: string }>();
    const workoutID = parseInt(id || "0");
    const { workout, exercises, loading, error, fetchExercises } = useWorkoutDetails(workoutID);
    const { exercises: allExercises } = useExercises();
    const [editingExercise, setEditingExercise] = useState<WorkoutExercise | null>(null);
    const [addingExercise, setAddingExercise] = useState(false);
    const [operationError, setOperaionError] = useState<string | null>(null);

    const handleAddExercise = async (data: {
        sets: number;
        reps: number;
        weight: number;
        notes: string;
        exercise_id: number;
    }) => {
        try {
            setOperaionError(null);

            const response = await fetch(
                `http://localhost:8080/api/v1/workouts/${workoutID}/exercises`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json'},
                    credentials: 'include',
                    body: JSON.stringify(data),
                }
            );

            if (!response.ok) throw new Error('Ошибка добавления');
            await fetchExercises();
            setAddingExercise(false);
            return true;
        } catch (err) {
            setOperaionError(err instanceof Error ? err.message : "Неизвестная ошибка");
            return false;
        }
    };


    const handleUpdateExercise = async (data: {
        sets: number;
        reps: number;
        weight: number;
        notes: string;
        exercise_id: number;
    }) => {
        if (!editingExercise) return false;

        try {
            setOperaionError(null);

            const response = await fetch(
                `http://localhost:8080/api/v1/workouts/${workoutID}/exercises/${editingExercise.id}`,
                {
                    method: 'PUT',
                    credentials: 'include',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(data),
                }
            );

            if (!response.ok) throw new Error("Ошибка обновления");
            await fetchExercises();
            return true;
        } catch (err) {
            setOperaionError(err instanceof Error ? err.message : "Неизвестная ошибка");
            return false;
        }
    }

    const handleDeleteExercises = async (exerciseID: number) => {
        if (!window.confirm("Вы уверены что хотите удалить это упражнение?")) {
            return;
        }

        try {
            setOperaionError(null);

            const response = await fetch(
                `http://localhost:8080/api/v1/workouts/${workoutID}/exercises/${exerciseID}`, {
                    method: 'DELETE',
                    credentials: 'include',
                }
            );

            if (!response.ok) throw new Error("Ошибка удаления упражнения");
            await fetchExercises();
        } catch (err) {
            setOperaionError(err instanceof Error ? err.message : "Неизвестная ошибка");
        }
    };

    if (loading) return <div>Загрузка...</div>;
    if (error) return <div className="text-red-500">{error}</div>;
    if (!workout) return <div>Тренировка не найдена</div>

    const dateOnly = workout.date.split("T")[0]

    return (
        <div className="container mx-auto p-4">
            <Link to="/workouts" className="text-blue-500 hover:underline mb-4 inline-block">
                Назад к списку тренировок
            </Link>

            <h1 className="text-2xl font-bold mb-2">
                Тренировка {dateOnly}
            </h1>

            <div className="bg-white rounded-lg shadow p-4 mb-6">
                <h2 className="text-lg font-semibold mb-2">
                    Описание
                </h2>

                <p className="text-gray-600">
                    {workout.notes || <span className="text-gray-400">Нет описания</span>}
                </p>
            </div>

            {operationError && (
                <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-4">
                    {operationError}
                </div>
            )}

            <h2 className="text-xl font-bold mb-4">Упражнения</h2>

            <button
                onClick={() => setAddingExercise(true)}
                className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 mb-4"
            >
                Добавить упражнение
            </button>

            {exercises.length === 0 ? (
                <p className="text-gray-500">Нет упражнений в этой трениировке</p>
            ) : (
                <div className="grid gap-4">
                    {exercises.map((exercise) => (
                        <div key={exercise.id} className="bg-white rounded-lg shadow p-4 mb-4">
                            <h3 className="font-bold text-lg">
                                {exercise.exercise?.name}
                            </h3>

                            <div className="flex space-x-2">
                                <button
                                    onClick={() => setEditingExercise(exercise)}
                                    className="text-blue-500 hover:text-blue-700"
                                >
                                    Редактировать
                                </button>

                                <button
                                    onClick={() => handleDeleteExercises(exercise.id)}
                                    className="text-red-500 hover:text-red-700"
                                >
                                    Удалить
                                </button>
                            </div>

                            <p className="text-gray-600">
                                {exercise.exercise?.description}
                            </p>

                            <div className="mt-2 grid grid-cols-3 gap-2">
                                <div>
                                    <span className="text-sm text-gray-500">
                                        Подходы
                                    </span>
                                    <p>{exercise.sets}</p>
                                </div>
                                <div>
                                    <span className="text-sm text-gray-500">
                                        Повторения
                                    </span>
                                    <p>{exercise.reps}</p>
                                </div>
                                <div>
                                    <span className="text-sm text-gray-500">
                                        Вес (кг)
                                    </span>
                                    <p>{exercise.weight}</p>
                                </div>
                            </div>
                        </div>
                    ))}
                    {
                        editingExercise && (
                            <EditWorkoutExerciseForm
                                exercise={editingExercise}
                                onSave={handleUpdateExercise}
                                onClose={() => setEditingExercise(null)}
                                allExercises={allExercises}
                            />
                        )
                    }

                </div>
            )}

            {
                addingExercise && (
                    <AddExerciseToWorkoutForm
                        // workoutID={workoutID}
                        onAdd={handleAddExercise}
                        onClose={() => setAddingExercise(false)}
                        allExercises={allExercises}
                    />
                )
            }

        </div>
    );
}
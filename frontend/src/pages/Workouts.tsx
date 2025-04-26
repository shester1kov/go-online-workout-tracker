import { useState } from "react";
import { useWorkouts } from "../hooks/useWorkouts";
import { Workout } from "../models/workouts";
import { AddWorkoutForm } from "../components/AddWorkoutForm";
import { EditWorkoutForm } from "../components/EditWorkuotForm";
import { Link } from "react-router-dom";

export default function Workouts() {
    const { workouts, loading, error, fetchWorkouts, updateWorkout, deleteWorkout } = useWorkouts();

    const [editingWorkout, setEditingWorkout] = useState<Workout | null>(null);

    if (loading) return <div>Загрузка...</div>;
    if (error) return <div className="text-red-500">{error}</div>;

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-6">Мои тренировки</h1>

            <AddWorkoutForm onSuccess={() => fetchWorkouts()} />

            {workouts.length === 0 ? (
                <p className="text-gray-500">
                    Нет тренировок
                </p>
            ) : (
                <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                    {workouts.map((workout) => (
                        <WorkoutCard
                            key={workout.id}
                            workout={workout}
                            onEdit={() => setEditingWorkout(workout)}
                            onDelete={deleteWorkout}
                        />
                    ))}
                </div>
            )}


            {editingWorkout && (
                <EditWorkoutForm
                    workout={editingWorkout}
                    onClose={() => setEditingWorkout(null)}
                    onSave={async (id, data) => {
                        const success = await updateWorkout(id, data);
                        if (success) fetchWorkouts();
                        return success;
                    }}
                />
            )}
        </div>
    );
}

function WorkoutCard({
    workout,
    onEdit,
    onDelete
}: {
    workout: Workout;
    onEdit: () => void;
    onDelete: (id: number) => void;
}) {
    const dateOnly = workout.date.split("T")[0]

    return (
        <div className="bg-white rounded-lg shadow p-4">
            <Link to={`/workouts/${workout.id}`} className="block">
                <h2 className="font-bold text-lg hover:text-blue-500">
                    Тренировка {dateOnly}
                </h2>

                <p className="text-gray-600 mt-2">
                    {workout.notes || <span className="text-gray-400">Нет описания</span>}
                </p>
            </Link>

            <div className="flex mt-2">
                <button
                    onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        onEdit();
                    }}
                    className="mt-2 bg-blue-500 text-white py-1 px-3 rounded hover:bg-blue-600 mr-1"
                >
                    Редактировать
                </button>

                <button
                    onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        onDelete(workout.id)
                    }}
                    className="mt-2 bg-red-500 text-white py-1 px-3 rounded hover:bg-red-600"
                >
                    Удалить
                </button>
            </div>

        </div>
    )
}
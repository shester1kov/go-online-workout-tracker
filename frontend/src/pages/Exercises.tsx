import { useExercises } from "../hooks/useExercises";
import { Exercise } from "../models/exercise";

export default function Exercises() {
    const { exercises, loading, error } = useExercises()

    if (loading) return <div>Загрузка...</div>
    if (error) return <div className="text-red-500">{error}</div>

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-6">
                Упражнения
            </h1>
            <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                {exercises.map(exercise => (
                    <ExerciseCard key={exercise.id} exercise={exercise} />
                ))}
            </div>
        </div>
    )
}

function ExerciseCard({ exercise }: { exercise: Exercise }) {
    return (
        <div className="bg-white rounded-lg shadow p-4">
            <h2 className="font-bold text-lg">
                {exercise.name}
            </h2>
            <p className="text-gray-600 mt-2">
                {exercise.description}
            </p>
        </div>
    )
}
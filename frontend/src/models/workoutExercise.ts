export interface WorkoutExerciseItem {
    id: number;
    name: string;
    description: string;
}

export interface WorkoutExercise {
    id: number;
    workout_id: number;
    exercise_id: number;
    sets: number;
    reps: number;
    weight: number;
    notes: string;
    created_at: string;
    exercise?: WorkoutExerciseItem;
}
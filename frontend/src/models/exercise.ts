export interface Exercise {
    id: number,
    name: string,
    description: string,
    category_id: number,
    created_at: string,
    updated_at: string,
}

export interface ExerciseList {
    exercises: Exercise[],
    total: number,
    page: number,
    limit: number,
}

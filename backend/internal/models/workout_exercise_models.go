package models

import "time"

type WorkoutExercise struct {
	ID         int                  `json:"id"`
	WorkoutID  int                  `json:"workout_id"`
	ExerciseID int                  `json:"exercise_id"`
	Sets       int                  `json:"sets"`
	Reps       int                  `json:"reps"`
	Weight     float64              `json:"weight"`
	Notes      string               `json:"notes"`
	CreatedAt  time.Time            `json:"created_at"`
	Exercise   *WorkoutExerciseItem `json:"exercise,omitempty"`
}

type WorkoutExerciseRequest struct {
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Notes      string  `json:"notes"`
}

type WorkoutExerciseResponse struct {
	ID         int                  `json:"id"`
	WorkoutID  int                  `json:"workout_id"`
	ExerciseID int                  `json:"exercise_id"`
	Sets       int                  `json:"sets"`
	Reps       int                  `json:"reps"`
	Weight     float64              `json:"weight"`
	Notes      string               `json:"notes"`
	CreatedAt  time.Time            `json:"created_at"`
	Exercise   *WorkoutExerciseItem `json:"exercise,omitempty"`
}

type WorkoutExerciseItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

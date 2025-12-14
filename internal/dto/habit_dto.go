package dto

type CreateHabitRequest struct {
	UserID uint64 `json:"user_id"`
	Name   string `json:"name"`
	Info   string `json:"info"`
}

type UpdateHabitRequest struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

type DeleteHabitRequest struct {
	UserID  uint64 `json:"user_id"`
	HabitID uint64 `json:"habit_id"`
}

package dto

type CreateHabitRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Info   string `json:"info" binding:"required"`
}

type UpdateHabitRequest struct {
	Name string `json:"name" binding:"required"`
	Info string `json:"info" binding:"required"`
}

type DeleteHabitRequest struct {
	UserID  uint64 `json:"user_id" binding:"required"`
	HabitID uint64 `json:"habit_id" binding:"required"`
}

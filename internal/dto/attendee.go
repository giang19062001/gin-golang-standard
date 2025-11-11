package dto

type CreateAttendeeDto struct {
	UserId  int `json:"userId" binding:"required"`
	EventId int `json:"eventId" binding:"required"`
}

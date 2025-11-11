package models

type Event struct {
	Id          int        `json:"id"`
	OwnerId     int        `json:"ownerId" binding:"required"`
	Name        string     `json:"name" binding:"required,min=3"`
	Description string     `json:"description" binding:"required,min=10"`
	Date        string     `json:"date" binding:"required,datetime=2006-01-02"` // YYYY-MM-DD
	Location    string     `json:"location" binding:"required,min=3"`
	EventImgs   []EventImg `json:"eventImgs,omitempty"` // * Khi EventImgs == nil hoặc len(EventImgs) == 0, thì JSON sẽ bỏ qua field này
}

type EventImg struct {
	Id       int    `json:"id"`
	EventId  int    `json:"eventId" binding:"required"`
	Filepath string `json:"filepath" binding:"required"`
}

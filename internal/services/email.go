package services

import (
	"encoding/json"
	"fmt"

	"github.com/giang19062001/gin-golang-standard/internal/models"
	rabbitmq "github.com/giang19062001/gin-golang-standard/pkg/rabbit"
)

type IEmailService interface {
	SendEmailHandler(email, name string) error
}

type EmailService struct {
	rabbitMq *rabbitmq.MqService
}

func NewEmailService(rabbitMq *rabbitmq.MqService) IEmailService {
	return &EmailService{rabbitMq: rabbitMq}
}

func (ser *EmailService) SendEmailHandler(email, name string) error {
	fmt.Printf("Đang chuẩn bị gửi email %s .... \n", email)

	msg := models.EmailMessage{
		To:      email,
		Subject: fmt.Sprintf("Hello *%s*", name),
		Body:    "Đơn đăng ký sự kiện của bạn đã thành công.",
	}

	// chuyển struct msg thành chuỗi JSON
	payload, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("marshal json lỗi: %s", err)
		return err
	}

	// gọi rabbit để đẩy nội dung sự kiện (send email)
	err = ser.rabbitMq.PushEmailEvent(email, payload)
	if err != nil {
		return err
	}
	return nil
}

package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/giang19062001/gin-golang-standard/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IEmailService interface {
	SendEmailHandler(email string) error
	Close()
}

type EmailService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewEmailService() (IEmailService, error) {
	log.Print("Đang kết nối RabbitMQ.....")

	// Thay đổi URL nếu bạn dùng user/pass hoặc port khác
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("lỗi mở channel RabbitMQ: %w", err)
	}

	// Khai báo queue (tên giống topic cũ để dễ theo dõi)
	queueName := "send-email"
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("lỗi khai báo queue: %w", err)
	}

	fmt.Print("Kết nối RabbitMQ thành công")
	return &EmailService{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

func (s *EmailService) SendEmailHandler(email string) error {
	fmt.Printf("Email sẽ được gửi qua RabbitMQ: %s", email)

	emailArray := []string{email}

	msg := models.EmailMessage{
		To:      emailArray,
		Subject: "Chào mừng",
		Body:    "Thành viên mới",
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal json lỗi: %w", err)
	}

	// Publish vào queue
	err = s.channel.Publish(
		"",      // exchange (rỗng = default exchange)
		s.queue, // routing key = tên queue
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent, // lưu vào disk để không mất khi restart
		},
	)
	if err != nil {
		return fmt.Errorf("gửi message vào RabbitMQ lỗi: %w", err)
	}

	fmt.Printf("Đã gửi message vào queue %s thành công", s.queue)
	return nil
}

// Đóng kết nối khi app shutdown
func (s *EmailService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	fmt.Print("Đã đóng kết nối RabbitMQ")
}

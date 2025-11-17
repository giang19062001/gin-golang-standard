package rabbitmq

import (
	"errors"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MqService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func Init(rabbitUrl string) (*MqService, error) {
	log.Print("Đang kết nối RabbitMQ.....")

	// Thay đổi URL nếu bạn dùng user/pass hoặc port khác
	conn, err := amqp.Dial(rabbitUrl)
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối RabbitMQ: %w", err)
	}

	// mở 1 channel để giao tiếp
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("lỗi mở channel RabbitMQ: %w", err)
	}

	// Khai báo tên queue
	queueName := "send-email"
	_, err = ch.QueueDeclare(
		queueName, // tên queue
		true,      // durable -> true: queue sẽ được lưu trên disk, không mất khi RabbitMQ restart.
		false,     // false: queue không bị xóa tự động khi không có consumer
		false,     // exclusive -> false: queue không bị giới hạn chỉ cho 1 connection, nhiều consumer có thể dùng
		false,     // no-wait -> false: phải chờ RabbitMQ phản hồi việc khai báo queue
		nil,       // arguments -> nil : không truyền thêm option đặc biệt
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("lỗi khai báo queue: %w", err)
	}

	fmt.Print("Kết nối RabbitMQ thành công")
	return &MqService{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}
func (ser *MqService) PushEmailEvent(email string, payload []byte) error {

	if ser.channel == nil || ser.channel.IsClosed() {
		fmt.Println("channel RabbitMQ đã bị đóng hoặc chưa khởi tạo")
		return errors.New("channel RabbitMQ đã bị đóng hoặc chưa khởi tạo")
	}

	// * Các loại exchange trong RabbitMq
	// ? Direct	-> Routing key khớp chính xác	Logging theo mức độ
	// ? Fanout	-> Broadcast đến tất cả queue	Thông báo hệ thống, chat broadcast
	// ? Topic	-> Pattern matching với * và #	Phân loại theo chủ đề, module
	// ? Headers -> Dựa trên header key-value

	// Publish vào queue
	err := ser.channel.Publish(
		"",        // "" → dùng default exchange -> Direct
		ser.queue, // routing key = tên queue
		false,     // false → không bắt buộc phải có queue nhận.
		false,     // false → không yêu cầu consumer phải nhận ngay.
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent, // lưu vào disk để không mất khi restart
		},
	)
	if err != nil {
		fmt.Printf("gửi message vào RabbitMQ lỗi: %s", err)
		return err
	}

	fmt.Printf("Đã gửi message vào queue %s thành công", ser.queue)
	return nil
}

func (s *MqService) Close() {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	fmt.Print("Đã đóng kết nối RabbitMQ")
}

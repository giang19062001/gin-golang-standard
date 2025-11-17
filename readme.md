# tạo ra file cấu hình hot reload
air init 

# cài swagger CLI nếu chưa có
go install github.com/swaggo/swag/cmd/swag@latest

# chạy swagger
swag init --dir ./ --generalInfo swagger.go --parseDependency --parseInternal --parseDepth 2 -o ./docs

# khởi tạo go project
 go mod init github.com/giang19062001/gin-golang 

# cài đặt global cho migrate CLI
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# khởi tạo 1 migration cho 1 table
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_users_table

# cách chạy migration
go run cmd/migrate/main.go up
 or
go run cmd/migrate/main.go down

# khởi tạo rabbimq
docker-compose up -d

# tắt rabbimq
docker-compose down -v
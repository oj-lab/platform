build:
	go build -o bin/migrate_user.exe migration/migrate_user.go
	go build -o bin/user_service.exe user-service/application.go
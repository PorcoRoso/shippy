package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/labstack/gommon/log"
	pb "github.com/porcorosso/shippy/user-service/proto/user"
)

func main() {
	// 连接到数据库
	db, err := CreateConnection()
	fmt.Printf("%+v\n", db)
	fmt.Printf("%+v\n", err)

	defer db.Close()
	if err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	repo := &UserRepository{db}
	// 自动检查 User 结构是否变化
	db.AutoMigrate(&pb.User{})
	s := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	s.Init()

	t := TokenService{repo}
	pb.RegisterUserServiceHandler(s.Server(), &handler{repo, &t})

	if err := s.Run(); err != nil {
		log.Fatalf("user service error: %v\n", err)
	}

}



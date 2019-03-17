package main

import (
	"github.com/micro/go-micro"
	"log"
	pb "github.com/porcorosso/shippy/user-service/proto/user"
)

func main() {
	// 连接到数据库
	db, err := CreateConnection()
	log.Printf("%v\n", db)
	log.Printf("%v\n", err)

	defer db.Close()
	if err != nil {
		log.Fatalf("connect error: %v\n", err)
	}

	repo := &UserRepository{db}
	// 自动检查 User 结构是否变化
	db.AutoMigrate(&pb.User{})

	srv := micro.NewService(
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)

	srv.Init()

	// 获取broker实例
	//pubSub := srv.Server().Options().Broker
	publisher := micro.NewPublisher(topic, srv.Client())

	// 注册handler
	t := TokenService{repo}
	//pb.RegisterUserServiceHandler(srv.Server(), &handler{repo, &t, pubSub})
	pb.RegisterUserServiceHandler(srv.Server(), &handler{repo, &t, publisher})
	if err := srv.Run(); err != nil {
		log.Fatalf("user service error: %v\n", err)
	}

}



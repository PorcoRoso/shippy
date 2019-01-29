package main

import (
	"log"
	"os"
	"golang.org/x/net/context"
	pb "github.com/porcorosso/shippy/user-service/proto/user"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
)

func main() {

	cmd.Init()

	// 创建 user-service 微服务的客户端
	client := pb.NewUserServiceClient("go.micro.srv.user", microclient.DefaultClient)


	r, err := client.Create(context.TODO(), &pb.User{
		Name: "evan",
		Email: "evan@bbc",
		Password: "test123",
		Company: "bbc",
	})
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %v", r.User.Id)

	getAll, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}

	authResp, err := client.Auth(context.TODO(), &pb.User{
		Email: "evan@bbc",
		Password: "test123",
	})

	if err != nil {
		log.Fatalf("auth failed: %v", err)
	}
	log.Println("token:", authResp.Token)

	os.Exit(0)


}


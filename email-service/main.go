package main

import (
	"github.com/micro/go-micro"
	"log"
	pb "github.com/porcorosso/shippy/user-service/proto/user"
	_ "github.com/micro/go-plugins/broker/nats"
	"context"
)

const topic = "user.created"

type Subscriber struct {}

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.email"),
		micro.Version("latest"),
	)

	srv.Init()

	/*pubSub := srv.Server().Options().Broker
	if err := pubSub.Connect(); err != nil {
		log.Fatalf("broker connect error : %v\n", err)
	}

	_, err := pubSub.Subscribe(topic, func(pub broker.Publication) error {
		var user *pb.User
		if err := json.Unmarshal(pub.Message().Body, &user); err != nil {
			log.Println("unmarshal error: %v", err)
			return err
		}
		log.Printf("[Create User]: %v\n", user)
		go sendEmail(user)

		return nil
	})

	if err != nil {
		log.Printf("subscribe error : %v\n", err)
	}*/

	micro.RegisterSubscriber(topic, srv.Server(), new(Subscriber))

	if err := srv.Run(); err != nil {
		log.Fatalf("srv run error: %v\n", err)
	}
}

func (sub *Subscriber) Process(ctx context.Context, user *pb.User) error {
	log.Println("[Picked up a new message]")
	log.Println("[Sending email to]:", user.Name)
	return nil
}

func sendEmail(user *pb.User) error{
	log.Printf("Sending a email to %s ...", user.Name)
	return nil
}

package main

import (
	"context"
	pb "github.com/porcorosso/shippy/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"fmt"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/micro/go-micro"
)

const topic = "user.created"

type handler struct {
	repo Repository
	tokenService Authable
	Publisher micro.Publisher
	//PubSub broker.Broker
}

func (h *handler) Get(ctx context.Context, req *pb.User, resp *pb.Response) error {
	u, err := h.repo.Get(req.Id)
	if err != nil {
		return err
	}

	resp.User = u
	return nil
}

func (h *handler) GetAll(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	u, err := h.repo.GetAll()
	if err != nil {
		return err
	}
	resp.Users = u
	return nil
}

func (h *handler) Create(ctx context.Context, req *pb.User, resp *pb.Response) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPwd)
	err = h.repo.Create(req)
	if err != nil {
		return err
	}
	resp.User = req

	// 发布消息，消息包括用户信息，通知邮件模块
	if err := h.Publisher.Publish(ctx, req); err != nil {
		return err
	}
	return nil
}

// 发送消息
/*func (h *handler) publishEvent(user *pb.User) error {
	body, err := json.Marshal(user)
	if err != nil {
		log.Println("marshal user info error : ", err)
		return err
	}

	msg := &broker.Message{
		Header: map[string]string {
			"id": user.Id,
		},
		Body: body,
	}

	// 发布user.created topic消息
	if err := h.PubSub.Publish(topic, msg); err != nil {
		log.Fatalf("[pub] failed: %v\n", err)
	}

	log.Printf("publish user.created topic success. msg is : %v\n", user)

	return nil
}*/


func (h *handler) Auth(ctx context.Context, req *pb.User, resp *pb.Token) error {
	u, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	fmt.Println(&u)
	fmt.Println("db: %s", u.Password)
	fmt.Println("req: %s",req.Password)
	// 密码验证
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return err
	}

	t, err := h.tokenService.Encode(u)
	if err != nil {
		return err
	}
	resp.Token = t
	return nil
}

func (h *handler) ValidateToken(ctx context.Context, req *pb.Token, resp *pb.Token) error {
	// decode token
	claims, err := h.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	resp.Valid = true
	return nil
}
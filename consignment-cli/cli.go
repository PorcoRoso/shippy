package main

import (
	"io/ioutil"
	"encoding/json"
	"github.com/micro/go-micro/cmd"
	"log"
	"errors"
	microclient "github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	pb "github.com/porcorosso/shippy/consignment-service/proto/consignment"
	vpb "github.com/porcorosso/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro/metadata"
)

const (
	ADDRESS           = "localhost:50051"
	DEFAULT_INFO_FILE = "./consignment.json"
	DEFAULT_INFO_FILE1 = "./vessel.json"
)

// 读取 consignment.json 中记录的货物信息
func parseFileConsignment(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment *pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("consignment.json file content error.")
	}
	return consignment, nil
}

func parseFileVessel(fileName string) (*vpb.Vessel, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var vessel *vpb.Vessel
	err = json.Unmarshal(data, &vessel)
	if err != nil {
		return nil, errors.New("vessel.json file content error.")
	}
	return vessel, nil
}

func main() {

	cmd.Init()
	// 创建微服务的客户端，简化了手动 Dial 连接服务端的步骤
	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)
	vClient := vpb.NewVesselServiceClient("go.micro.srv.vessel", microclient.DefaultClient)
	// 解析货物信息
	consignment, err := parseFileConsignment(DEFAULT_INFO_FILE)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 解析轮船信息
	vessel, err := parseFileVessel(DEFAULT_INFO_FILE1)
	if err != nil {
		log.Fatalf("parse info file error: %v", err)
	}

	// 调用RPC,将货物存储到我们自己的仓库
	vResp, err := vClient.Create(context.Background(), vessel)
	if err != nil {
		log.Fatalf("create vessel error: %v", err)
	}
	// 货轮是否创建成功
	log.Printf("created: %t", vResp.Created)

	// 创建带有用户token的context，consignment-service服务端将从中取出token，解密取出用户身份
	// 用户token
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7ImlkIjoiNjNhYjRhNjgtYTBiOC00NWUzLTg5MjAtMjAyNzIzOTM2MzBjIiwibmFtZSI6ImV2YW4iLCJjb21wYW55IjoiYmJjIiwiZW1haWwiOiJldmFuQGJiYyIsInBhc3N3b3JkIjoiJDJhJDEwJEViWlhBOFZWdDc4d1MuekJaVXBsdnVQUzhacFpSd3dDcGVkMm9XQUxTLzY0ZVNwekVjemFpIn0sImV4cCI6MTU0OTE4MzY2OCwiaXNzIjoiZ28ubWljcm8uc3J2LnVzZXIifQ.VPGMSnJj-jOEtJsiCpdHnx3EPkrsmoA0FPcqUt2nlmI"
	tokenContext := metadata.NewContext(context.Background(), map[string]string{
		"Token":token,
	})

	// 调用RPC
	// 将货物存储到我们自己的仓库
	resp, err := client.CreateConsignment(tokenContext, consignment)
	if err != nil {
		log.Fatalf("create consignment error: %v", err)
	}

	// 新货物是否托运成功
	log.Printf("created: %t", resp.Created)

	// 列出目前所有托运的货物
	resp, err = client.GetConsignments(tokenContext, &pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignments: %v", err)
	}
	for _, c := range resp.Consignments {
		log.Printf("%+v", c)
	}

}
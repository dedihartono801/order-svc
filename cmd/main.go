package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dedihartono801/order-svc/pkg/client"
	"github.com/dedihartono801/order-svc/pkg/config"
	"github.com/dedihartono801/order-svc/pkg/db"
	"github.com/dedihartono801/order-svc/pkg/service"
	pb "github.com/dedihartono801/protobuf/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	productSvc := client.InitProductServiceClient(c.ProductSvcUrl)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Order Svc on", c.Port)

	s := service.Server{
		H:          h,
		ProductSvc: productSvc,
	}

	opts := []grpc.ServerOption{}
	tls := true

	if tls {
		certFile := "ssl/order-svc/server.crt"
		kefFile := "ssl/order-svc/server.pem"

		creds, err := credentials.NewServerTLSFromFile(certFile, kefFile)

		if err != nil {
			log.Fatalf("Failed loading certificates: %v\n", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterOrderServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

package client

import (
	"context"
	"fmt"
	"log"

	pb "github.com/dedihartono801/protobuf/product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ProductServiceClient struct {
	Client pb.ProductServiceClient
}

func InitProductServiceClient(url string) ProductServiceClient {
	opts := []grpc.DialOption{}
	tls := true

	if tls {
		certFile := "../ssl/product-svc/ca.crt"

		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("Error while loading CA trust certificates: %v\n", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	cc, err := grpc.Dial(url, opts...)

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := ProductServiceClient{
		Client: pb.NewProductServiceClient(cc),
	}

	return c
}

func (c *ProductServiceClient) FindOne(productId int64) (*pb.FindOneResponse, error) {
	req := &pb.FindOneRequest{
		Id: productId,
	}

	return c.Client.FindOne(context.Background(), req)
}

func (c *ProductServiceClient) DecreaseStock(productId int64, orderId int64) (*pb.DecreaseStockResponse, error) {
	req := &pb.DecreaseStockRequest{
		Id:      productId,
		OrderId: orderId,
	}

	return c.Client.DecreaseStock(context.Background(), req)
}

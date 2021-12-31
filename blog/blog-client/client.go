package main

import (
	"context"
	"fmt"
	"log"

	"github.com/liridonrama/grpc-go-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("localhost:6543", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "",
			AuthorId: primitive.NewObjectID().Hex(),
			Title:    "Some title",
			Content:  "Some content",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Panicln("Error while getting response from server", err)
	}

	fmt.Println("Blog has been created:", res.GetBlog())

	rBReq := &blogpb.ReadBlogRequest{
		BlogId: "",
	}

	rBRes, err := c.ReadBlog(context.Background(), rBReq)
	if err != nil {
		log.Panicln("Error while getting response from server", err)
	}

	fmt.Println("Blog has been retrieved:", rBRes.GetBlog())
}

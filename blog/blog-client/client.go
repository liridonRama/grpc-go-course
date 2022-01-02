package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
		BlogId: "61cf25475edd00e9eac3ef95",
	}

	rBRes, err := c.ReadBlog(context.Background(), rBReq)
	if err != nil {
		log.Println("Error while getting response from server", err)
	} else {
		fmt.Println("Blog has been retrieved:", rBRes.GetBlog())
	}

	uBReq := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:      "61cf25475edd00e9eac3ef95",
			Content: "ueli, hans mauerer",
		},
	}

	uBRes, err := c.UpdateBlog(context.Background(), uBReq)
	if err != nil {
		log.Println("Error while getting response from server", err)
	} else {
		fmt.Println("Blog has been updated:", uBRes.GetBlog())
	}

	dBRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: "61d00f562023fb557458d377",
	})
	if err != nil {
		log.Println("Error while getting response from server", err)
	} else {
		fmt.Println("Blog has been deleted:", dBRes.GetBlogId())
	}

	// list blogs
	rStream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Panicln("ListBlog — Error while trying to send request", err)
	}

	for {
		blogM, err := rStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panicln("Error while trying to read stream", err)
		}

		fmt.Println("blog reveiced:", blogM.GetBlog())

		time.Sleep(time.Second)
	}

	fmt.Println("All Blogs retrieved")
}

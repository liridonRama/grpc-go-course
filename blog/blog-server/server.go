package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/liridonrama/grpc-go-course/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	log.Println("Create Blog Request RPC Call")

	blog := req.GetBlog()

	authorId, err := primitive.ObjectIDFromHex(blog.GetAuthorId())
	if err != nil {
		return nil, err
	}

	data := BlogItem{
		ID:       primitive.NewObjectID(),
		AuthorID: authorId,
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	res := &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID.Hex(),
			Title:    data.Title,
			Content:  data.Content,
		},
	}

	return res, nil
}
func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	log.Println("Read Blog Request RPC Call")

	id := req.GetBlogId()

	blogId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "the provided id: %v is not a objectid string", id)
	}

	result := collection.FindOne(ctx, primitive.M{
		"_id": blogId,
	})

	if result.Err() == mongo.ErrNoDocuments {
		return nil, status.Errorf(codes.NotFound, "Blog with the id: %v not found", id)
	}

	var blog blogpb.Blog

	err = result.Decode(&blog)
	if err != nil {
		return nil, err
	}

	res := &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blog.Id,
			AuthorId: blog.AuthorId,
			Title:    blog.Title,
			Content:  blog.Content,
		},
	}

	return res, nil
}

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID primitive.ObjectID `bson:"authorId,omitempty"`
	Content  string             `bson:"content,omitempty"`
	Title    string             `bson:"title,omitempty"`
}

var collection *mongo.Collection

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Connecting to mongodb")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Panicln("Error while trying to connect to mongo")
	}

	collection = client.Database("grpc-go-course").Collection("blog")

	log.Println("Blog Service Started")
	mux, err := net.Listen("tcp", ":6543")
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	s := grpc.NewServer()
	defer s.Stop()

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		err = s.Serve(mux)
		if err != nil {
			log.Fatal("failed to listen", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Println("Stopping the server")
	s.Stop()

	log.Println("close mongodb connection")
	client.Disconnect(context.Background())

	log.Println("Closing the listener")
	mux.Close()
}

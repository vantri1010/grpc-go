package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"grpc-go/blog/blogpb"
)

func main() {
	// Define CLI flags
	operation := flag.String("operation", "", "CRUD operation: create, read, update, delete")
	blogID := flag.String("id", "", "ID of the blog (required for read, update, delete)")
	authorID := flag.String("author", "", "Author ID of the blog (required for create/update)")
	title := flag.String("title", "", "Title of the blog (required for create/update)")
	content := flag.String("content", "", "Content of the blog (required for create/update)")
	flag.Parse()

	if *operation == "" {
		fmt.Println("Error: operation flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Connect to gRPC server
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	switch *operation {
	case "create":
		if *authorID == "" || *title == "" || *content == "" {
			log.Fatalf("Error: author, title, and content are required for create operation")
		}
		blog := &blogpb.Blog{
			AuthorId: *authorID,
			Title:    *title,
			Content:  *content,
		}
		createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
		if err != nil {
			log.Fatalf("Unexpected error: %v", err)
		}
		fmt.Printf("Blog has been created: %v\n", createBlogRes)

	case "read":
		if *blogID == "" {
			log.Fatalf("Error: id is required for read operation")
		}
		readBlogRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: *blogID})
		if err != nil {
			log.Fatalf("Error while reading the blog: %v", err)
		}
		fmt.Printf("Blog found: %v\n", readBlogRes)

	case "update":
		if *blogID == "" || *authorID == "" || *title == "" || *content == "" {
			log.Fatalf("Error: id, author, title, and content are required for update operation")
		}
		blog := &blogpb.Blog{
			Id:       *blogID,
			AuthorId: *authorID,
			Title:    *title,
			Content:  *content,
		}
		updateBlogRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: blog})
		if err != nil {
			log.Fatalf("Error while updating the blog: %v", err)
		}
		fmt.Printf("Blog has been updated: %v\n", updateBlogRes)

	case "delete":
		if *blogID == "" {
			log.Fatalf("Error: id is required for delete operation")
		}
		deleteBlogRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: *blogID})
		if err != nil {
			log.Fatalf("Error while deleting the blog: %v", err)
		}
		fmt.Printf("Blog has been deleted: %v\n", deleteBlogRes)

	default:
		log.Fatalf("Unknown operation: %v", *operation)
	}
}

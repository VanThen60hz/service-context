package main

import (
	"context"
	"log"

	sctx "github.com/VanThen60hz/service-context"
	s3c "github.com/VanThen60hz/service-context/component/s3c"
)

func main() {
	serviceCtx := sctx.NewServiceContext(
		sctx.WithName("my-service"),
		sctx.WithComponent(s3c.NewS3Component("s3")),
	)

	if err := serviceCtx.Load(); err != nil {
		log.Fatal(err)
	}

	// Use the S3 component
	s3Comp := serviceCtx.MustGet("s3").(*s3c.S3Component)
	url, err := s3Comp.Upload(context.TODO(), "", "images")
	if err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}

	log.Printf("File uploaded to: %s", url)

	_ = serviceCtx.Stop()
}

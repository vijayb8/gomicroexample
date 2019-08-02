package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/vijayb8/gomicroexample/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consigment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consigment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consigment, nil
}

type Service struct {
	repo Repository
}

func (s *Service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consigment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		Created:     true,
		Consignment: consigment,
	}, nil
}

func main() {
	repo := Repository{}

	//setup gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Cannot connect to gRPC server")
	}

	s := grpc.NewServer()

	pb.RegisterConsignmentServiceServer(s, &Service{repo})

	// register reflection
	reflection.Register(s)

	log.Println("server running")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

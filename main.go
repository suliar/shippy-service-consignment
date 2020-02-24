// shippy-service-consignment

package main

import (
	"context"
	pb "github.com/suliar/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

const port  = ":50051"

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later

type Repository struct {
	mu sync.RWMutex
	consignments []*pb.Consignment
}

// Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()

	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()

	return consignment, nil
}

// GetAll gets all consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you better idea
type service struct {
	repo repository
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()

	return &pb.Response{Consignments:consignments}, nil
}

// CreateConsignment - we created just one method on our service
// which is just one method, which takes a context and a request as
// an argument, these are handled by the grpc server


func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	//Return matching the 'Response' message we created in our
	//protobuf
	return &pb.Response{ Created:true, Consignment:consignment}, nil
}

func main() {
	repo := &Repository{}

	//Set up grpc server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen to: %v", err)
	}
	s := grpc.NewServer()


	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition

	pb.RegisterShippingServiceServer(s, &service{repo})

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
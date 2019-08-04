package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "../pb"
)

const (
	port = ":50051"
)

type fbService struct {
	listOfFb map[string]pb.Feedback
}

//implemen all API of PassengerFeedback service
func (fs *fbService) AddFb(ctx context.Context, req *pb.Feedback) (*pb.AddFbkResponse, error) {
	log.Println("Got request add feedback")
	fs.listOfFb[req.BookingCode] = *req
	// log.Println(fs.listOfFb)
	return &pb.AddFbkResponse{BookingCode: req.BookingCode, PassengerId: req.PassengerId}, status.Errorf(codes.OK, "Add successful")
}

func (fs *fbService) GetFbByPasId(ctx context.Context, req *pb.PassengerId) (*pb.ListOfFeedback, error) {
	log.Println("Got request get feedback by passenger id")
	var listFb []*pb.Feedback
	for _, fb := range fs.listOfFb {
		if fb.PassengerId == req.PassengerId {
			listFb = append(listFb, &pb.Feedback{BookingCode: fb.BookingCode,
				PassengerId: fb.PassengerId,
				Feedback:    fb.Feedback})
		}
	}
	if len(listFb) == 0 {
		return nil, status.Errorf(codes.NotFound, "Cant find any feedback by this passenger ID")
	}

	return &pb.ListOfFeedback{ListOfFeedback: listFb}, status.Errorf(codes.OK, "Get successful")
}

func (fs *fbService) GetFbByBookingCode(ctx context.Context, req *pb.BookingCode) (*pb.Feedback, error) {
	log.Println("Got request get feedback by booking code")
	if value, ok := fs.listOfFb[req.BookingCode]; ok == true {
		return &value, status.Errorf(codes.OK, "Get successful")
	}

	return nil, status.Errorf(codes.NotFound, "Cant find any feedback by this booking code")
}

func (fs *fbService) DeleteByPasId(ctx context.Context, req *pb.PassengerId) (*pb.DeleteByPasIdResponse, error) {
	log.Println("Got request delete feedback by passenger id")
	var cnt int32
	for bc, fb := range fs.listOfFb {
		if fb.PassengerId == req.PassengerId {
			delete(fs.listOfFb, bc)
			cnt++
		}
	}

	if cnt == 0 {
		return nil, status.Errorf(codes.NotFound, "Cant find any feedback by this passenger ID")
	}

	return &pb.DeleteByPasIdResponse{NumberOfFbDeleted: cnt}, status.Errorf(codes.OK, "Delete successful")

}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterPassengerFeedbackServer(server, &fbService{listOfFb: make(map[string]pb.Feedback)})
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

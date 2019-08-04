package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "../pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	adress       = "localhost:50051"
	bookingCode1 = "123asd"
	passengerId1 = 1
	feedback1    = "vui tinh"
	bookingCode2 = "123asd1"
	passengerId2 = 1
	feedback2    = "1 sao"
)

func addFeedback(fb *pb.Feedback, client pb.PassengerFeedbackClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.AddFb(ctx, fb)

	if verifyThenPrintError(err) == false {
		return err
	}

	fmt.Println("User ", response.PassengerId, " sent feedback succesful for booking code: ", response.BookingCode)
	return nil
}

func printFbByPasId(pid int32, client pb.PassengerFeedbackClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.GetFbByPasId(ctx, &pb.PassengerId{PassengerId: pid})

	if verifyThenPrintError(err) == false {
		return
	}

	fmt.Printf("List feedback from Passenger ID %d:\n", pid)
	for _, fb := range response.ListOfFeedback {
		printFb(fb)
	}
}

func printFb(fb *pb.Feedback) {
	fmt.Println("####################")
	fmt.Println("Booking code: ", fb.BookingCode)
	fmt.Println("Feedback: ", fb.Feedback)
	fmt.Println("####################")
}

func printFbByBookingCode(bc string, client pb.PassengerFeedbackClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.GetFbByBookingCode(ctx, &pb.BookingCode{BookingCode: bc})

	if verifyThenPrintError(err) == false {
		return
	}

	fmt.Printf("List feedback from booking code %s:\n", bc)
	printFb(response)
}

func deleteFbByPasId(pid int32, client pb.PassengerFeedbackClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.DeleteByPasId(ctx, &pb.PassengerId{PassengerId: pid})

	if verifyThenPrintError(err) == false {
		return
	}

	fmt.Println("Deleted ", response.NumberOfFbDeleted, " feedback(s) from database")
}

func verifyThenPrintError(err error) bool {
	if err != nil {
		grpcErr, ok := status.FromError(err)
		if ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				fmt.Println("Not Found: ", grpcErr.Message())
			case codes.OK:
				fmt.Println("OK: ", grpcErr.Message())
				return true
			default:
				fmt.Println("Unexpected error: ", grpcErr.Code())
			}
		} else {
			fmt.Println("Failed to call server", err)
		}
		return false
	}
	return true
}

func main() {
	connection, err := grpc.Dial(adress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connection.Close()

	client := pb.NewPassengerFeedbackClient(connection)

	fb1 := pb.Feedback{BookingCode: bookingCode1,
		PassengerId: passengerId1,
		Feedback:    feedback1}
	fb2 := pb.Feedback{BookingCode: bookingCode2,
		PassengerId: passengerId2,
		Feedback:    feedback2}

	addFeedback(&fb1, client)
	addFeedback(&fb2, client)
	printFbByPasId(passengerId1, client)
	printFbByBookingCode(bookingCode1, client)
	deleteFbByPasId(passengerId1, client)
	printFbByPasId(1, client)
}

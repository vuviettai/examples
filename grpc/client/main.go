package main

import (
	"context"
	pb "grpc-client/proto"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	serverUrl := os.Getenv("GRPC_SERVER_URL")
	conn, err := grpc.NewClient(serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRouteGuideClient(conn)
	runRouteChat(c)
	consensusClient := pb.NewConsensusApiClient((conn))
	initTransactionStream(consensusClient)
}
func initTransactionStream(client pb.ConsensusApiClient) {
	txs := []*pb.ExternalTransaction{
		{Namespace: "EvmOs", TxBytes: []byte{1}},
		{Namespace: "EvmOs", TxBytes: []byte{2}},
		{Namespace: "EvmOs", TxBytes: []byte{3}},
		{Namespace: "EvmOs", TxBytes: []byte{4}},
		{Namespace: "EvmOs", TxBytes: []byte{5}},
		{Namespace: "EvmOs", TxBytes: []byte{6}},
		{Namespace: "EvmOs", TxBytes: []byte{7}},
		{Namespace: "EvmOs", TxBytes: []byte{8}},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.InitTransaction(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.InitTransaction failed: %v", err)
			}
			txs := in.Transactions
			log.Printf("Got commited transactions %s", txs)
		}
	}()
	for _, tx := range txs {
		if err := stream.Send(tx); err != nil {
			log.Fatalf("client.Sendtx: stream.Send(%v) failed: %v", tx, err)
		}
	}
	stream.CloseSend()
	<-waitc
}
func runRouteChat(client pb.RouteGuideClient) {
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", note, err)
		}
	}
	stream.CloseSend()
	<-waitc
}

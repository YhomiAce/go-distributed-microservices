package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest)(*logs.LogResponse, error) {
	input := req.GetLogEntry()
	entry := data.LogEntry {
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(entry)
	if err != nil {
		return &logs.LogResponse{Result: "Failed to insert using GRPC"}, err
	}
	res := &logs.LogResponse{Result: "Logged using GRPC"}
	return res, nil
}

func (app *Application) grpcListen() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s",grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for GRPC %v\n", err)
	}
	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{
		Models: app.Models,
	})
	log.Printf("GRPC server started on PORT: %s", grpcPort)

	 err = s.Serve(listener);
	 if err != nil {
		log.Fatalf("Failed to Listen for gRPC %v", err)
	 }
}
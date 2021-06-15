package main

import (
	"database/sql"
	"net"
	"flag"

	"google.golang.org/grpc"
	"github.com/AshutoshKumar/entain/sports/db"
	"github.com/AshutoshKumar/entain/sports/service"
	"github.com/AshutoshKumar/entain/sports/proto/sports"
	log "github.com/sirupsen/logrus"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9000", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9000")
	if err != nil {
		return err
	}

	sportsDB, err := sql.Open("sqlite3", "./db/sports.db")
	if err != nil {
		return err
	}

	sportsRepo := db.NewSportsRepo(sportsDB)
	if err := sportsRepo.Init(); err != nil {
		return err
	}

	grpcServer := grpc.NewServer()


	sports.RegisterSportsServer(grpcServer, service.NewSportsService(sportsRepo))

	log.Infof("gRPC server listening on: %s", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"github.com/odysseia-greek/agora/eupalinos/stomion"
	"github.com/odysseia-greek/agora/plato/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const standardPort = ":50060"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	logging.System(`
   ___  __ __  ____   ____  _      ____  ____    ___   _____
  /  _]|  |  ||    \ /    || |    |    ||    \  /   \ / ___/
 /  [_ |  |  ||  o  )  o  || |     |  | |  _  ||     (   \_ 
|    _]|  |  ||   _/|     || |___  |  | |  |  ||  O  |\__  |
|   [_ |  :  ||  |  |  _  ||     | |  | |  |  ||     |/  \ |
|     ||     ||  |  |  |  ||     | |  | |  |  ||     |\    |
|_____| \__,_||__|  |__|__||_____||____||__|__| \___/  \___|
	`)

	logging.System("ἀρχιτέκτων δὲ τοῦ ὀρύγματος τούτου ἐγένετο Μεγαρεὺς Εὐπαλῖνος Ναυστρόφου")
	logging.System("The designer of this work was Eupalinus son of Naustrophus, a Megarian")
	logging.System("Starting up...")

	logging.System("starting up.....")
	logging.System("starting up and getting env variables")

	config, err := stomion.CreateNewConfig()
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	// Check if HTTPS mode is enabled
	if config.TLSConfig != nil {
		logging.System("starting up in HTTPS mode")
		server = grpc.NewServer(grpc.Creds(credentials.NewTLS(config.TLSConfig)))
	} else {
		logging.System("starting up in HTTP mode")
		server = grpc.NewServer()
	}

	config.LoadStateFromDisk()
	config.StartAutoSave()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logging.System(fmt.Sprintf("system send a: %s stopping service", sig.String()))
		config.SaveStateToDisk()
		os.Exit(0)
	}()

	pb.RegisterEupalinosServer(server, config)

	if config.Streaming {
		config.StartBroadcasting()
	}

	config.PeriodStatsPrint()

	logging.System(fmt.Sprintf("Server listening on %s", port))
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

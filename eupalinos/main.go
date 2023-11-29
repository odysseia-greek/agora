package main

import (
	"github.com/odysseia-greek/agora/eupalinos/app"
	"github.com/odysseia-greek/agora/eupalinos/config"
	pb "github.com/odysseia-greek/agora/eupalinos/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const standardPort = ":50051"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = standardPort
	}

	log.Print(`
   ___  __ __  ____   ____  _      ____  ____    ___   _____
  /  _]|  |  ||    \ /    || |    |    ||    \  /   \ / ___/
 /  [_ |  |  ||  o  )  o  || |     |  | |  _  ||     (   \_ 
|    _]|  |  ||   _/|     || |___  |  | |  |  ||  O  |\__  |
|   [_ |  :  ||  |  |  _  ||     | |  | |  |  ||     |/  \ |
|     ||     ||  |  |  |  ||     | |  | |  |  ||     |\    |
|_____| \__,_||__|  |__|__||_____||____||__|__| \___/  \___|
	`)

	log.Print("ἀρχιτέκτων δὲ τοῦ ὀρύγματος τούτου ἐγένετο Μεγαρεὺς Εὐπαλῖνος Ναυστρόφου")
	log.Print("The designer of this work was Eupalinus son of Naustrophus, a Megarian")
	log.Print("Starting up...")

	env := os.Getenv("ENV")

	// Create the configuration based on the environment
	eupalinosConf, err := config.CreateNewConfig(env)
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var server *grpc.Server

	// Check if HTTPS mode is enabled
	if eupalinosConf.TLSConfig != nil {
		server = grpc.NewServer(grpc.Creds(credentials.NewTLS(eupalinosConf.TLSConfig)))
	} else {
		server = grpc.NewServer()
	}

	queueServer := &app.EupalinosHandler{
		DiexodosMap: make([]*app.Diexodos, 0),
		Config:      eupalinosConf,
	}

	queueServer.LoadStateFromDisk()
	queueServer.StartAutoSave()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("system send a: %s stopping service", sig.String())
		queueServer.SaveStateToDisk()
		os.Exit(0)
	}()

	pb.RegisterEupalinosServer(server, queueServer)

	if queueServer.Config.Streaming {
		queueServer.StartBroadcasting()
	}

	log.Printf("Server listening on %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"github.com/odysseia-greek/agora/plato/logging"
	"github.com/odysseia-greek/agora/theofrastos/futikon"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//https://patorjk.com/software/taag/#p=display&f=Crawford2&t=theofrastos
	logging.System(`
 ______  __ __    ___   ___   _____  ____    ____  _____ ______   ___   _____
|      ||  |  |  /  _] /   \ |     ||    \  /    |/ ___/|      | /   \ / ___/
|      ||  |  | /  [_ |     ||   __||  D  )|  o  (   \_ |      ||     (   \_ 
|_|  |_||  _  ||    _]|  O  ||  |_  |    / |     |\__  ||_|  |_||  O  |\__  |
  |  |  |  |  ||   [_ |     ||   _] |    \ |  _  |/  \ |  |  |  |     |/  \ |
  |  |  |  |  ||     ||     ||  |   |  .  \|  |  |\    |  |  |  |     |\    |
  |__|  |__|__||_____| \___/ |__|   |__|\_||__|__| \___|  |__|   \___/  \___|
`)
	logging.System("\"Εἰ μὲν ἀμαθὴς εἶ, φρονίμως ποιεῖς, εἰ δὲ πεπαίδευσαι, ἀφρόνως.\"")
	logging.System("\"If you are an ignorant man, you are acting wisely; but if you have had any education, you are behaving like a fool.\"")
	logging.System("starting up and getting environment variables...")

	handler, err := futikon.CreateNewConfig()
	if err != nil {
		logging.Error(err.Error())
		log.Fatal("death has found me")
	}

	if err := handler.WatchConfigMapChanges(ctx); err != nil && ctx.Err() == nil {
		logging.Errorf("failed to watch ConfigMap: %v", err)
	}
}

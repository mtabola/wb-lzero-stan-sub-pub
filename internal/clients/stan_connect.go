package clients

import (
	"log"
	"stan/config"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func StanNewConnection(cfg *config.Config, name string) (*nats.Conn, stan.Conn) {
	opts := []nats.Option{nats.Name(name)}

	nc, err := nats.Connect(cfg.Addr, opts...)
	if err != nil {
		log.Fatal(err)
	}
	sc, err := stan.Connect(cfg.Cluster, cfg.User, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			nc.Close()
			log.Fatalf("Connection lost, reason: %v", reason)
		}))

	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, cfg.Addr)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", cfg.Addr, cfg.Cluster, cfg.User)
	return nc, sc
}

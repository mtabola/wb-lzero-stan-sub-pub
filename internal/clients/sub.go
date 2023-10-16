package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"stan/config"
	"stan/internal/models"

	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
)

func sentData(json_b []byte) {
	log.Println("Sending...")
	resp, err := http.Post("http://localhost:8100/orders", "application/json", bytes.NewBuffer(json_b))
	if err != nil {
		log.Panic("Connection error: ", err)
	}
	var res = models.JSONRequest{}
	json.NewDecoder(resp.Body).Decode(&res)

	log.Printf("Response: %v", res.Data)
}

func NewSubscriberListener(cfg *config.Config) {

	nc, sc := StanNewConnection(cfg, "STAN Subscrber")

	defer nc.Close()

	defer sc.Close()

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)

	notSentedData := make(map[int][]byte)
	countForNotSented := 0

	mcb := func(msg *stan.Msg) {
		req := models.JSONResponse{
			Data:      string(msg.Data),
			Sequence:  msg.Sequence,
			Timestamp: msg.Timestamp,
		}
		json_d, err := json.Marshal(req)

		if err != nil {
			log.Panic("Marshal error: ", err)
		}

		_, err = exec.Command("lsof", "-nP", "-i4TCP:8100").CombinedOutput()
		if err != nil {
			notSentedData[countForNotSented] = json_d
			countForNotSented++
			log.Println("Wait of senting")
		} else {
			if len(notSentedData) != 0 {
				for k, v := range notSentedData {
					sentData(v)
					delete(notSentedData, k)
				}
			}
			sentData(json_d)
		}
	}

	sub, err := sc.QueueSubscribe(cfg.Channel, "", mcb, startOpt, stan.DurableName(""))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s]\n", cfg.Channel, cfg.User)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			sub.Unsubscribe()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}

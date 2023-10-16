package clients

import (
	"log"
	"os"
	"stan/config"
)

func NewPublisher(cfg *config.Config, filePath string) {
	nc, sc := StanNewConnection(cfg, "STAN Publisher")

	defer nc.Close()

	defer sc.Close()

	bfile, err := os.ReadFile(string(filePath))

	if err != nil {
		log.Fatal("File read failed")
	}
	err = sc.Publish(cfg.Channel, bfile)
	if err != nil {
		log.Fatalf("Error during publish: %v\n", err)
	}
	log.Printf("Published on [%s] file: '%s'\n", cfg.Channel, filePath)
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	paspoid "github.com/paspoid/server-api"
)

func main() {
	_ = godotenv.Load()

	baseUrl := os.Getenv("PASPOID_BASE_URL")

	apiKey := os.Getenv("PASPOID_API_KEY")
	apiSecret := os.Getenv("PASPOID_API_SECRET")
	servicePublicId := os.Getenv("PASPOID_SERVICE_PUBLIC_ID")

	transactionType := os.Getenv("PASPOID_TRANSACTION_TYPE")

	c := paspoid.NewClient(
		baseUrl,
		apiKey,
		apiSecret,
	)

	ctx := context.Background()

	getKeyResp, err := c.GetKey(ctx, servicePublicId, transactionType)
	if err != nil {
		log.Fatalf("failed to get key: %v", err)
	}

	fmt.Println("key:", getKeyResp.Key)
	fmt.Println("validation_window:", getKeyResp.ValidationWindow)

	validateResp, err := c.Validate(ctx, getKeyResp.Key)
	if err != nil {
		log.Fatalf("failed to validate: %v", err)
	}

	fmt.Println("status:", validateResp.Status)
	if validateResp.DataType != nil {
		fmt.Println("data_type:", *validateResp.DataType)
	}
	if validateResp.DataValue != nil {
		fmt.Println("data_value:", *validateResp.DataValue)
	}
	fmt.Println("phone_data:", string(validateResp.PhoneData))
	fmt.Println("device_data:", string(validateResp.DeviceData))
}

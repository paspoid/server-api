package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	paspoid "github.com/paspoid/server-api"
)

func main() {
	if err := godotenv.Overload(); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

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

	validationWindow, err := time.ParseDuration(getKeyResp.ValidationWindow)
	if err != nil {
		log.Fatalf(
			"invalid validation_window %q: %v",
			getKeyResp.ValidationWindow,
			err,
		)
	}
	if validationWindow <= 0 {
		log.Fatalf("validation_window must be positive: %s", validationWindow)
	}

	pollInterval := validationWindow / 10
	if pollInterval < time.Second {
		pollInterval = time.Second
	}
	if pollInterval > 5*time.Second {
		pollInterval = 5 * time.Second
	}

	fmt.Println("poll_interval:", pollInterval)

	pollCtx, cancel := context.WithTimeout(ctx, validationWindow)
	defer cancel()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-pollCtx.Done():
			log.Fatalf(
				"validation window expired after %s",
				validationWindow,
			)

		case <-ticker.C:
			validateResp, err := c.Validate(pollCtx, getKeyResp.Key)
			if err != nil {
				log.Printf("validate request failed, retrying: %v", err)
				continue
			}

			fmt.Println("status:", validateResp.Status)

			switch validateResp.Status {
			case "incomplete":
				continue

			case "success":
				printValidationResult(validateResp)
				return

			case "failed":
				log.Fatal("validation failed")

			default:
				log.Printf(
					"unknown validation status %q, retrying",
					validateResp.Status,
				)
			}
		}
	}
}

func printValidationResult(validateResp *paspoid.ValidateResponse) {
	if validateResp.DataType != nil {
		fmt.Println("data_type:", *validateResp.DataType)
	}
	if validateResp.DataValue != nil {
		fmt.Println("data_value:", *validateResp.DataValue)
	}
	fmt.Println("phone_data:", string(validateResp.PhoneData))
	fmt.Println("device_data:", string(validateResp.DeviceData))
}

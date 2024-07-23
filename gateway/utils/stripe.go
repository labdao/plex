package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/billing/meterevent"
)

func setupStripeClient() error {
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey == "" {
		return fmt.Errorf("STRIPE_SECRET_KEY environment variable not set")
	}
	stripe.Key = apiKey
	return nil
}

func RecordUsage(stripeCustomerID string, usage int64) error {
	err := setupStripeClient()
	if err != nil {
		return fmt.Errorf("failed to set up Stripe client: %v", err)
	}

	params := &stripe.BillingMeterEventParams{
		EventName: stripe.String("compute_units"),
		Payload: map[string]string{
			"value":              strconv.FormatInt(usage, 10),
			"stripe_customer_id": stripeCustomerID,
		},
		Identifier: stripe.String(fmt.Sprintf("usage-%d", time.Now().Unix())),
	}
	_, err = meterevent.New(params)
	if err != nil {
		return fmt.Errorf("failed to record usage: %v", err)
	}
	return nil
}

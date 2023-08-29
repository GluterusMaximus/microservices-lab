package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"

	kafkautils "github.com/GluterusMaximus/ci/services/notify/kafka-utils"
	"github.com/GluterusMaximus/ci/services/notify/payloads"
	"github.com/sirupsen/logrus"
)

var (
	kafkaHost     string
	consumerTopic string
	brevoKey      string
	logger        *logrus.Logger
)

func init() {
	kafkaHost = os.Getenv("KAFKA_HOST")
	consumerTopic = os.Getenv("CONSUMER_TOPIC")
	brevoKey = os.Getenv("BREVO_API_KEY")
	logger = logrus.New()
}

type Order struct {
	ID             int
	RequestorEmail string
	Status         string
}

func main() {
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic occurred: %v", r)
			debug.PrintStack()
		}
	}()

	logger.Info("Application started")
	kafkautils.NewKafkaConsumer([]string{kafkaHost}, consumerTopic, "email-service", brevoKey, logger, HandleKafkaMessage)

	select {}
}

func HandleKafkaMessage(message []byte) {
	order, err := ParseKafkaMessage(message)
	if err != nil {
		logger.WithField("message", message).Errorf("failed to parse Kafka message: %v", err)
		return
	}

	err = SendEmailViaBrevo(order, brevoKey) // key is delegated to func cause it also need to be tested separated
	if err != nil {
		logger.WithFields(logrus.Fields{
			"orderID": order.ID,
			"email":   order.RequestorEmail,
		}).Errorf("failed to send email notification: %v", err)
	}
}

func ParseKafkaMessage(message []byte) (Order, error) {
	var order Order

	err := json.Unmarshal(message, &order)
	if err != nil {
		return Order{}, fmt.Errorf("failed to parse Kafka message: %w", err)
	}

	return order, nil
}

func SendHttpRequest(client *http.Client, method, url string, headers map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// SendEmailViaBrevo sends an email using the Brevo API.
// It takes an order, an API key.
// It returns an error if sending the email fails.
func SendEmailViaBrevo(order Order, apiKey string) error {

	emailContent := payloads.EmailContent{
		Sender: payloads.Sender{
			Name:  "Alex Newman",
			Email: "poctester24@gmail.com",
		},
		To: []payloads.Recipient{
			{
				Email: order.RequestorEmail,
				Name:  "Dear Recipient",
			},
		},
		Subject:     fmt.Sprintf("Hello %s", order.RequestorEmail),
		HtmlContent: fmt.Sprintf("<html><head></head><body><p>Your order %d changed status</p>Your order %d changed it status to %s.</p></body></html>", order.ID, order.ID, order.Status),
	}

	jsonData, err := json.Marshal(emailContent)
	if err != nil {
		return err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	headers := map[string]string{
		"accept":       "application/json",
		"api-key":      apiKey,
		"content-type": "application/json",
	}

	responseBody, err := SendHttpRequest(client, "POST", "https://api.brevo.com/v3/smtp/email", headers, jsonData)
	if err != nil {
		return err
	}
	logger.Infof("Response Body: %s", responseBody)

	logger.Info("Email sent successfully")

	return nil
}

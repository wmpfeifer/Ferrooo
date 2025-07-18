package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type ProcessorClient struct {
    defaultURL  string
    fallbackURL string
    client      *fasthttp.Client
}

type ProcessorRequest struct {
    CorrelationID uuid.UUID `json:"correlationId"`
    Amount        float64   `json:"amount"`
    RequestedAt   string    `json:"requestedAt"`
}

func NewProcessorClient(defaultURL, fallbackURL string) *ProcessorClient {
    return &ProcessorClient{
        defaultURL:  defaultURL,
        fallbackURL: fallbackURL,
        client: &fasthttp.Client{
            ReadTimeout:  time.Second * 10,
            WriteTimeout: time.Second * 10,
        },
    }
}

func (c *ProcessorClient) ProcessPayment(ctx context.Context, correlationID uuid.UUID, amount float64, requestedAt time.Time, useDefault bool) error {
    url := c.fallbackURL
    if useDefault {
        url = c.defaultURL
    }

    reqData := ProcessorRequest{
        CorrelationID: correlationID,
        Amount:        amount,
        RequestedAt:   requestedAt.Format(time.RFC3339),
    }

    jsonData, err := json.Marshal(reqData)
    if err != nil {
        return err
    }

    req := fasthttp.AcquireRequest()
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseRequest(req)
    defer fasthttp.ReleaseResponse(resp)

    req.SetRequestURI(url + "/payments")
    req.Header.SetMethod("POST")
    req.Header.SetContentType("application/json")
    req.SetBody(jsonData)

    err = c.client.DoTimeout(req, resp, time.Second*30)
    if err != nil {
        return err
    }

    if resp.StatusCode() >= 400 {
        return fmt.Errorf("processor returned status: %d", resp.StatusCode())
    }

    return nil
}
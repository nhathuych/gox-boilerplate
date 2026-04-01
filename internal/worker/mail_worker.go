package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/nhathuych/gox-boilerplate/internal/config"
)

func mailIdempotencyKey(messageID string) string {
	return fmt.Sprintf("idempotency:mail:%s", messageID)
}

// MailJob is the JSON payload consumed from the mail queue.
type MailJob struct {
	MessageID string `json:"message_id"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
}

// MailWorker consumes mail jobs from RabbitMQ and skips duplicates using Redis (SETNX on message_id).
type MailWorker struct {
	cfg *config.Config
	rdb *redis.Client
	log *zap.Logger
	ch  *amqp.Channel
}

func NewMailWorker(cfg *config.Config, rdb *redis.Client, log *zap.Logger) *MailWorker {
	return &MailWorker{cfg: cfg, rdb: rdb, log: log}
}

// Run blocks until the AMQP channel is closed or the connection fails.
func (w *MailWorker) Run(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	w.ch = ch
	if _, err := ch.QueueDeclare(
		w.cfg.Worker.MailQueueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	if err := ch.Qos(w.cfg.RabbitMQ.PrefetchCount, 0, false); err != nil {
		return err
	}
	msgs, err := ch.Consume(
		w.cfg.Worker.MailQueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for d := range msgs {
		w.handleDelivery(d)
	}
	return nil
}

func (w *MailWorker) handleDelivery(d amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var job MailJob
	if err := json.Unmarshal(d.Body, &job); err != nil {
		w.log.Warn("invalid mail job payload", zap.Error(err))
		_ = d.Nack(false, false)
		return
	}
	if job.MessageID == "" {
		w.log.Warn("missing message_id")
		_ = d.Nack(false, false)
		return
	}

	key := mailIdempotencyKey(job.MessageID)
	// Idempotency: only the first successful SETNX proceeds; duplicates ack without work.
	ok, err := w.rdb.SetNX(ctx, key, "1", 7*24*time.Hour).Result()
	if err != nil {
		w.log.Error("redis idempotency check", zap.Error(err))
		_ = d.Nack(true, true)
		return
	}
	if !ok {
		w.log.Info("skip duplicate mail job", zap.String("message_id", job.MessageID))
		_ = d.Ack(false)
		return
	}

	w.log.Info("process mail job",
		zap.String("message_id", job.MessageID),
		zap.String("to", job.To),
		zap.String("subject", job.Subject),
	)
	_ = d.Ack(false)
}

// Shutdown closes the consumer channel so Run returns.
func (w *MailWorker) Shutdown(ctx context.Context) error {
	if w.ch == nil {
		return nil
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = w.ch.Close()
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

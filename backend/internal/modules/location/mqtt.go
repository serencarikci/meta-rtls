package location

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTWorker struct {
	client mqtt.Client
	svc    *Service
	logger *slog.Logger
	topic  string
}

func NewMQTTWorker(broker, clientID, topic string, svc *Service, logger *slog.Logger) *MQTTWorker {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(3 * time.Second)

	w := &MQTTWorker{svc: svc, logger: logger, topic: topic}
	opts.SetDefaultPublishHandler(func(_ mqtt.Client, msg mqtt.Message) {
		w.onMessage(msg)
	})
	w.client = mqtt.NewClient(opts)
	return w
}

func (w *MQTTWorker) Start() {
	token := w.client.Connect()
	token.Wait()
	if token.Error() != nil {
		w.logger.Warn("mqtt connect failed (simulator can still run local)", "err", token.Error())
		return
	}
	sub := w.client.Subscribe(w.topic, 0, nil)
	sub.Wait()
	if sub.Error() != nil {
		w.logger.Warn("mqtt subscribe failed", "err", sub.Error())
		return
	}
	w.logger.Info("mqtt subscribed", "topic", w.topic)
}

func (w *MQTTWorker) Stop() {
	if w.client != nil && w.client.IsConnected() {
		w.client.Disconnect(250)
	}
}

func (w *MQTTWorker) Publish(event PositionEvent) {
	if w.client == nil || !w.client.IsConnected() {
		return
	}
	body, err := json.Marshal(event)
	if err != nil {
		return
	}
	topic := "rtls/" + event.TenantID + "/location"
	w.client.Publish(topic, 0, false, body)
}

func (w *MQTTWorker) onMessage(msg mqtt.Message) {
	var event PositionEvent
	if err := json.Unmarshal(msg.Payload(), &event); err != nil {
		w.logger.Warn("bad mqtt payload", "err", err)
		return
	}
	w.svc.HandleEvent(context.Background(), event)
}

func (w *MQTTWorker) Connected() bool {
	return w.client != nil && w.client.IsConnected()
}

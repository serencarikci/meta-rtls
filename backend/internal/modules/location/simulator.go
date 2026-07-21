package location

import (
	"context"
	"log/slog"
	"math"
	"sync"
	"time"
)

const (
	SIM_TICK    = 2 * time.Second
	SIM_FLOOR_W = 80.0
	SIM_FLOOR_H = 40.0
	DEMO_TENANT = "warehouse-s"
)

type Simulator struct {
	repo   *Repository
	svc    *Service
	mqtt   *MQTTWorker
	logger *slog.Logger

	mu       sync.Mutex
	running  bool
	stopCh   chan struct{}
	tenantID string
	floorID  string
	tags     []SimTag
}

func NewSimulator(repo *Repository, svc *Service, mqtt *MQTTWorker, logger *slog.Logger) *Simulator {
	return &Simulator{repo: repo, svc: svc, mqtt: mqtt, logger: logger}
}

func (s *Simulator) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()

	tenantID, floorID, tags, err := s.repo.EnsureDemoTags(ctx, DEMO_TENANT)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.tenantID = tenantID
	s.floorID = floorID
	s.tags = tags
	s.running = true
	s.stopCh = make(chan struct{})
	s.mu.Unlock()

	go s.loop()
	s.logger.Info("location simulator started", "tenant", DEMO_TENANT, "tags", len(tags))
	return nil
}

func (s *Simulator) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	close(s.stopCh)
	s.running = false
}

func (s *Simulator) Running() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Simulator) loop() {
	ticker := time.NewTicker(SIM_TICK)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.tick()
		}
	}
}

func (s *Simulator) tick() {
	s.mu.Lock()
	tenantID := s.tenantID
	tags := s.tags
	s.mu.Unlock()

	for i := range tags {
		tags[i].X += tags[i].DirX
		tags[i].Y += tags[i].DirY

		if tags[i].X < 1 || tags[i].X > SIM_FLOOR_W-1 {
			tags[i].DirX *= -1
			tags[i].X = math.Max(1, math.Min(SIM_FLOOR_W-1, tags[i].X))
		}
		if tags[i].Y < 1 || tags[i].Y > SIM_FLOOR_H-1 {
			tags[i].DirY *= -1
			tags[i].Y = math.Max(1, math.Min(SIM_FLOOR_H-1, tags[i].Y))
		}

		event := PositionEvent{
			TenantID:  tenantID,
			TagID:     tags[i].TagID,
			TagCode:   tags[i].TagCode,
			FloorID:   tags[i].FloorID,
			X:         tags[i].X,
			Y:         tags[i].Y,
			Timestamp: time.Now().UTC(),
		}

		s.mqtt.Publish(event)
		if !s.mqtt.Connected() {
			s.svc.HandleEvent(context.Background(), event)
		}
	}

	s.mu.Lock()
	s.tags = tags
	s.mu.Unlock()
}

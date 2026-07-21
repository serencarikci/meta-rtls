package location

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Service struct {
	repo   *Repository
	hub    *Hub
	logger *slog.Logger

	mu           sync.Mutex
	latest       map[string]LivePosition
	lastZone     map[string]string
	zonesByFloor map[string][]ZoneBox
}

func NewService(repo *Repository, hub *Hub, logger *slog.Logger) *Service {
	return &Service{
		repo:         repo,
		hub:          hub,
		logger:       logger,
		latest:       map[string]LivePosition{},
		lastZone:     map[string]string{},
		zonesByFloor: map[string][]ZoneBox{},
	}
}

func (s *Service) HandleEvent(ctx context.Context, event PositionEvent) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	if err := s.repo.InsertEvent(ctx, event); err != nil {
		s.logger.Warn("location insert failed", "err", err)
	}

	zoneCode, zoneID := s.findZone(ctx, event.TenantID, event.FloorID, event.X, event.Y)
	s.updateZoneEvents(ctx, event, zoneID)

	pos := LivePosition{
		TenantID: event.TenantID,
		TagID:    event.TagID,
		TagCode:  event.TagCode,
		FloorID:  event.FloorID,
		X:        event.X,
		Y:        event.Y,
		ZoneCode: zoneCode,
		Updated:  event.Timestamp,
	}

	key := event.TenantID + "|" + event.TagID
	s.mu.Lock()
	s.latest[key] = pos
	s.mu.Unlock()

	s.hub.Broadcast(pos)
}

func (s *Service) LatestForTenant(tenantID string) []LivePosition {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []LivePosition
	for _, pos := range s.latest {
		if pos.TenantID == tenantID {
			out = append(out, pos)
		}
	}
	return out
}

func (s *Service) findZone(ctx context.Context, tenantID, floorID string, x, y float64) (code, id string) {
	key := tenantID + "|" + floorID
	s.mu.Lock()
	zones, ok := s.zonesByFloor[key]
	s.mu.Unlock()
	if !ok {
		loaded, err := s.repo.ListZonesForFloor(ctx, tenantID, floorID)
		if err != nil {
			s.logger.Warn("load zones failed", "err", err)
			return "", ""
		}
		s.mu.Lock()
		s.zonesByFloor[key] = loaded
		zones = loaded
		s.mu.Unlock()
	}

	for _, z := range zones {
		if x >= z.MinX && x <= z.MaxX && y >= z.MinY && y <= z.MaxY {
			return z.Code, z.ID
		}
	}
	return "", ""
}

func (s *Service) updateZoneEvents(ctx context.Context, event PositionEvent, newZoneID string) {
	key := event.TenantID + "|" + event.TagID
	s.mu.Lock()
	oldZoneID := s.lastZone[key]
	s.lastZone[key] = newZoneID
	s.mu.Unlock()

	if oldZoneID == newZoneID {
		return
	}
	if oldZoneID != "" {
		_ = s.repo.InsertZoneEvent(ctx, event.TenantID, event.TagID, oldZoneID, "ZONE_EXITED", event.Timestamp)
	}
	if newZoneID != "" {
		_ = s.repo.InsertZoneEvent(ctx, event.TenantID, event.TagID, newZoneID, "ZONE_ENTERED", event.Timestamp)
	}
}

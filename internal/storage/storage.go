package storage

import (
	"context"
	"fmt"

	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/interfaces"
	"github.com/and67o/otus_project/internal/model"
	_ "github.com/go-sql-driver/mysql" // nolint: gci
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

const (
	driverName = "mysql"
	clickCount = "count_click"
	showCount  = "count_show"
)

func New(config configuration.DBConf) (interfaces.Storage, error) {
	db, err := sqlx.Open(driverName, dataSourceName(config))
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("check ping db: %w", err)
	}

	return &Storage{db: db}, nil
}

func dataSourceName(config configuration.DBConf) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s",
		config.User,
		config.Pass,
		config.Host,
		config.Port,
		config.DBName,
	)
}

func (s *Storage) AddBanner(ctx context.Context, b *model.BannerPlace) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO rotation (id_banner, id_slot) VALUES (?, ?)",
		b.BannerID,
		b.SlotID,
	)
	if err != nil {
		return fmt.Errorf("add banner storage: %w", err)
	}

	return nil
}

func (s *Storage) DeleteBanner(ctx context.Context, b *model.BannerPlace) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM rotation WHERE id_banner = ? AND id_slot = ?",
		b.BannerID,
		b.SlotID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Banners(ctx context.Context, slotID int64, groupID int64) ([]model.Banner, error) {
	var banners []model.Banner

	sql := fmt.Sprintf("SELECT r.id_banner, r.id_slot, s.count_show, s.count_click " +
		"FROM rotation r " +
		"LEFT join statistics s on r.id_banner = s.id_banner AND r.id_slot = s.id_slot and r.id_slot = ? " +
		"WHERE s.id_group = ? AND r.status = ?",
	)
	res, err := s.db.QueryContext(ctx, sql, slotID, groupID, model.BannerStatusActive)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	for res.Next() {
		var b model.Banner

		err = res.Scan(
			&b.IDBanner,
			&b.IDSlot,
			&b.CountShow,
			&b.CountClick,
		)
		if err != nil {
			return nil, err
		}

		b.Status = model.BannerStatusActive

		banners = append(banners, b)
	}

	return banners, nil
}

func (s *Storage) IncShowCount(ctx context.Context, slotID int64, groupID int64, bannerID int64) error {
	return s.incCount(ctx, slotID, groupID, bannerID, showCount)
}

func (s *Storage) IncClickCount(ctx context.Context, slotID int64, groupID int64, bannerID int64) error {
	return s.incCount(ctx, slotID, groupID, bannerID, clickCount)
}

func (s *Storage) AddStatistics(ctx context.Context, stat *model.Statistics) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO statistics (id_slot, id_banner, id_group, count_click, count_show) VALUES (?, ?, ?, ?, ?)",
		stat.IDSlot,
		stat.IDBanner,
		stat.IDGroup,
		stat.CountClick,
		stat.CountShow,
	)
	if err != nil {
		return fmt.Errorf("add banner storage: %w", err)
	}

	return nil
}

func (s *Storage) DeleteStatistics(ctx context.Context, stat *model.Statistics) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM statistics WHERE id_banner = ? AND id_slot = ? AND id_group = ?",
		stat.IDBanner,
		stat.IDSlot,
		stat.IDGroup,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetStatistics(ctx context.Context, stat *model.Statistics) (*model.Statistics, error) {
	var statistics model.Statistics

	res, err := s.db.QueryContext(ctx, "SELECT id_banner, id_slot, id_group, count_show, count_click FROM statistics WHERE id_banner = ? AND id_slot = ? AND id_group = ?",
		stat.IDBanner,
		stat.IDSlot,
		stat.IDGroup,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Close()
	}()

	for res.Next() {
		err = res.Scan(&statistics.IDBanner,
			&statistics.IDSlot,
			&statistics.IDGroup,
			&statistics.CountShow,
			&statistics.CountClick,
		)
		if err != nil {
			return nil, err
		}
	}

	return &statistics, nil
}

func (s *Storage) incCount(ctx context.Context, slotID int64, groupID int64, bannerID int64, value string) error {
	sql := fmt.Sprintf("INSERT INTO statistics (id_slot, id_banner, id_group, %s) "+
		"VALUES (?, ?, ?, 1) "+
		"ON DUPLICATE KEY UPDATE %s = %s + 1", value, value, value)
	_, err := s.db.ExecContext(ctx, sql, slotID, bannerID, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateStatus(ctx context.Context, status model.BannerStatus, b *model.BannerPlace) error {
	_, err := s.db.ExecContext(ctx, "UPDATE rotation set status = ? WHERE id_banner = ? AND id_slot = ?",
		status,
		b.BannerID,
		b.SlotID,
	)
	if err != nil {
		return err
	}

	return nil
}

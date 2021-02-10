package storage

import (
	"fmt"
	"github.com/and67o/otus_project/internal/interfaces"

	"github.com/and67o/otus_project/internal/configuration"
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

func (s *Storage) AddBanner(b *model.BannerPlace) error {
	_, err := s.db.Exec("INSERT INTO rotation (id_banner, id_slot) VALUES (?, ?)",
		b.BannerID,
		b.SlotID,
	)
	if err != nil {
		return fmt.Errorf("add banner storage: %w", err)
	}

	return nil
}

func (s *Storage) DeleteBanner(b *model.BannerPlace) error {
	return s.updateStatus(model.BannerStatusDeleted, b)
}

func (s *Storage) Banners(slotID int64, groupID int64) ([]model.Banner, error) {
	var banners []model.Banner

	sql := fmt.Sprintf("SELECT r.id_banner, r.id_slot, s.count_show, s.count_click " +
		"FROM rotation r " +
		"LEFT join statistics s on r.id_banner = s.id_banner AND r.id_slot = s.id_slot and r.id_slot = ? " +
		"WHERE s.id_group = ? AND r.status = ?",
	)
	res, err := s.db.Query(sql, slotID, groupID, model.BannerStatusActive)
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

func (s *Storage) IncShowCount(slotID int64, groupID int64, bannerID int64) error {
	return s.incCount(slotID, groupID, bannerID, showCount)
}

func (s *Storage) IncClickCount(slotID int64, groupID int64, bannerID int64) error {
	return s.incCount(slotID, groupID, bannerID, clickCount)
}

func (s *Storage) AddStatistics(stat *model.Statistics) error {
	_, err := s.db.Exec("INSERT INTO statistics (id_slot, id_banner, id_group, count_click, count_show) VALUES (?, ?, ?, ?, ?)",
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

func (s *Storage) DeleteStatistics(stat *model.Statistics) error {
	_, err := s.db.Exec("DELETE FROM statistics WHERE id_banner = ? AND id_slot = ? AND id_group = ?",
		stat.IDBanner,
		stat.IDSlot,
		stat.IDGroup,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetStatistics(stat *model.Statistics) (*model.Statistics, error) {
	var statistics model.Statistics

	res, err := s.db.Query("SELECT id_banner, id_slot, id_group, count_show, count_click FROM statistics WHERE id_banner = ? AND id_slot = ? AND id_group = ?",
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

func (s *Storage) incCount(slotID int64, groupID int64, bannerID int64, value string) error {
	sql := fmt.Sprintf("INSERT INTO statistics (id_slot, id_banner, id_group, %s) "+
		"VALUES (?, ?, ?, 1) "+
		"ON DUPLICATE KEY UPDATE %s = %s + 1", value, value, value)
	_, err := s.db.Exec(sql, slotID, bannerID, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) updateStatus(status model.BannerStatus, b *model.BannerPlace) error {
	_, err := s.db.Exec("UPDATE rotation set status = ? WHERE id_banner = ? AND id_slot = ?",
		status,
		b.BannerID,
		b.SlotID,
	)
	if err != nil {
		return err
	}

	return nil
}

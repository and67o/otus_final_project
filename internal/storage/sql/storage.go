package sql

import (
	"fmt"
	"github.com/and67o/otus_project/internal/configuration"
	"github.com/and67o/otus_project/internal/model"
	_ "github.com/go-sql-driver/mysql" // nolint: gci
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

type StorageAction interface {
	AddBanner(b *model.BannerPlace) error
	DeleteBanner(b *model.BannerPlace) error
	Banners(slotId int64, groupId int64) ([]model.Banner, error)
	IncShowCount(slotId int64, groupId int64, bannerId int64) error
	IncClickCount(slotId int64, groupId int64, bannerId int64) error
}

const driverName = "mysql"
const format = "2006-01-02 15:04:05"
const clickCount = "count_click"
const showCount = "count_show"

func New(config configuration.DBConf) (StorageAction, error) {
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
		b.BannerId,
		b.SlotId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteBanner(b *model.BannerPlace) error {
	_, err := s.db.Exec("DELETE FROM rotation WHERE id_banner = ? AND id_slot = ?",
		b.BannerId,
		b.SlotId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Banners(slotId int64, groupId int64) ([]model.Banner, error) {
	var err error
	var banners []model.Banner

	sql := fmt.Sprintf("SELECT r.id_banner, r.id_slot, s.count_show, s.count_click " +
		"from rotation r " +
		"left join statistics s on r.id_banner = s.id_banner and r.id_slot = s.id_slot and r.id_slot = ? " +
		"where id_group = ?")

	res, err := s.db.Query(sql, slotId, groupId)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var b model.Banner

		err = res.Scan(
			&b.ID,
			&b.SlotID,
			&b.ShowCount,
			&b.ClickCount,
		)
		if err != nil {
			return nil, err
		}

		banners = append(banners, b)
	}

	return banners, nil
}

func (s *Storage) IncShowCount(slotId int64, groupId int64, bannerId int64) error {
	err := s.incCount(slotId, groupId, bannerId, showCount)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) IncClickCount(slotId int64, groupId int64, bannerId int64) error {
	err := s.incCount(slotId, groupId, bannerId, clickCount)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) incCount(slotId int64, groupId int64, bannerId int64, value string) error {
	sql := fmt.Sprintf("INSERT INTO statistics (id_slot, id_banner, id_group, %s) "+
		"VALUES (?, ?, ?, 1) "+
		"ON DUPLICATE KEY UPDATE %s = %s + 1", value, value,value)
	_, err := s.db.Exec(sql,  slotId, bannerId, groupId)
	if err != nil {
		return err
	}

	return nil
}

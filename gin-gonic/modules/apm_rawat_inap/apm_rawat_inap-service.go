package apm_rawat_inap

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ApmRawatInapService interface {
	Create(input *ApmRawatInapCreate) (*ApmRawatInap, error)
	Publish(id string) (*ApmRawatInap, error)
}

type apmRawatInapService struct {
	db *gorm.DB
	nc *nats.Conn
}

func NewApmRawatInapService(db *gorm.DB, nc *nats.Conn) ApmRawatInapService {
	return &apmRawatInapService{db: db, nc: nc}

}

func (s *apmRawatInapService) Create(input *ApmRawatInapCreate) (*ApmRawatInap, error) {
	tglAntrian := time.Now().Format("2006-01-02")

	var lastApm ApmRawatInap
	var lastNo int

	err := s.db.Select("no_antrian").
		Where("tgl_antrian = ?", tglAntrian).
		Order("id desc").
		First(&lastApm).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err == nil && len(lastApm.NoAntrian) > 1 {
		if num, err := strconv.Atoi(lastApm.NoAntrian[1:]); err == nil {
			lastNo = num
		}
	}

	noAntrian := fmt.Sprintf("U%03d", lastNo+1)

	apm := &ApmRawatInap{
		ID:           time.Now().Format("20060102150405"),
		TglAntrian:   tglAntrian,
		WaktuAntrian: time.Now().Format("15:04:05"),
		NoAntrian:    noAntrian,
		NoRm:         input.NoRm,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(apm).Error; err != nil {
		return nil, err
	}

	return apm, nil
}

func (s *apmRawatInapService) Publish(id string) (*ApmRawatInap, error) {
	var apm ApmRawatInap
	query := s.db.Where("id = ?", id)
	if err := query.First(&apm).Error; err != nil {
		return nil, err
	}

	payload := ApmRawatInapPayloadNats{
		ID:           apm.ID,
		TglAntrian:   apm.TglAntrian,
		WaktuAntrian: apm.WaktuAntrian,
		NoAntrian:    apm.NoAntrian,
	}

	dataJson, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal mengkonversi data ke JSON: %v", err)
	}
	// publisher mengirim perintah dengan topic "antrian.ranap"
	if _, err_req := s.nc.Request("antrian.ranap", dataJson, 500*time.Millisecond); err_req != nil {
		return nil, fmt.Errorf("gagal publish ke NATS: %v", err_req)
	}

	log.Printf("Berhasil publish apm_rawat_inap untuk ID: %s", id)
	return &apm, nil
}

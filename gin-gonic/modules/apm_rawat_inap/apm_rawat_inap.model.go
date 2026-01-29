package apm_rawat_inap

import (
	"time"
)

type ApmRawatInap struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	TglAntrian   string    `json:"tgl_antrian"`
	WaktuAntrian string    `json:"waktu_antrian"`
	NoAntrian    string    `json:"no_antrian"`
	NoRm         string    `json:"no_rm"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ApmRawatInapCreate struct {
	NoRm string `json:"no_rm" binding:"required"`
}

type ApmRawatInapPayloadNats struct {
	ID           string `json:"id"`
	TglAntrian   string `json:"tgl_antrian"`
	WaktuAntrian string `json:"waktu_antrian"`
	NoAntrian    string `json:"no_antrian"`
}

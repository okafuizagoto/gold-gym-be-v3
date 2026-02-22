package goldgym

import (
	"gopkg.in/guregu/null.v3/zero"
)

type GetGoldUser struct {
	GoldId            int    `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldEmail         string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
	GoldPassword      string `gorm:"column:gold_password" db:"gold_password" json:"gold_password"`
	GoldNama          string `gorm:"column:gold_nama" db:"gold_nama" json:"gold_nama"`
	GoldNomorHp       string `gorm:"column:gold_nomorhp" db:"gold_nomorhp" json:"gold_nomorhp"`
	GoldNomorKartu    string `gorm:"column:gold_nomorkartu" db:"gold_nomorkartu" json:"gold_nomorkartu"`
	GoldCvv           string `gorm:"column:gold_cvv" db:"gold_cvv" json:"gold_cvv"`
	GoldExpireddate   string `gorm:"column:gold_expireddate" db:"gold_expireddate" json:"gold_expireddate"`
	GoldPemegangKartu string `gorm:"column:gold_namapemegangkartu" db:"gold_namapemegangkartu" json:"gold_namapemegangkartu"`
}

type GetGoldUsers struct {
	GoldId            int    `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldEmail         string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
	GoldPassword      string `gorm:"column:gold_password" db:"gold_password" json:"gold_password"`
	GoldNama          string `gorm:"column:gold_nama" db:"gold_nama" json:"gold_nama"`
	GoldNomorHp       string `gorm:"column:gold_nomorhp" db:"gold_nomorhp" json:"gold_nomorhp"`
	GoldNomorKartu    string `gorm:"column:gold_nomorkartu" db:"gold_nomorkartu" json:"gold_nomorkartu"`
	GoldCvv           string `gorm:"column:gold_cvv" db:"gold_cvv" json:"gold_cvv"`
	GoldExpireddate   string `gorm:"column:gold_expireddate" db:"gold_expireddate" json:"gold_expireddate"`
	GoldPemegangKartu string `gorm:"column:gold_namapemegangkartu" db:"gold_namapemegangkartu" json:"gold_namapemegangkartu"`
	GoldOTP           string `gorm:"column:gold_otp" db:"gold_otp" json:"gold_otp"`
}

type GetGoldUserss struct {
	GoldId                  int         `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldEmail               string      `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
	GoldPassword            string      `gorm:"column:gold_password" db:"gold_password" json:"gold_password"`
	GoldNama                string      `gorm:"column:gold_nama" db:"gold_nama" json:"gold_nama"`
	GoldNomorHp             string      `gorm:"column:gold_nomorhp" db:"gold_nomorhp" json:"gold_nomorhp"`
	GoldNomorKartu          string      `gorm:"column:gold_nomorkartu" db:"gold_nomorkartu" json:"gold_nomorkartu"`
	GoldCvv                 string      `gorm:"column:gold_cvv" db:"gold_cvv" json:"gold_cvv"`
	GoldExpireddate         string      `gorm:"column:gold_expireddate" db:"gold_expireddate" json:"gold_expireddate"`
	GoldPemegangKartu       string      `gorm:"column:gold_namapemegangkartu" db:"gold_namapemegangkartu" json:"gold_namapemegangkartu"`
	GoldValidasiYN          string      `gorm:"column:gold_validasiyn" db:"gold_validasiyn" json:"gold_validasiyn"`
	GoldToken               zero.String `gorm:"column:gold_token" db:"gold_token" json:"gold_token"`
	GoldOtp                 string      `gorm:"column:gold_otp" db:"gold_otp" json:"gold_otp"`
	GoldUpdatedBy           string      `gorm:"column:gold_updated_by" db:"gold_updated_by" json:"gold_updated_by"`
	GoldUpdatedAt           string      `gorm:"column:gold_updated_at" db:"gold_updated_at" json:"gold_updated_at"`
	GoldLastLogin           string      `gorm:"column:gold_last_login" db:"gold_last_login" json:"gold_last_login"`
	GoldLastLoginHost       string      `gorm:"column:gold_last_login_host" db:"gold_last_login_host" json:"gold_last_login_host"`
	GoldForceChangePassword int         `gorm:"column:gold_force_change_password" db:"gold_force_change_password" json:"gold_force_change_password"`
}

type LoginUser struct {
	GoldNama          string `gorm:"column:gold_nama" db:"gold_nama" json:"gold_nama"`
	GoldNomorHp       string `gorm:"column:gold_nomorhp" db:"gold_nomorhp" json:"gold_nomorhp"`
	GoldNomorKartu    string `gorm:"column:gold_nomorkartu" db:"gold_nomorkartu" json:"gold_nomorkartu"`
	GoldCvv           string `gorm:"column:gold_cvv" db:"gold_cvv" json:"gold_cvv"`
	GoldExpireddate   string `gorm:"column:gold_expireddate" db:"gold_expireddate" json:"gold_expireddate"`
	GoldPemegangKartu string `gorm:"column:gold_namapemegangkartu" db:"gold_namapemegangkartu" json:"gold_namapemegangkartu"`
	GoldToken         string `json:"gold_token"`
}

type LogUser struct {
	GoldEmail    string `json:"gold_email"`
	GoldPassword string `json:"gold_password"`
}

// type DeleteSubsHeader struct {
// 	GoldId int `gorm:"column:gold_id" json:"gold_id"`
// 	// GoldMenuId int `gorm:"column:gold_menuid" json:"gold_menuid"`
// }

type UpdatePassword struct {
	GoldPassword string `gorm:"column:gold_password" db:"gold_password" json:"gold_password"`
	GoldEmail    string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
	GoldOTP      string `gorm:"column:gold_otp" db:"gold_otp" json:"gold_otp"`
}

type UpdateNama struct {
	GoldNama  string `gorm:"column:gold_nama" db:"gold_nama" json:"gold_nama"`
	GoldEmail string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
}

type UpdateKartu struct {
	GoldNomorKartu string `gorm:"column:gold_nomorkartu" db:"gold_nomorkartu" json:"gold_nomorkartu"`
	GoldCvv        string `gorm:"column:gold_cvv" db:"gold_cvv" json:"gold_cvv"`
	GoldEmail      string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
}

type Logout struct {
	GoldEmail string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
}

type GetSubsWithUser struct {
	GoldId              int         `gorm:"column:gold_id" db:"gold_id" json:"gold_id"`
	GoldMenuId          zero.String `gorm:"column:gold_menuid" db:"gold_menuid" json:"gold_menuid"`
	GoldEmail           string      `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
	GoldNama            string      `gorm:"column:gold_nama" db:"gold_nama" json:"gold_nama"`
	GoldNomorHp         string      `gorm:"column:gold_nomorhp" db:"gold_nomorhp" json:"gold_nomorhp"`
	GoldExpireddate     string      `gorm:"column:gold_expireddate" db:"gold_expireddate" json:"gold_expireddate"`
	GoldNamaPaket       zero.String `gorm:"column:gold_namapaket" db:"gold_namapaket" json:"gold_namapaket"`
	GoldNamaLayanan     zero.String `gorm:"column:gold_namalayanan" db:"gold_namalayanan" json:"gold_namalayanan"`
	GoldHarga           zero.Float  `gorm:"column:gold_harga" db:"gold_harga" json:"gold_harga"`
	GoldListLatihan     zero.String `gorm:"column:gold_listlatihan" db:"gold_listlatihan" json:"gold_listlatihan"`
	GoldJumlahpertemuan zero.Int    `gorm:"column:gold_jumlahpertemuan" db:"gold_jumlahpertemuan" json:"gold_jumlahpertemuan"`
	GoldDurasi          zero.Int    `gorm:"column:gold_durasi" db:"gold_durasi" json:"gold_durasi"`
	GoldStatuslangganan zero.String `gorm:"column:gold_statuslangganan" db:"gold_statuslangganan" json:"gold_statuslangganan"`
}

type GetValidationGoldOTP struct {
	GoldOTP string `gorm:"column:gold_otp" db:"gold_otp" json:"gold_otp"`
}

type UpdateValidationOTP struct {
	GoldEmail string `gorm:"column:gold_email" db:"gold_email" json:"gold_email"`
}

func (GetGoldUser) TableName() string {
	return "data_peserta"
}

func (GetGoldUserss) TableName() string {
	return "data_peserta"
}

func (LoginUser) TableName() string {
	return "data_peserta"
}

func (GetGoldUsers) TableName() string {
	return "data_peserta"
}

func (GetValidationGoldOTP) TableName() string {
	return "data_peserta"
}

func (UpdateValidationOTP) TableName() string {
	return "data_peserta"
}

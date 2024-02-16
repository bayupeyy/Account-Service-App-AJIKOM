package users

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID       int `gorm:"primaryKey"`
	Nama     string
	HP       string
	Email    string
	Password string
	Alamat   string
	Saldo    float64
}

type RiwayatTopUp struct {
	UserID    int `gorm:"references:ID"`
	Amount    float64
	Timestamp time.Time
}

type RiwayatTransfer struct {
	UserID    int `gorm:"references:ID"`
	Penerima  string
	Amount    float64
	Timestamp time.Time
}

func AutoMigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&Users{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&RiwayatTopUp{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&RiwayatTransfer{}); err != nil {
		return err
	}
	return nil
}

func (u *Users) GantiPassword(connection *gorm.DB, newPassword string) (bool, error) {
	query := connection.Table("users").Where("hp = ?", u.HP).Update("password", newPassword)
	if err := query.Error; err != nil {
		return false, err
	}

	return query.RowsAffected > 0, nil
}

func Register(connection *gorm.DB, newUser Users) (bool, error) {
	query := connection.Create(&newUser)
	if err := query.Error; err != nil {
		return false, err
	}

	return query.RowsAffected > 0, nil
}

// Fungsi untuk Update Profil
func UpdateProfil(db *gorm.DB, userID int, profilBaru Users) error {
	//Mencari user berdasarkan ID

	var user Users //Membuat Variabel user dari Tbl Users
	err := db.First(&user, userID).Error
	if err != nil {
		return err
	}

	// Update profil user dengan data baru jika tidak kosong
	if profilBaru.Nama != "" { //Melakukan pengecekan apakah ada input kosong
		user.Nama = profilBaru.Nama
	}
	if profilBaru.HP != "" {
		user.HP = profilBaru.HP
	}
	if profilBaru.Email != "" {
		user.Email = profilBaru.Email
	}
	if profilBaru.Password != "" {
		user.Password = profilBaru.Password
	}
	if profilBaru.Alamat != "" {
		user.Alamat = profilBaru.Alamat
	}

	//Menyimpan perubahan ke dalam DB
	err = db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// Fungsi untuk Menampilkan user profil
func MenampilkanProfilUser(db *gorm.DB, userID int) error {
	var user Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	fmt.Println("=== Profil Pengguna ===")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("Nama: %s\n", user.Nama)
	fmt.Printf("Nomor HP: %s\n", user.HP)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Alamat: %s\n", user.Alamat)
	fmt.Printf("Saldo: %.2f\n", user.Saldo)

	return nil
}

// Fungsi Delete Users
func DeleteUser(db *gorm.DB, userID int) error {
	// Temukan user berdasarkan ID
	var user Users
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	// Hapus user dari database
	if err := db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

// TopUpSaldo menambahkan saldo pengguna berdasarkan ID pengguna dan jumlah saldo yang ditambahkan
func TopUpSaldo(db *gorm.DB, userID int, amount float64) error {

	var user Users
	if err := db.First(&user).Error; err != nil {
		return err
	}

	// Tambahkan saldo
	user.Saldo += amount

	// Simpan perubahan ke database
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	// Simpan riwayat topup ke database
	if err := SimpanRiwayatTopUp(db, userID, amount); err != nil {
		return err
	}

	return nil
}

// Simpan Riwayat TopUp
func SimpanRiwayatTopUp(db *gorm.DB, userID int, amount float64) error {
	history := RiwayatTopUp{
		UserID:    userID,
		Amount:    amount,
		Timestamp: time.Now(),
	}
	if err := db.Create(&history).Error; err != nil {
		return err
	}
	return nil
}

func GetTopUpHistory(db *gorm.DB, userID int) ([]RiwayatTopUp, error) {
	var history []RiwayatTopUp
	if err := db.Where("user_id = ?", userID).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

// Fungsi Transfer Saldo
func TransferSaldo(db *gorm.DB, senderID int, receiverHP int, amount float64) (bool, error) {
	var sender Users
	var receiver Users

	// Mencari pengguna pengirim berdasarkan ID
	if err := db.First(&sender, senderID).Error; err != nil {
		return false, err
	}

	// Mencari pengguna penerima berdasarkan nomor HP
	if err := db.Where("hp = ?", receiverHP).First(&receiver).Error; err != nil {
		return false, err
	}

	// Memastikan saldo pengirim mencukupi untuk transfer
	if sender.Saldo < amount {
		return false, nil
	}

	// Melakukan pengurangan saldo dari pengirim
	sender.Saldo -= amount
	if err := db.Save(&sender).Error; err != nil {
		return false, err
	}
	// Menambahkan saldo ke penerima
	receiver.Saldo += amount
	if err := db.Save(&receiver).Error; err != nil {
		return false, err
	}

	// Simpan riwayat topup ke database
	if err := SimpanRiwayatTransfer(db, sender.ID, receiver.Nama, amount); err != nil {
		return false, err
	}

	return true, nil
}

// Simpan Riwayat Transfer
func SimpanRiwayatTransfer(db *gorm.DB, userID int, penerima string, amount float64) error {
	history := RiwayatTransfer{
		UserID:    userID,
		Amount:    amount,
		Penerima:  penerima,
		Timestamp: time.Now(),
	}
	if err := db.Create(&history).Error; err != nil {
		return err
	}
	return nil
}

func SemuaRiwayatTransfer(db *gorm.DB, userID int) ([]RiwayatTransfer, error) {
	var history []RiwayatTransfer
	if err := db.Where("user_id = ?", userID).Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

// Melihat User dengan No Hp
func LihatProfilPenggunaByHP(db *gorm.DB, hp string) error {
	var user Users
	if err := db.Where("hp = ?", hp).First(&user).Error; err != nil {
		return err
	}

	// Tampilkan profil pengguna
	fmt.Println("=== Profil Pengguna ===")
	fmt.Printf("Nama: %s\n", user.Nama)
	fmt.Printf("Nomor HP: %s\n", user.HP)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Alamat: %s\n", user.Alamat)

	return nil
}

// func LihatUser(db *gorm.DB)
func Login(connection *gorm.DB, hp string, password string) (Users, error) {
	var result Users
	err := connection.Where("hp = ? AND password = ?", hp, password).First(&result).Error
	if err != nil {
		return Users{}, err
	}

	return result, nil
}

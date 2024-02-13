package users

import (
	"fmt"

	"gorm.io/gorm"
)

type Users struct {
	ID       int `gorm:"primaryKey"`
	Nama     string
	HP       string
	Email    string
	Password string
	Alamat   string
	Saldo    string
}

type Barang struct {
	ID_Barang     int `gorm:"primaryKey"`
	UserID        int
	Nama_Barang   string
	Harga         float64
	Jumlah_barang int
}

func AutoMigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&Users{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&Barang{}); err != nil {
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

func Login(connection *gorm.DB, hp string, password string) (Users, error) {
	var result Users
	err := connection.Where("hp = ? AND password = ?", hp, password).First(&result).Error
	if err != nil {
		return Users{}, err
	}

	return result, nil
}

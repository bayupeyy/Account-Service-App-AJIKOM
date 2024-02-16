package main

import (
	"be21/config"
	"be21/users"
	"fmt"
)

func main() {
	database := config.InitMysql()
	config.Migrate(database)
	var input int
	for input != 99 {
		fmt.Println("Selamat Datang di ToApp")
		fmt.Println("Pilih menu")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("0. Exit")
		fmt.Print("Masukkan pilihan:")
		fmt.Scanln(&input)
		if input == 1 {
			var isRunning bool = true
			for isRunning {
				var hp string
				var password string
				var loggedIn users.Users
				fmt.Println("Masukkan HP")
				fmt.Scanln(&hp)
				fmt.Println("Masukkan Password")
				fmt.Scanln(&password)
				loggedIn, err := users.Login(database, hp, password)
				if err == nil {
					fmt.Println("Selamat Datang,", loggedIn.Nama)

					// Menu input barang
					var inputMenu int
					for inputMenu != 99 {
						fmt.Println("Menu")
						fmt.Println("1. Read Account")
						fmt.Println("2. Update Account")
						fmt.Println("3. Delete Account")
						fmt.Println("4. TopUp")
						fmt.Println("5. Transfer")
						fmt.Println("6. History TopUp")
						fmt.Println("7. History Transfer")
						fmt.Println("8. Lihat Profil Lain")
						fmt.Println("0. Kembali ke Menu Utama")
						fmt.Print("Masukkan Pilihan : ")
						fmt.Scanln(&inputMenu)

						if inputMenu == 1 {
							//Menampilkan profil User
							err := users.MenampilkanProfilUser(database, loggedIn.ID)
							if err != nil {
								fmt.Println("Gagal Menampilkan Profil: ", err)
							}
						} else if inputMenu == 2 {
							//Update Profil Baru
							var profilBaru users.Users
							fmt.Println(" Update Profil User ")
							fmt.Print("Masukkan nama baru: ")
							fmt.Scanln(&profilBaru.Nama)
							fmt.Print("Masukkan nomor HP baru: ")
							fmt.Scanln(&profilBaru.HP)
							fmt.Print("Masukkan Email baru: ")
							fmt.Scanln(&profilBaru.Email)
							fmt.Print("Masukkan password baru: ")
							fmt.Scanln(&profilBaru.Password)
							fmt.Print("Masukkan alamat baru: ")
							fmt.Scanln(&profilBaru.Alamat)

							//Proses Update ke Fungsi dan Masuk ke DB
							err := users.UpdateProfil(database, loggedIn.ID, profilBaru)
							if err != nil {
								fmt.Println("Gagal mengupdate profil: ", err)
							} else {
								fmt.Println("Profil berhasil diupdate")
							}

						} else if inputMenu == 3 {
							err := users.DeleteUser(database, loggedIn.ID)
							if err != nil {
								fmt.Println("Gagal menghapus user:", err)
							} else {
								fmt.Println("User berhasil dihapus")
							}
						} else if inputMenu == 4 {
							var amount float64
							fmt.Println("Masukkan jumlah saldo yang ingin ditambahkan:")
							fmt.Scanln(&amount)
							err := users.TopUpSaldo(database, loggedIn.ID, amount)
							if err != nil {
								fmt.Println("Gagal top-up saldo:", err)
							} else {
								fmt.Println("Top-up saldo berhasil")
							}
						} else if inputMenu == 5 {
							var receiverHP int
							var amount float64
							fmt.Println("Transfer Saldo")
							fmt.Println("Masukkan Nomor HP Penerima:")
							fmt.Scanln(&receiverHP)
							fmt.Println("Masukkan Jumlah Saldo yang Akan Ditransfer:")
							fmt.Scanln(&amount)

							success, err := users.TransferSaldo(database, loggedIn.ID, receiverHP, amount)
							if err != nil {
								fmt.Println("Gagal melakukan transfer saldo:", err)
							} else if success {
								fmt.Println("Transfer saldo berhasil!")
							} else {
								fmt.Println("Gagal melakukan transfer saldo. Saldo tidak mencukupi.")
							}
						} else if inputMenu == 6 {
							history, err := users.GetTopUpHistory(database, loggedIn.ID)
							if err != nil {
								fmt.Println("Gagal mendapatkan Riwayat topup: ", err)
							} else {
								fmt.Println("=== Riwayat TopUp ===")
								for _, h := range history {
									fmt.Printf("Rp.%.2f, Waktu Transaksi: %s\n", h.Amount, h.Timestamp)
								}
							}
						} else if inputMenu == 7 {
							history, err := users.SemuaRiwayatTransfer(database, loggedIn.ID)
							if err != nil {
								fmt.Println("Gagal mendapatkan Riwayat Transfer: ", err)
							} else {
								fmt.Println("=== Riwayat Transfer ===")
								for _, h := range history {
									fmt.Printf("Rp.%.2f, Penerima: %s, Waktu Transaksi: %s\n", h.Amount, h.Penerima, h.Timestamp)
								}
							}
						} else if inputMenu == 8 {
							var userHP string
							fmt.Println("Masukkan nomor HP pengguna lain:")
							fmt.Scanln(&userHP)

							if err := users.LihatProfilPenggunaByHP(database, userHP); err != nil {
								fmt.Println("Gagal menampilkan profil pengguna:", err)
							}
						} else if inputMenu == 0 {
							// Kembali ke Menu Utama
							break
						} else {
							fmt.Println("Pilihan Tidak Sesuai.")
						}
					}
				} else {
					var inputExit string
					fmt.Print("Input 'EXIT' untuk kembali ke menu sebelumnya ")
					fmt.Scanln(&inputExit)
					if inputExit == "EXIT" {
						isRunning = false
					}
				}
			}

			// kalo sukses welcome, kalo gagal isi lagi
		} else if input == 2 {
			var newUser users.Users
			fmt.Print("Masukkan nama ")
			fmt.Scanln(&newUser.Nama)
			fmt.Print("Masukkan nomor HP ")
			fmt.Scanln(&newUser.HP)
			fmt.Print("Masukkan Email")
			fmt.Scanln(&newUser.Email)
			fmt.Print("Masukkan password ")
			fmt.Scanln(&newUser.Password)
			fmt.Print("Masukkan alamat ")
			fmt.Scanln(&newUser.Alamat)
			success, err := users.Register(database, newUser)
			if err != nil {
				fmt.Println("terjadi kesalahan(tidak bisa mendaftarkan pengguna)", err.Error())
			}

			if success {
				fmt.Printf("Selamat %s, anda telah terdaftar\n", newUser.Nama)
			}
		}
	}
	fmt.Println("Exited! Thank you")

}

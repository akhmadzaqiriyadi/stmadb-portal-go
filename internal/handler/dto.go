// internal/handler/dto.go
package handler

// GenericResponse adalah struktur dasar untuk semua respons JSON.
type GenericResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty agar tidak tampil jika nil
}

// LoginResponse adalah struktur data untuk respons login yang sukses.
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

// ProfileData adalah representasi data profil yang akan dikirim ke klien.
// Ini penting untuk menyembunyikan field seperti password.
type ProfileData struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	// Tambahkan detail Teacher atau Student di sini jika perlu
}

// ProfileResponse adalah struktur lengkap untuk respons get profile.
type ProfileResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    ProfileData `json:"data"`
}
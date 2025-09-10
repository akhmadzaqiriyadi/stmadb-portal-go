package database

import (
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db" 
)

// NewClient membuat instance baru dari Prisma Client.
func NewClient() *db.PrismaClient {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		// Sebaiknya gunakan panic di sini saat startup jika koneksi gagal
		panic(err)
	}
	return client
}
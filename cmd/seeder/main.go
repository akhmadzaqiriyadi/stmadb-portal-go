// cmd/seeder/main.go
package main

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"

	// Path yang benar ke client yang di-generate
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db"
)

func main() {
	log.Println("🌱 Starting database seeding...")
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer client.Disconnect()

	ctx := context.Background()

	// === HASH PASSWORDS ===
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	teacherPassword, _ := bcrypt.GenerateFromPassword([]byte("teacher123"), bcrypt.DefaultCost)
	studentPassword, _ := bcrypt.GenerateFromPassword([]byte("student123"), bcrypt.DefaultCost)

	// === SEED USERS ===
	log.Println("Seeding users...")

	// 1. Admin User
	adminUser, err := client.User.UpsertOne(
		db.User.Username.Equals("admin"),
	).Create(
		db.User.Username.Set("admin"),
		db.User.Password.Set(string(adminPassword)),
		db.User.Role.Set("admin"),
	).Update().Exec(ctx)
	if err != nil {
		log.Fatalf("could not seed admin user: %v", err)
	}

	// 2. Teacher User & Profile
	teacherUser, err := client.User.UpsertOne(
		db.User.Username.Equals("teacher001"),
	).Create(
		db.User.Username.Set("teacher001"),
		db.User.Password.Set(string(teacherPassword)),
		db.User.Role.Set("teacher"), // <-- Menggunakan string
	).Update().Exec(ctx)
	if err != nil {
		log.Fatalf("could not seed teacher user: %v", err)
	}

	// Create teacher profile
	_, err = client.Teacher.UpsertOne(
		db.Teacher.UserID.Equals(teacherUser.ID),
	).Create(
		db.Teacher.FullName.Set("Budi Santoso, S.Pd"),
		db.Teacher.EmploymentStatus.Set("ASN"),
		db.Teacher.User.Link(db.User.ID.Equals(teacherUser.ID)),
		db.Teacher.Nip.Set("196801011990031001"),
		db.Teacher.Nik.Set("3201010101680001"),
		db.Teacher.PhoneNumber.Set("081234567890"),
	).Update().Exec(ctx)
	if err != nil {
		log.Fatalf("could not seed teacher profile: %v", err)
	}

	// 3. Student User & Profile
	studentUser, err := client.User.UpsertOne(
		db.User.Username.Equals("student001"),
	).Create(
		db.User.Username.Set("student001"),
		db.User.Password.Set(string(studentPassword)),
		db.User.Role.Set("student"),
	).Update().Exec(ctx)
	if err != nil {
		log.Fatalf("could not seed student user: %v", err)
	}

	// Create student profile
	_, err = client.Student.UpsertOne(
		db.Student.UserID.Equals(studentUser.ID),
	).Create(
		db.Student.Nis.Set("2024001"),
		db.Student.FullName.Set("Siti Nurhaliza"),
		db.Student.Gender.Set("P"), // <-- Menggunakan string
		db.Student.User.Link(db.User.ID.Equals(studentUser.ID)),
		db.Student.Nisn.Set("0012345678"),
		db.Student.Address.Set("Jl. Pendidikan No. 123, Jakarta"),
		db.Student.PhoneNumber.Set("081234567891"),
	).Update().Exec(ctx)
	if err != nil {
		log.Fatalf("could not seed student profile: %v", err)
	}

	log.Println("✅ Users seeded successfully")
	log.Println("🎉 Database seeding completed successfully!")
	log.Printf("Admin User ID: %s, Teacher User ID: %s, Student User ID: %s\n", adminUser.ID, teacherUser.ID, studentUser.ID)
}
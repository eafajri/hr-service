package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Open file for writing
	file, err := os.Create("users.csv")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	// Initialize CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	err = writer.Write([]string{"email", "password", "role"})
	if err != nil {
		log.Fatalf("failed to write header: %v", err)
	}

	// Write user records
	for i := 1; i <= 100; i++ {
		email := fmt.Sprintf("user%d@example.com", i)
		password := email
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}

		role := "employee"
		if i%2 == 0 {
			role = "admin"
		}

		err = writer.Write([]string{email, string(hashedPassword), role})
		if err != nil {
			log.Fatalf("failed to write record: %v", err)
		}
	}

	fmt.Println("✅ users.csv has been generated successfully.")
}

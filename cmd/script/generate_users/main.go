package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	file, err := os.Create("./cmd/script/users_dummy.csv")
	if err != nil {
		log.Fatalf("Cannot create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"id", "username", "full_name", "password_hash", "salary", "role"})
	if err != nil {
		log.Fatalf("Cannot write header: %v", err)
	}

	for i := 1; i <= 100; i++ {
		fullName := fmt.Sprintf("User %03d", i)
		username := fmt.Sprintf("user%03d", i)
		plainPassword := fmt.Sprintf("%d%s", i, fullName)
		randomFactor := rand.Intn(51)
		salary := float64(5000+randomFactor*100) * 1000

		hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password for user %d: %v", i, err)
		}

		record := []string{
			strconv.Itoa(i),
			username,
			fullName,
			string(hashed),
			fmt.Sprintf("%.2f", salary),
			"employee",
		}

		err = writer.Write(record)
		if err != nil {
			log.Fatalf("Failed to write record for user %d: %v", i, err)
		}
	}

	fmt.Println("CSV file 'users_dummy.csv' created successfully.")
}

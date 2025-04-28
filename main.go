package main

import (
	"bufio"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

type Phrase struct {
	Text   string
	Author string
}

func main() {
	// Ruta del archivo de texto
	filePath := "frases.txt"

	// Conexión a la base de datos
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Leer el archivo y procesar las frases
	err = processFile(filePath, db)
	if err != nil {
		log.Fatalf("Error al procesar el archivo: %v", err)
	}
}

func processFile(filePath string, db *sql.DB) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "—", 2)
		if len(parts) != 2 {
			log.Printf("Línea inválida: %s", line)
			continue
		}

		phrase := Phrase{
			Text:   strings.TrimSpace(parts[0]),
			Author: strings.TrimSpace(parts[1]),
		}

		err := savePhraseToDB(db, phrase)
		if err != nil {
			log.Printf("Error al guardar la frase: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func savePhraseToDB(db *sql.DB, phrase Phrase) error {
	query := `
		INSERT INTO phrases (phrase, author, owner, state, created_at)
		VALUES ($1, $2, 'cheojeg', 'published', NOW())
	`
	_, err := db.Exec(query, phrase.Text, phrase.Author)
	return err
}

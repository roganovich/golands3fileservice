package handlers

import (
	"golands3fileservice/pkg/models"
	"golands3fileservice/pkg/database"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"io"
	"github.com/google/uuid"
	"time"
	"path/filepath"
	"strings"
)

// get team by id
func getOneFile(fileName string) (error, models.File) {
	var file models.File
	err := database.DB.QueryRow("SELECT * FROM files WHERE name = $1", fileName).Scan(
		&file.ID,

		&file.Title,
		&file.Filename,
		&file.Extension,
		&file.Size,
		&file.DateCreate,
	)

	if  err != nil {
		log.Println("Ошибка в getOneFile", fileName, err.Error())
	}

	return err, file
}

// @Summary Загрузить медиафайл
// @Description Загрузка медиафайла
// @Tags Медиафайлы
// @Param file formData file true "Загруженный файл"
// @Success 200 {object} models.File
// @Failure 400 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 415 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/files/upload [post]
func Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			// Устанавливаем заголовки для CORS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Отправляем успешный ответ
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodPost {
			// Загрузка файла
			formFile, fileHeader, errFile := r.FormFile("file")
			if errFile != nil {
				log.Println("Не удалось прочитать файл")
				http.Error(w, "Не удалось прочитать файл", http.StatusBadRequest)
				return
			}
			defer formFile.Close()

			newUUID := uuid.New()
			fileName := newUUID.String()
			dstPath := filepath.Join("./public/uploads/", fileName)

			f, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
			if err != nil {
				log.Println("Не удалось открыть файл")
				http.Error(w, "Не удалось открыть файл", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			fileSize, err := io.Copy(f, formFile)
			if err != nil {
				log.Println("Не удалось скопировать файл")
				http.Error(w, "Не удалось скопировать файл", http.StatusInternalServerError)
				return
			}

			createdAt := time.Now()
			mimeType := getMIMEType(fileHeader.Filename)

			var file models.File
			file.ID = newUUID
			file.Title = fileName
			file.Filename = dstPath
			file.Extension = mimeType
			file.Size = fileSize
			file.DateCreate = createdAt

			errInsert := database.DB.QueryRow("INSERT INTO files (id, title, filename, extension, size) VALUES ($1, $2, $3, $4, $5) RETURNING id", file.ID, file.Title, file.Filename, file.Extension, file.Size).Scan(&file.ID)
			if errInsert != nil {
				http.Error(w, "Возникла ошибка при добавлении изображении (#i)", http.StatusBadRequest)
				return
			}

			json.NewEncoder(w).Encode(file)
			return
		}

		// Если метод не поддерживается
		w.Header().Set("Allow", "POST, OPTIONS")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getMIMEType(filename string) string {
	extFile := filepath.Ext(filename)
	extData := strings.Split(extFile, ".")
	ext := ""
	if len(extData) > 0 {
		ext = extData[1]
	} else {
		log.Println("Расширение файла не удалось получить:" + extFile)
	}

	return ext
}

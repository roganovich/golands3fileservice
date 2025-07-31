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
	"github.com/gorilla/mux"
	"fmt"
	"github.com/minio/minio-go/v7"
)

// @Summary Информация о файле
// @Description Метод просмотра информации о файлах
// @Tags Медиафайлы
// @Accept  json
// @Produce octet-stream
// @Param fileName query string true "File name"
// @Success 200 {file} file
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 409 {string} string "Conflict (multiple files)"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/files/{id}/view [get]
func ViewMedia() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileId, _ := vars["id"]
		var file models.File
		err := database.DB.QueryRow("SELECT id, title, filename, extension, size, date_create FROM files WHERE id = $1", fileId).Scan(
			&file.ID,
			&file.Title,
			&file.Filename,
			&file.Extension,
			&file.Size,
			&file.DateCreate,
		)

		if  err != nil {
			log.Fatalf("Ошибка получения информации о файле %v, %v", fileId, err.Error())
			return
		}
		json.NewEncoder(w).Encode(file)
	}
}

// @Summary Скачать файл
// @Description Метод скачивания файлов
// @Tags Медиафайлы
// @Accept  json
// @Produce octet-stream
// @Param fileName query string true "File name"
// @Success 200 {file} file
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 409 {string} string "Conflict (multiple files)"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/files/{id}/download [get]
func DownloadMedia() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileId, _ := vars["id"]

		var file models.File
		err := database.DB.QueryRow("SELECT id, title, filename, extension, size, date_create FROM files WHERE id = $1", fileId).Scan(
			&file.ID,
			&file.Title,
			&file.Filename,
			&file.Extension,
			&file.Size,
			&file.DateCreate,
		)
		if  err != nil {
			log.Fatalf("Ошибка получения информации о файле %v, %v", fileId, err.Error())
			return
		}
		// Получение файла из MinIO
		bucketName := os.Getenv("MINIO_ROOT_BUCKET_NAME")
		object, err := database.MinioClient.GetObject(
			r.Context(),
			bucketName,
			file.ID.String(),
			minio.GetObjectOptions{},
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка скачивания файла %v", fileId), http.StatusInternalServerError)
			return
		}
		defer object.Close()

		// Отправка файла
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Title))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
		io.Copy(w, object)
	}
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
func UploadMedia() http.HandlerFunc {
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
				http.Error(w, "Не удалось прочитать файл", http.StatusBadRequest)
				return
			}
			defer formFile.Close()

			// Генерация уникального ключа
			newUUID := uuid.New()
			fileName := fileHeader.Filename
			dstPath := filepath.Join("./public/uploads/", fileName)
			mimeType := getMIMEType(fileHeader.Filename)
			createdAt := time.Now()
			fileSize := fileHeader.Size
			// Сохранение метаданных
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

			// Загрузка в MinIO
			bucketName := os.Getenv("MINIO_ROOT_BUCKET_NAME")
			_, err := database.MinioClient.PutObject(
				r.Context(),
				bucketName,
				newUUID.String(), // objectKey
				formFile,
				fileSize,
				minio.PutObjectOptions{ContentType: fileHeader.Header.Get("Content-Type")},
			)
			if err != nil {
				http.Error(w, "Возникла ошибка при загрузке изображения (#us3)", http.StatusInternalServerError)
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

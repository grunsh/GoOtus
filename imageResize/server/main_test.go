package main

import (
	"bytes"
	"context"
	"image/color"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"imageResize/internal/cache"
	"imageResize/internal/processor"
	"imageResize/internal/storage"
)

var (
	testImage   []byte
	imageServer *httptest.Server
)

func TestMain(m *testing.M) {
	// 1. Инициализация тестовых ресурсов
	testImage = createTestImage(800, 600)
	imageServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(testImage)
	}))
	defer imageServer.Close()

	os.Exit(m.Run())
}

func TestImageProcessing(t *testing.T) {
	// 2. Настройка тестового сервера
	storage := storage.NewMemoryStorage()
	cache := cache.NewLRUCache(10, storage)
	imgProcessor := processor.NewImageProcessor(cache)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// Регистрируем обработчики
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/fill/", func(w http.ResponseWriter, r *http.Request) {
		// Парсим параметры из URL
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 5 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		width, err := strconv.Atoi(parts[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		height, err := strconv.Atoi(parts[3])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := strings.Join(parts[4:], "/")
		if url == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Обрабатываем изображение
		data, contentType, err := imgProcessor.ProcessImage(r.Context(), url, width, height)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Write(data)
	})

	// Запускаем сервер
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()
	defer server.Shutdown(context.Background())

	// Ждем готовности сервера
	waitForServerReady(t, "http://localhost:8081/health")

	// 3. Тестовые случаи
	t.Run("Successful image resize", func(t *testing.T) {
		imageURL := imageServer.URL[7:] // убираем "http://"
		resp, err := http.Get("http://localhost:8081/fill/300/200/" + imageURL + "/test.jpg")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))

		imgData, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.True(t, len(imgData) > 0)
		assert.True(t, isJPEG(imgData))
	})
}

func createTestImage(width, height int) []byte {
	img := imaging.New(width, height, color.White)
	var buf bytes.Buffer
	imaging.Encode(&buf, img, imaging.JPEG)
	return buf.Bytes()
}

func waitForServerReady(t *testing.T, url string) {
	client := http.Client{Timeout: 100 * time.Millisecond}
	for start := time.Now(); time.Since(start) < 5*time.Second; time.Sleep(100 * time.Millisecond) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
	}
	t.Fatal("Server did not become ready in time")
}

func isJPEG(data []byte) bool {
	return len(data) > 2 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
}

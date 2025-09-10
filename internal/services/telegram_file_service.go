package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// TelegramFileService handles file downloads from Telegram
type TelegramFileService struct {
	botToken string
	logger   domain.Logger
	client   *http.Client
}

// NewTelegramFileService creates a new Telegram file service
func NewTelegramFileService(botToken string, logger domain.Logger) *TelegramFileService {
	return &TelegramFileService{
		botToken: botToken,
		logger:   logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// DownloadFile downloads a file from Telegram servers to a temporary location
func (s *TelegramFileService) DownloadFile(document *domain.TelegramDocument) (string, error) {
	s.logger.Info("Starting file download", "file_id", document.FileID, "filename", document.FileName)
	
	// 1. Get file info from Telegram
	fileInfo, err := s.getFileInfo(document.FileID)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %v", err)
	}
	
	if fileInfo.FilePath == "" {
		return "", fmt.Errorf("file path not available from Telegram")
	}
	
	// 2. Download file from Telegram servers
	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", s.botToken, fileInfo.FilePath)
	
	resp, err := s.client.Get(downloadURL)
	if err != nil {
		s.logger.Error("Failed to download file from Telegram", "error", err, "url", downloadURL)
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: HTTP %d", resp.StatusCode)
	}
	
	// 3. Create temporary file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("telegram_file_%d_%s", time.Now().Unix(), document.FileName))
	
	file, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer file.Close()
	
	// 4. Copy file content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(tempFile) // Clean up on error
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	
	s.logger.Info("File downloaded successfully", 
		"filename", document.FileName, 
		"size", document.FileSize, 
		"temp_path", tempFile)
	
	return tempFile, nil
}

// getFileInfo gets file information from Telegram Bot API
func (s *TelegramFileService) getFileInfo(fileID string) (*domain.TelegramFile, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", s.botToken, fileID)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Telegram API returned status %d", resp.StatusCode)
	}
	
	var result struct {
		OK     bool                `json:"ok"`
		Result *domain.TelegramFile `json:"result"`
		Description string           `json:"description,omitempty"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	
	if !result.OK {
		return nil, fmt.Errorf("Telegram API error: %s", result.Description)
	}
	
	return result.Result, nil
}

// CleanupFile removes a temporary file
func (s *TelegramFileService) CleanupFile(filePath string) error {
	if filePath == "" {
		return nil
	}
	
	err := os.Remove(filePath)
	if err != nil {
		s.logger.Error("Failed to cleanup temporary file", "file", filePath, "error", err)
		return err
	}
	
	s.logger.Info("Temporary file cleaned up", "file", filePath)
	return nil
}

// GetFileSize returns file size in a human readable format
func (s *TelegramFileService) GetFileSize(size int) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}
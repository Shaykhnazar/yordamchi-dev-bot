package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"

	"yordamchi-dev-bot/internal/domain"
)

// FileExtractor handles content extraction from various file types
type FileExtractor struct {
	logger domain.Logger
}

// NewFileExtractor creates a new file extraction service
func NewFileExtractor(logger domain.Logger) *FileExtractor {
	return &FileExtractor{
		logger: logger,
	}
}

// ExtractContent extracts text content from files based on their type
func (e *FileExtractor) ExtractContent(filePath, fileName string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	
	e.logger.Info("Extracting content from file", "file", fileName, "type", ext)
	
	switch ext {
	case ".txt", ".md":
		return e.extractTextFile(filePath)
	case ".pdf":
		return e.extractPDFContent(filePath)
	case ".docx":
		return e.extractWordContent(filePath)
	case ".xlsx", ".xls":
		return e.extractExcelContent(filePath)
	default:
		return "", fmt.Errorf("unsupported file type: %s. Supported formats: TXT, MD, PDF, DOCX, XLSX", ext)
	}
}

// extractTextFile reads plain text files
func (e *FileExtractor) extractTextFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		e.logger.Error("Failed to read text file", "error", err)
		return "", fmt.Errorf("failed to read text file: %v", err)
	}
	
	text := string(content)
	e.logger.Info("Text file extracted", "length", len(text))
	return text, nil
}

// extractPDFContent extracts text from PDF files
func (e *FileExtractor) extractPDFContent(filePath string) (string, error) {
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		e.logger.Error("Failed to open PDF", "error", err)
		return "", fmt.Errorf("failed to open PDF: %v", err)
	}
	defer file.Close()

	var content strings.Builder
	totalPages := reader.NumPage()
	
	e.logger.Info("Processing PDF", "pages", totalPages)
	
	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}
		
		text, err := page.GetPlainText(nil)
		if err != nil {
			e.logger.Warn("Failed to extract text from page", "page", pageNum, "error", err)
			continue
		}
		
		content.WriteString(text)
		content.WriteString("\n\n")
	}
	
	result := content.String()
	e.logger.Info("PDF extracted", "pages", totalPages, "length", len(result))
	
	if result == "" {
		return "", fmt.Errorf("no text content found in PDF")
	}
	
	return result, nil
}

// extractWordContent extracts text from DOCX files
func (e *FileExtractor) extractWordContent(filePath string) (string, error) {
	doc, err := docx.ReadDocxFile(filePath)
	if err != nil {
		e.logger.Error("Failed to read DOCX file", "error", err)
		return "", fmt.Errorf("failed to read DOCX file: %v", err)
	}
	defer doc.Close()

	docx := doc.Editable()
	content := docx.GetContent()
	
	e.logger.Info("DOCX extracted", "length", len(content))
	
	if content == "" {
		return "", fmt.Errorf("no text content found in DOCX file")
	}
	
	return content, nil
}

// extractExcelContent extracts data from Excel files
func (e *FileExtractor) extractExcelContent(filePath string) (string, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		e.logger.Error("Failed to open Excel file", "error", err)
		return "", fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer file.Close()

	var content strings.Builder
	sheets := file.GetSheetList()
	
	e.logger.Info("Processing Excel file", "sheets", len(sheets))
	
	for _, sheetName := range sheets {
		content.WriteString(fmt.Sprintf("Sheet: %s\n", sheetName))
		content.WriteString("=" + strings.Repeat("=", len(sheetName)+7) + "\n\n")
		
		rows, err := file.GetRows(sheetName)
		if err != nil {
			e.logger.Warn("Failed to read sheet", "sheet", sheetName, "error", err)
			continue
		}
		
		for rowIndex, row := range rows {
			// Skip empty rows
			if len(row) == 0 {
				continue
			}
			
			// Check if row has any content
			hasContent := false
			for _, cell := range row {
				if strings.TrimSpace(cell) != "" {
					hasContent = true
					break
				}
			}
			
			if !hasContent {
				continue
			}
			
			// Format row data
			content.WriteString(fmt.Sprintf("Row %d: ", rowIndex+1))
			content.WriteString(strings.Join(row, " | "))
			content.WriteString("\n")
		}
		
		content.WriteString("\n")
	}
	
	result := content.String()
	e.logger.Info("Excel extracted", "sheets", len(sheets), "length", len(result))
	
	if result == "" {
		return "", fmt.Errorf("no data found in Excel file")
	}
	
	return result, nil
}

// GetSupportedFormats returns list of supported file formats
func (e *FileExtractor) GetSupportedFormats() []string {
	return []string{"TXT", "MD", "PDF", "DOCX", "XLSX"}
}

// IsSupported checks if a file format is supported
func (e *FileExtractor) IsSupported(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	supportedExts := map[string]bool{
		".txt":  true,
		".md":   true,
		".pdf":  true,
		".docx": true,
		".xlsx": true,
		".xls":  true,
	}
	return supportedExts[ext]
}

// ValidateFile performs basic validation on uploaded files
func (e *FileExtractor) ValidateFile(document *domain.TelegramDocument) error {
	// Check file size (limit to 20MB)
	const maxSize = 20 * 1024 * 1024 // 20MB
	if document.FileSize > maxSize {
		return fmt.Errorf("file too large (%.1fMB). Maximum size: 20MB", float64(document.FileSize)/(1024*1024))
	}
	
	// Check if file type is supported
	if !e.IsSupported(document.FileName) {
		return fmt.Errorf("unsupported file type. Supported formats: %s", strings.Join(e.GetSupportedFormats(), ", "))
	}
	
	return nil
}
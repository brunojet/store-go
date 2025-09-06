package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brunojet/store-go/shared/internal/domain"
)

// DefaultStorageService implementa domain.StorageHandler
type DefaultStorageService struct {
	BasePath string // Caminho base para armazenamento local
}

// NewDefaultStorageService cria uma nova instância do serviço de storage
func NewDefaultStorageService(basePath string) *DefaultStorageService {
	return &DefaultStorageService{
		BasePath: basePath,
	}
}

// DeleteFile deleta um arquivo do storage
func (s *DefaultStorageService) DeleteFile(ctx context.Context, caminho string) error {
	fullPath := filepath.Join(s.BasePath, caminho)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("erro ao deletar arquivo %s: %w", fullPath, err)
	}

	fmt.Printf("Arquivo deletado com sucesso: %s\n", fullPath)
	return nil
}

// FileExists verifica se um arquivo existe no storage
func (s *DefaultStorageService) FileExists(ctx context.Context, caminho string) bool {
	fullPath := filepath.Join(s.BasePath, caminho)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// UploadFile faz upload de um arquivo para o storage
func (s *DefaultStorageService) UploadFile(ctx context.Context, caminho string, dados []byte) error {
	fullPath := filepath.Join(s.BasePath, caminho)

	// Criar diretórios se não existirem
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
	}

	// Escrever arquivo
	if err := os.WriteFile(fullPath, dados, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo %s: %w", fullPath, err)
	}

	fmt.Printf("Arquivo enviado com sucesso: %s (%d bytes)\n", fullPath, len(dados))
	return nil
}

// Garantir que implementa a interface
var _ domain.StorageHandler = (*DefaultStorageService)(nil)

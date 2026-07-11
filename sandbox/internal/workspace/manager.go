package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// Manager는 호스트 파일 시스템의 작업 공간(Workspace) 관리 및 파일 CRUD를 담당합니다.
// Docker SDK 의존성을 전혀 가지지 않는 순수한 호스트 FS 계층입니다.
type Manager struct {
	baseDir string
}

// NewManager는 지정된 기본 디렉토리를 바탕으로 Workspace Manager를 생성합니다.
func NewManager(baseDir string) *Manager {
	return &Manager{
		baseDir: baseDir,
	}
}

// GetPath는 지정된 ID에 해당하는 샌드박스의 호스트 내 절대 경로를 반환합니다.
func (m *Manager) GetPath(id string) string {
	abs, err := filepath.Abs(filepath.Join(m.baseDir, id))
	if err != nil {
		return filepath.Join(m.baseDir, id)
	}
	return abs
}

// CreateWorkspace는 새로운 샌드박스를 위한 호스트 디렉토리를 생성합니다.
func (m *Manager) CreateWorkspace(id string) error {
	path := m.GetPath(id)
	return os.MkdirAll(path, 0755)
}

// RemoveWorkspace는 샌드박스의 호스트 디렉토리와 내부 모든 파일을 재귀적으로 제거합니다.
func (m *Manager) RemoveWorkspace(id string) error {
	path := m.GetPath(id)
	return os.RemoveAll(path)
}

// ReadFile은 지정된 샌드박스 작업 공간 내부의 특정 파일 내용을 읽어 반환합니다.
func (m *Manager) ReadFile(id string, filename string) ([]byte, error) {
	path := filepath.Join(m.GetPath(id), filename)
	return os.ReadFile(path)
}

// WriteFile은 지정된 샌드박스 작업 공간 내부에 파일을 생성하거나 갱신합니다.
func (m *Manager) WriteFile(id string, filename string, content []byte) error {
	path := filepath.Join(m.GetPath(id), filename)
	// 디렉토리 구조가 필요할 경우 생성
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return os.WriteFile(path, content, 0644)
}

// DeleteFile은 지정된 샌드박스 작업 공간 내부의 특정 파일을 삭제합니다.
func (m *Manager) DeleteFile(id string, filename string) error {
	path := filepath.Join(m.GetPath(id), filename)
	return os.Remove(path)
}

// ListFiles는 샌드박스 작업 공간 내부의 모든 파일 목록을 상대 경로로 조회합니다.
func (m *Manager) ListFiles(id string) ([]string, error) {
	path := m.GetPath(id)
	var files []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(path, filePath)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk files under %s: %w", path, err)
	}
	return files, nil
}

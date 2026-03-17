package pipeline

import (
    "fmt"
    "lab1/internal/domain"
    "strings"
    "sync"
)

// Pipeline этапы обработки данных
type Pipeline struct {
    errorChan chan error
}

// NewPipeline создает новый pipeline
func NewPipeline() *Pipeline {
    return &Pipeline{
        errorChan: make(chan error, 100),
    }
}

// Stage1 - первый этап: фильтрация и валидация постов
func (p *Pipeline) Stage1(input <-chan domain.Post) <-chan domain.Post {
    output := make(chan domain.Post)

    go func() {
        defer close(output)
        
        for post := range input {
            // Валидация: пропускаем только посты с непустым заголовком
            if strings.TrimSpace(post.Title) == "" {
                p.errorChan <- fmt.Errorf("пост #%d имеет пустой заголовок", post.ID)
                continue
            }
            
            // Фильтрация: оставляем только посты с ID < 20 (для примера)
            if post.ID > 20 {
                continue
            }
            
            output <- post
        }
    }()

    return output
}

// Stage2 - второй этап fan-out: обработка каждого поста в отдельной горутине
func (p *Pipeline) Stage2(input <-chan domain.Post) <-chan domain.FormattedPost {
    output := make(chan domain.FormattedPost)
    
    // Fan-out: запускаем workers
    const numWorkers = 3
    var wg sync.WaitGroup
    
    // Создаем канал для результатов от workers
    results := make(chan domain.FormattedPost, numWorkers)
    
    // Запускаем workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go p.worker(input, results, &wg, i)
    }
    
    // Fan-in: собираем результаты
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Перенаправляем результаты в выходной канал
    go func() {
        for result := range results {
            output <- result
        }
        close(output)
    }()
    
    return output
}

// worker обрабатывает посты
func (p *Pipeline) worker(input <-chan domain.Post, results chan<- domain.FormattedPost, wg *sync.WaitGroup, workerID int) {
    defer wg.Done()
    
    for post := range input {
        // Форматирование данных
        formatted := domain.FormattedPost{
            PostID:      post.ID,
            UserID:      post.UserID,
            ShortTitle:  p.truncateString(post.Title, 30),
            PreviewBody: p.truncateString(post.Body, 50),
        }
        
        results <- formatted
    }
}

// truncateString обрезает строку до указанной длины
func (p *Pipeline) truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "..."
}

// Run запускает pipeline
func (p *Pipeline) Run(posts []domain.Post, limit int) ([]domain.FormattedPost, []error) {
    // Создаем входной канал
    input := make(chan domain.Post)
    
    // Запускаем pipeline
    stage1Out := p.Stage1(input)
    stage2Out := p.Stage2(stage1Out)
    
    // Загружаем данные в pipeline
    go func() {
        defer close(input)
        
        // Ограничиваем количество постов
        for i, post := range posts {
            if limit > 0 && i >= limit {
                break
            }
            input <- post
        }
    }()
    
    // Собираем результаты
    var results []domain.FormattedPost
    var errors []error
    
    // Собираем результаты из pipeline
    for result := range stage2Out {
        results = append(results, result)
    }
    
    // Собираем ошибки
    close(p.errorChan)
    for err := range p.errorChan {
        errors = append(errors, err)
    }
    
    return results, errors
}
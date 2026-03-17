package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "lab1/internal/client"
    "lab1/internal/pipeline"
)

func main() {
    // Определяем флаги командной строки
    var (
        limit = flag.Int("limit", 10, "количество постов для обработки (по умолчанию 10)")
        format = flag.String("format", "simple", "формат вывода: simple или detailed (по умолчанию simple)")
    )
    flag.Parse()

    // Создаем клиент
    apiClient := client.NewJSONPlaceholderClient()

    // Получаем данные
    fmt.Println("Загрузка данных из API...")
    posts, err := apiClient.FetchPosts()
    if err != nil {
        log.Fatalf("Ошибка при получении данных: %v", err)
    }

    fmt.Printf("Загружено %d постов\n\n", len(posts))

    // Создаем и запускаем pipeline
    p := pipeline.NewPipeline()
    results, errors := p.Run(posts, *limit)

    // Выводим результаты
    fmt.Printf("Обработано %d постов\n", len(results))
    fmt.Printf("Ошибок: %d\n\n", len(errors))

    // Выводим отформатированные результаты
    fmt.Println("РЕЗУЛЬТАТЫ:")
    for i, result := range results {
        if *format == "detailed" {
            fmt.Printf("%d. Post #%d (User: %d)\n   Title: %s\n   Preview: %s\n\n",
                i+1, result.PostID, result.UserID, result.ShortTitle, result.PreviewBody)
        } else {
            fmt.Printf("%d. %s\n", i+1, result.ShortTitle)
        }
    }

    // Выводим ошибки если были
    if len(errors) > 0 {
        fmt.Println("\nОШИБКИ:")
        for i, err := range errors {
            fmt.Printf("%d. %v\n", i+1, err)
        }
        os.Exit(1)
    }
}
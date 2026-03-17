package client

import (
    "encoding/json"
    "fmt"
    "net/http"
    "lab1/internal/domain"
)

// APIClient интерфейс для работы с API
type APIClient interface {
    FetchPosts() ([]domain.Post, error)
    FetchPost(id int) (*domain.Post, error)
}

// JSONPlaceholderClient реализация клиента для JSONPlaceholder API
type JSONPlaceholderClient struct {
    baseURL string
    httpClient *http.Client
}

// NewJSONPlaceholderClient создает новый экземпляр клиента
func NewJSONPlaceholderClient() *JSONPlaceholderClient {
    return &JSONPlaceholderClient{
        baseURL: "https://jsonplaceholder.typicode.com",
        httpClient: &http.Client{},
    }
}

// FetchPosts получает все посты
func (c *JSONPlaceholderClient) FetchPosts() ([]domain.Post, error) {
    resp, err := c.httpClient.Get(c.baseURL + "/posts")
    if err != nil {
        return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
    }

    var posts []domain.Post
    if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
        return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
    }

    return posts, nil
}

// FetchPost получает пост по ID
func (c *JSONPlaceholderClient) FetchPost(id int) (*domain.Post, error) {
    resp, err := c.httpClient.Get(fmt.Sprintf("%s/posts/%d", c.baseURL, id))
    if err != nil {
        return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
    }

    var post domain.Post
    if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
        return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
    }

    return &post, nil
}
package domain

import "fmt"

// Post представляет модель данных из JSONPlaceholder
type Post struct {
    UserID int    `json:"userId"`
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

// FormattedPost представляет отформатированную версию поста
type FormattedPost struct {
    PostID      int
    ShortTitle  string
    PreviewBody string
    UserID      int
}

// String возвращает строковое представление форматированного поста
func (fp FormattedPost) String() string {
    return fmt.Sprintf("Post #%d (User: %d)\nTitle: %s\nPreview: %s\n---",
        fp.PostID, fp.UserID, fp.ShortTitle, fp.PreviewBody)
}
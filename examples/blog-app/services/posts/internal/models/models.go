package models

import (
	"time"
)

// Posts represents a posts in the system
type Posts struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	// Add your fields here
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePostsRequest represents the request to create a new posts
type CreatePostsRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	// Add your fields here
}

// UpdatePostsRequest represents the request to update a posts
type UpdatePostsRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Add your fields here
}

// PostsResponse represents the API response for a posts
type PostsResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	// Add your fields here
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListPostsResponse represents the API response for listing postss
type ListPostsResponse struct {
	Items      []PostsResponse `json:"items"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
}

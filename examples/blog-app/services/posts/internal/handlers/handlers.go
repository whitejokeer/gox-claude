package posts

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
)

// Handlers wraps the service with HTTP handlers
type Handlers struct {
    service *Service
}

// NewHandlers creates new HTTP handlers
func NewHandlers(service *Service) *Handlers {
    return &Handlers{service: service}
}

// HandleGet handles GET /posts/:id
func (h *Handlers) HandleGet(c *gin.Context) {
    id := c.Param("id")
    
    item, err := h.service.Get(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, item)
}

// HandleList handles GET /posts
func (h *Handlers) HandleList(c *gin.Context) {
    items, err := h.service.List(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, items)
}

// HandleCreate handles POST /posts
func (h *Handlers) HandleCreate(c *gin.Context) {
    var item Posts
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.service.Create(c.Request.Context(), &item); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, item)
}

// HandleUpdate handles PUT /posts/:id
func (h *Handlers) HandleUpdate(c *gin.Context) {
    id := c.Param("id")
    
    var item Posts
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.service.Update(c.Request.Context(), id, &item); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.Status(http.StatusOK)
}

// HandleDelete handles DELETE /posts/:id
func (h *Handlers) HandleDelete(c *gin.Context) {
    id := c.Param("id")
    
    if err := h.service.Delete(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.Status(http.StatusNoContent)
}
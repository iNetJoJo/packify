package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"packify/internal/models"
	"packify/internal/services"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom renderer for Echo
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Map page template names to their corresponding content templates
	contentTemplateMap := map[string]string{
		"home.html":               "content",
		"pack_sizes.html":         "content",
		"calculation_result.html": "calculation_result",
		"pack_sizes_table.html":   "pack_sizes_table",
	}

	// If this is a page template, render the content template directly
	if contentTemplate, ok := contentTemplateMap[name]; ok {
		// For partial templates, render them directly
		//TODO make this more generic in case we add more partials
		if name == "calculation_result.html" || name == "pack_sizes_table.html" {
			return t.templates.ExecuteTemplate(w, contentTemplate, data)
		}

		// For full page templates, we need to ensure that only the specific page template is loaded
		// Create a new template set for this render operation
		tmpl, err := template.New("").ParseGlob("templates/layouts/*.html")
		if err != nil {
			return err
		}

		// Parse only the specific page template
		pagePath := "templates/pages/" + name
		_, err = tmpl.ParseFiles(pagePath)
		if err != nil {
			return err
		}

		funcMap := template.FuncMap{
			"multiply": func(a, b uint64) uint64 {
				return a * b
			},
		}

		// Parse any partials that might be needed
		_, err = tmpl.Funcs(funcMap).ParseGlob("templates/partials/*.html")
		if err != nil {
			return err
		}

		// Now render the base layout with the specific page template
		return tmpl.ExecuteTemplate(w, "layouts/base.html", data)
	}

	// For other templates, render them directly
	return t.templates.ExecuteTemplate(w, name, data)
}

// NewTemplateRenderer creates a new template renderer
func NewTemplateRenderer() (*TemplateRenderer, error) {
	// Define template functions
	funcMap := template.FuncMap{
		"multiply": func(a, b uint64) uint64 {
			return a * b
		},
	}

	// Parse templates
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/**/*.html")
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{
		templates: tmpl,
	}, nil
}

// Handler contains all the HTTP handlers
type Handler struct {
	PackService *services.PackService
	Renderer    *TemplateRenderer
}

// NewHandler creates a new handler
func NewHandler(packService *services.PackService, renderer *TemplateRenderer) *Handler {
	return &Handler{
		PackService: packService,
		Renderer:    renderer,
	}
}

// RegisterRoutes registers all the routes
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	// API routes
	api := e.Group("/api")
	{
		// Pack calculation routes
		api.POST("/calculate", h.CalculatePacks)

		// Pack size management routes
		api.GET("/pack-sizes", h.GetPackSizes)
		api.POST("/pack-sizes", h.AddPackSize)
		api.PUT("/pack-sizes/:id", h.UpdatePackSize)
		api.DELETE("/pack-sizes/:id", h.DeletePackSize)
	}

	// Web UI routes
	e.GET("/", h.HomePage)
	e.GET("/calculate", h.CalculatePage)
	e.POST("/calculate", h.CalculatePagePost)
	e.GET("/pack-sizes", h.PackSizesPage)

	// Partial templates for HTMX
	e.GET("/pack-sizes/partial", h.PackSizesPartial)

	// Static files
	e.Static("/static", "static")
	e.Static("/css", "static/css")
}

type CalculateRequest struct {
	// Handle potential overflow for max uint64 values. If the value exceeds the maximum, it will wrap around to zero.
	// Consider using big.Int for handling extremely large numbers if required in the future.
	// Lets be real here, we are not going to have more than 2^64 items in a single order but then again, we are not here to judge.
	ItemsOrdered uint64 `json:"itemsOrdered"`
}

// CalculatePacks calculates the optimal packs for an order
func (h *Handler) CalculatePacks(c echo.Context) error {
	// Parse request

	req := new(CalculateRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request"))
	}

	// Validate request
	if req.ItemsOrdered <= 0 {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Items ordered must be greater than 0"))
	}

	// Calculate packs
	result, err := h.PackService.CalculatePacks(req.ItemsOrdered)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
	}

	response := CalculateResponse{
		TotalPacks:  result.TotalPacks,
		TotalItems:  result.TotalItems,
		ExcessItems: result.ExcessItems,
	}

	// Convert map to slice for better JSON formatting
	for size, count := range result.PackCounts {
		response.Packs = append(response.Packs, PackInfo{
			Size:  size,
			Count: count,
		})
	}

	return c.JSON(http.StatusOK, response)
}

// PackInfo Format response
type PackInfo struct {
	Size  uint64 `json:"size"`
	Count uint64 `json:"count"`
}
type CalculateResponse struct {
	Packs       []PackInfo `json:"packs"`
	TotalPacks  uint64     `json:"totalPacks"`
	TotalItems  uint64     `json:"totalItems"`
	ExcessItems uint64     `json:"excessItems"`
}

// GetPackSizes returns all pack sizes
func (h *Handler) GetPackSizes(c echo.Context) error {
	packSizes, err := h.PackService.GetPackSizes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, packSizes)
}

type AddPackSizeRequest struct {
	Size uint64 `form:"size" json:"size"`
}

// AddPackSize adds a new pack size
func (h *Handler) AddPackSize(c echo.Context) error {
	// Parse request

	req := new(AddPackSizeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request"))
	}

	// Validate request
	if req.Size <= 0 {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Pack size must be greater than 0"))
	}

	// Add pack size
	if err := h.PackService.AddPackSize(req.Size); err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, models.NewSuccessResponse("Pack size added successfully"))
}

type UpdatePackSizeRequest struct {
	IsAvailable bool `json:"isAvailable"`
}

// UpdatePackSize updates a pack size
func (h *Handler) UpdatePackSize(c echo.Context) error {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid ID"))
	}

	// Parse request

	req := new(UpdatePackSizeRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid request"))
	}

	// Update pack size
	if err := h.PackService.UpdatePackSize(uint(id), req.IsAvailable); err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, models.NewSuccessResponse("Pack size updated successfully"))
}

// DeletePackSize deletes a pack size
func (h *Handler) DeletePackSize(c echo.Context) error {
	// Parse ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.NewErrorResponse("Invalid ID"))
	}

	// Delete pack size
	if err := h.PackService.DeletePackSize(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, models.NewSuccessResponse("Pack size deleted successfully"))
}

// HomePage renders the home page
func (h *Handler) HomePage(c echo.Context) error {
	fmt.Println("Rendering home page")

	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"Title": "Home",
	})
}

// CalculatePage renders the calculate page
func (h *Handler) CalculatePage(c echo.Context) error {
	return c.Render(http.StatusOK, "calculate.html", map[string]interface{}{
		"Title":        "Calculate Packs",
		"ItemsOrdered": 1, // Default value
	})
}

type CalculatePagePostRequest struct {
	ItemsOrdered uint64 `form:"itemsOrdered" json:"itemsOrdered"`
}

// CalculatePagePost handles the calculate form submission
func (h *Handler) CalculatePagePost(c echo.Context) error {
	// Parse request

	req := new(CalculatePagePostRequest)

	if err := c.Bind(req); err != nil {
		return c.Render(http.StatusBadRequest, "calculation_result.html", map[string]interface{}{
			"Error": "Invalid request",
		})
	}

	if req.ItemsOrdered <= 0 {
		return c.Render(http.StatusBadRequest, "calculation_result.html", map[string]interface{}{
			"Error": "Items ordered must be a positive number",
		})
	}

	// Calculate packs
	result, err := h.PackService.CalculatePacks(req.ItemsOrdered)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "calculation_result.html", map[string]interface{}{
			"Error": err.Error(),
		})
	}

	var packs []PackInfo
	for size, count := range result.PackCounts {
		packs = append(packs, PackInfo{
			Size:  size,
			Count: count,
		})
	}

	// If this is an HTMX request, render just the result partial
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "calculation_result.html", map[string]interface{}{
			"ItemsOrdered": req.ItemsOrdered,
			"Result": CalculateResponse{
				Packs:       packs,
				TotalPacks:  result.TotalPacks,
				TotalItems:  result.TotalItems,
				ExcessItems: result.ExcessItems,
			},
		})
	}

	// Otherwise render the full page
	return c.Render(http.StatusOK, "calculate.html", map[string]interface{}{
		"Title":        "Calculate Packs",
		"ItemsOrdered": req.ItemsOrdered,
		"Result": CalculateResponse{
			Packs:       packs,
			TotalPacks:  result.TotalPacks,
			TotalItems:  result.TotalItems,
			ExcessItems: result.ExcessItems,
		},
	})
}

// PackSizesPage renders the pack sizes management page
func (h *Handler) PackSizesPage(c echo.Context) error {
	return c.Render(http.StatusOK, "pack_sizes.html", map[string]interface{}{
		"Title": "Manage Pack Sizes",
	})
}

// PackSizesPartial renders the pack sizes table partial
func (h *Handler) PackSizesPartial(c echo.Context) error {
	packSizes, err := h.PackService.GetPackSizes()
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "<div class='error'>Failed to load pack sizes</div>")
	}

	return c.Render(http.StatusOK, "pack_sizes_table.html", map[string]interface{}{
		"PackSizes": packSizes,
	})
}

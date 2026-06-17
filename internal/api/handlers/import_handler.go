package handlers

import (
	"fmt"
	"net/http"

	apimw "go-nagiosql/internal/api/middleware"
	"go-nagiosql/internal/services/logbook"
	"go-nagiosql/internal/services/nagimport"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// ImportHandler handles /api/v1/import endpoints.
type ImportHandler struct{ db *gorm.DB }

// NewImportHandler creates an ImportHandler.
func NewImportHandler(db *gorm.DB) *ImportHandler { return &ImportHandler{db: db} }

// importRequest is the JSON body for POST /api/v1/import.
type importRequest struct {
	File      string `json:"file"`
	ConfigID  uint8  `json:"config_id"`
	Overwrite bool   `json:"overwrite"`
}

// ImportResult summarises an import operation.
type ImportResult struct {
	Inserted int      `json:"inserted"`
	Updated  int      `json:"updated"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors"`
}

// Import godoc
// @Summary      Import a Nagios .cfg file into the database
// @Tags         import
// @Accept       json
// @Produce      json
// @Param        body  body  importRequest  true  "Import parameters"
// @Security     BearerAuth
// @Success      200  {object}  ImportResult
// @Router       /import [post]
func (h *ImportHandler) Import(c *echo.Context) error {
	var req importRequest
	if err := c.Bind(&req); err != nil {
		return BadRequest(c, "invalid body")
	}
	if req.File == "" {
		return BadRequest(c, "file is required")
	}

	objects, err := nagimport.ParseFile(req.File)
	if err != nil {
		return BadRequest(c, fmt.Sprintf("parsing file: %v", err))
	}

	result := ImportResult{}
	for _, obj := range objects {
		ok, wasNew, importErr := nagimport.ImportObject(h.db, obj, req.ConfigID, req.Overwrite)
		if importErr != nil {
			result.Errors = append(result.Errors, importErr.Error())
			continue
		}
		switch {
		case !ok:
			result.Skipped++
		case wasNew:
			result.Inserted++
		default:
			result.Updated++
		}
	}

	claims := apimw.ClaimsFromContext(c)
	logbook.Write(h.db, claims.DomainID, claims.Username, "import", "cfg", req.File,
		fmt.Sprintf("inserted=%d updated=%d skipped=%d errors=%d",
			result.Inserted, result.Updated, result.Skipped, len(result.Errors)))
	return c.JSON(http.StatusOK, result)
}

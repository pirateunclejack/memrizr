package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pirateXXOO/memrizr/account/model"
	"github.com/pirateXXOO/memrizr/account/model/apperrors"
)

// Image handler
func (h *Handler) Image(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)

	// limit overlay large request bodies
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	imageFileHeader, err := c.FormFile("imageFile")

	// check for error before checking for non-nil header
	if err != nil {
		// should be validation error
		log.Printf("Unable parse multipart/form-data: %+v", err)

		if err.Error() == "http: request body too large" {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": fmt.Sprintf("Max request body size is %v bytes\n", h.MaxBodyBytes),
			})
			return
		}
		e := apperrors.NewBadRequest("Unable to parse multipart/form-data")
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	mimeType := imageFileHeader.Header.Get("Content-Type")

	// Validate image mime-type is allowable
	if valid := IsAllowedImageType(mimeType); !valid {
		log.Println("Image is not an allowable mime-type")
		e := apperrors.NewBadRequest("imageFile must be 'image/jpeg' or 'image/png'")
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	ctx := c.Request.Context()

	updatedUser, err := h.UserService.SetProfileImage(ctx, authUser.UID, imageFileHeader)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imageUrl": updatedUser.ImageURL,
		"message":  "success",
	})
}

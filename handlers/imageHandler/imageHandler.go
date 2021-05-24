package imageHandler

import (
	"encoding/json"
	"fmt"
	"google-cloud-task-processor/dto"
	"google-cloud-task-processor/handlers/utilities"
	"google-cloud-task-processor/services/imageService"
	"net/http"
)

type ImageHandler struct {
	imageService *imageService.ImageService
}

func New(imageService *imageService.ImageService) *ImageHandler {
	return &ImageHandler{imageService: imageService}
}

func (h *ImageHandler) UploadImage(writer http.ResponseWriter, request *http.Request) {
	var image = &dto.ImageRequest{}
	if err := json.NewDecoder(request.Body).Decode(image); err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusBadRequest, fmt.Errorf("json decode failed"))
		return
	}
	response, err := h.imageService.UploadImage(image)

	if err != nil {
		utilities.ErrorJsonRespond(writer, http.StatusInternalServerError, err)
		return
	}
	utilities.RespondJson(writer, http.StatusCreated, response)
}
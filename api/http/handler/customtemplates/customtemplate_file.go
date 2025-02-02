package customtemplates

import (
	"net/http"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

type fileResponse struct {
	FileContent string
}

// @id CustomTemplateFile
// @summary Get Template stack file content.
// @description Retrieve the content of the Stack file for the specified custom template
// @description **Access policy**: authenticated
// @tags custom_templates
// @security ApiKeyAuth
// @security jwt
// @produce json
// @param id path int true "Template identifier"
// @success 200 {object} fileResponse "Success"
// @failure 400 "Invalid request"
// @failure 404 "Custom template not found"
// @failure 500 "Server error"
// @router /custom_templates/{id}/file [get]
func (handler *Handler) customTemplateFile(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	customTemplateID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid custom template identifier route variable", err)
	}

	customTemplate, err := handler.DataStore.CustomTemplate().CustomTemplate(portainer.CustomTemplateID(customTemplateID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a custom template with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find a custom template with the specified identifier inside the database", err)
	}

	entryPath := customTemplate.EntryPoint
	if customTemplate.GitConfig != nil {
		entryPath = customTemplate.GitConfig.ConfigFilePath
	}
	fileContent, err := handler.FileService.GetFileContent(customTemplate.ProjectPath, entryPath)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve custom template file from disk", err)
	}

	return response.JSON(w, &fileResponse{FileContent: string(fileContent)})
}

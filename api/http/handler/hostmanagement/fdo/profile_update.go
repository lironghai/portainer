package fdo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
	portainer "github.com/portainer/portainer/api"
)

// @id updateProfile
// @summary updates an existing FDO Profile
// @description updates an existing FDO Profile
// @description **Access policy**: administrator
// @tags intel
// @security jwt
// @produce json
// @param id path int true "FDO Profile identifier"
// @success 200 "Success"
// @failure 400 "Invalid request"
// @failure 409 "Profile name already exists"
// @failure 500 "Server error"
// @router /fdo/profiles/{id} [put]
func (handler *Handler) updateProfile(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	id, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Bad request", errors.New("missing 'id' query parameter"))
	}

	var payload createProfileFromFileContentPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return httperror.BadRequest("Invalid request payload", err)
	}

	profile, err := handler.DataStore.FDOProfile().FDOProfile(portainer.FDOProfileID(id))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a FDO Profile with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find a FDO Profile with the specified identifier inside the database", err)
	}

	isUnique, err := handler.checkUniqueProfileName(payload.Name, id)
	if err != nil {
		return httperror.InternalServerError(err.Error(), err)
	}
	if !isUnique {
		return &httperror.HandlerError{StatusCode: http.StatusConflict, Message: fmt.Sprintf("A profile with the name '%s' already exists", payload.Name), Err: errors.New("a profile already exists with this name")}
	}

	filePath, err := handler.FileService.StoreFDOProfileFileFromBytes(strconv.Itoa(int(profile.ID)), []byte(payload.ProfileFileContent))
	if err != nil {
		return httperror.InternalServerError("Unable to update profile", err)
	}
	profile.FilePath = filePath
	profile.Name = payload.Name

	err = handler.DataStore.FDOProfile().Update(profile.ID, profile)
	if err != nil {
		return httperror.InternalServerError("Unable to update profile", err)
	}

	return response.JSON(w, profile)
}

package users_transport_http

import (
	"net/http"

	core_logger "github.com/koliashka/golang-todoapp/internal/core/logger"
	core_http_response "github.com/koliashka/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/koliashka/golang-todoapp/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userId from path",
		)
	}

	user, err := h.usersService.GetUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userId",
		)
	}

	response := GetUserResponse(UserDtoFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)
}

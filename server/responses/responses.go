package responses

type ResponseBody struct {
	Data  interface{}   `json:"data"`
	Error ResponseError `json:"error"`
}

func GenerateResponse(body interface{}, err *ErrorType) (int, ResponseBody) {
	if body == nil {
		body = struct{}{}
	}
	responseBody := ResponseBody{Data: body}

	code := 200
	if err != nil {
		code = err.Code
		responseBody.Error.Type = err.Type
		responseBody.Error.Message = err.Message
		responseBody.Error.Details = err.Details
	}

	return code, responseBody
}

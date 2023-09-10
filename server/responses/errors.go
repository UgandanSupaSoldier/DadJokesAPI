package responses

type ResponseError struct {
	Type    string  `json:"type,omitempty"`
	Message string  `json:"message,omitempty"`
	Details *string `json:"details,omitempty"`
}

type ErrorType struct {
	Code int `json:"-"`
	ResponseError
}

func (e *ErrorType) Error() string {
	return e.Type
}

var (
	InvalidJSONError  = &ErrorType{400, ResponseError{"invalid_json", "The request body contains incorrectly formatted JSON data", nil}}
	InvalidQueryError = &ErrorType{400, ResponseError{"invalid_query", "The request query string is missing a required option", nil}}

	AuthMissingError  = &ErrorType{401, ResponseError{"unauthorized", "The request is missing an authorization header", nil}}
	AuthExpiredError  = &ErrorType{401, ResponseError{"unauthorized", "The authorization token provided has expired", nil}}
	AuthInvalidError  = &ErrorType{401, ResponseError{"unauthorized", "The authorization header is malformed or invalid", nil}}
	UnauthorizedError = &ErrorType{401, ResponseError{"unauthorized", "You are not authorized to access this resource", nil}}

	NotFoundError        = &ErrorType{404, ResponseError{"not_found", "The requested resource or object could not be found", nil}}
	TooManyRequestsError = &ErrorType{429, ResponseError{"too_many_requests", "Too many requests have been made to this resource", nil}}
)

func ErrorWithDetails(err *ErrorType, details string) *ErrorType {
	err.Details = &details
	return err
}

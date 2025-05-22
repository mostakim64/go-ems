package msgutil

type Data map[string]interface{}

type Msg struct {
	Data Data
}

func NewMessage() Msg {
	return Msg{
		Data: make(Data),
	}
}

func (m Msg) Set(key string, value interface{}) Msg {
	m.Data[key] = value
	return m
}

func (m Msg) Done() Data {
	return m.Data
}

func RequestBodyParseErrorMsg() Data {
	return NewMessage().Set("message", "Failed to parse request body").Done()
}

func JwtCreateErrorMsg() Data {
	return NewMessage().Set("message", "Failed to create JWT token").Done()
}

func SomethingWentWrongMsg() Data {
	return NewMessage().Set("message", "Something went wrong").Done()
}

func ExpectationFailedMsg() Data {
	return NewMessage().Set("message", "Expectation failed").Done()
}

func AccessForbiddenMsg() Data {
	return NewMessage().Set("message", "Access forbidden").Done()
}

func UnprocessableEntityMsg() Data {
	return NewMessage().Set("message", "Unprocessable entity").Done()
}

func InvalidRequestMsg() Data {
	return NewMessage().Set("message", "Invalid Request").Done()
}

func PermissionError() Data {
	return NewMessage().Set("message", "Operation not permitted").Done()
}

func RefreshTokenNotFound() Data {
	return NewMessage().Set("message", "refresh token not found").Done()
}

func EventNotFound() Data {
	return NewMessage().Set("message", "Event not found").Done()
}

func UserAlreadyExists() Data {
	return NewMessage().Set("message", "User already exists").Done()
}

func UserNotFound() Data {
	return NewMessage().Set("message", "User not found").Done()
}

func InvalidLoginCredentials() Data {
	return NewMessage().Set("message", "Invalid login credentials").Done()
}

// TODO: fix msg
func UserUnauthorized() Data {
	return NewMessage().Set("message", "User unauthorized").Done()
}

func UserUpdatedSuccessfully() Data {
	return NewMessage().Set("message", "User updated successfully").Done()
}

func UserDeletedSuccessfully() Data {
	return NewMessage().Set("message", "User deleted successfully").Done()
}

package errutil

import (
	"errors"
)

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrInvalidInput            = errors.New("invalid input")
	ErrUserAlreadyExist        = errors.New("user already exists")
	ErrInvalidLoginCredentials = errors.New("invalid login credentials")
	ErrCreateJwt               = errors.New("failed to create JWT token")
	ErrAccessTokenSign         = errors.New("failed to sign access_token")
	ErrRefreshTokenSign        = errors.New("failed to sign refresh_token")

	ErrInvalidEmail              = errors.New("invalid email")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrStoreTokenUuid            = errors.New("failed to store token uuid")
	ErrUpdateMetaData            = errors.New("failed to update metadata")
	ErrNoContextUser             = errors.New("failed to get user from context")
	ErrInvalidRefreshToken       = errors.New("invalid refresh_token")
	ErrInvalidAccessToken        = errors.New("invalid access_token")
	ErrInvalidPasswordResetToken = errors.New("invalid reset_token")
	ErrInvalidVerifyEmailToken   = errors.New("invalid verification token")
	ErrInvalidVToken             = errors.New("invalid reset_token")
	ErrInvalidRefreshUuid        = errors.New("invalid refresh_uuid")
	ErrInvalidAccessUuid         = errors.New("invalid refresh_uuid")
	ErrInvalidJwtSigningMethod   = errors.New("invalid signing method while parsing jwt")
	ErrParseJwt                  = errors.New("failed to parse JWT token")
	ErrDeleteOldTokenUuid        = errors.New("failed to delete old token uuids")
	ErrSendingEmail              = errors.New("failed to send email")
	ErrUserCreate                = errors.New("failed to create user")
	ErrUserNotFound              = errors.New("user not found")

	ErrUserUpdate                       = errors.New("failed to update user")
	ErrInvalidOtpNonce                  = errors.New("invalid otp nonce")
	ErrInvalidOtp                       = errors.New("invalid otp")
	ErrInvalidAuNumber                  = errors.New("not an australian number")
	ErrInvalidPhoneNumber               = errors.New("invalid phone number")
	ErrEmailAlreadyInUse                = errors.New("email already in use")
	ErrPhoneAlreadyInUse                = errors.New("phone already in use")
	ErrEmailUpdateNotAllowed            = errors.New("email update not allowed")
	ErrPhoneUpdateNotAllowed            = errors.New("phone update not allowed")
	ErrNoPhoneNumberIsSet               = errors.New("no phone number is set")
	ErrNoEmailIsSet                     = errors.New("no email is set")
	ErrEmailAlreadyVerified             = errors.New("email already verified")
	ErrEmailMisMatched                  = errors.New("email mismatched")
	ErrPhoneAlreadyVerified             = errors.New("phone already verified")
	ErrPasswordUpdateNotAllowed         = errors.New("password update not allowed")
	ErrPhoneNumberUpdateNotAllowed      = errors.New("phone number update not allowed")
	ErrInvalidLoginProvider             = errors.New("invalid login provider")
	ErrUserAlreadyRegistered            = errors.New("user already registered")
	ErrUserAlreadyRegisteredViaGoogle   = errors.New("user already registered via google")
	ErrUserAlreadyRegisteredViaFacebook = errors.New("user already registered via facebook")
	ErrUserAlreadyRegisteredViaApple    = errors.New("user already registered via apple")
	ErrLoginAttemptWithAppleProvider    = errors.New("invalid login attempt")
	ErrLoginAttemptWithGoogleProvider   = errors.New("invalid login attempt")
	ErrLoginAttemptWithFacebookProvider = errors.New("invalid login attempt")
	ErrInvalidAuthorizationToken        = errors.New("invalid authorization token")
	ErrInvalidPasswordFormat            = errors.New("minimum 8 characters with at least 1 uppercase letter(A-Z), 1 lowercase letter(a-z), 1 number(0-9) and 1 special character(.!@#~$%^&*()+|_<>)")
	ErrInvalidRegion                    = errors.New("invalid region")
	ErrInvalidLineLoginCountry          = errors.New("invalid login country")
)

func Exists(err error, errs []error) bool {
	for _, e := range errs {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}

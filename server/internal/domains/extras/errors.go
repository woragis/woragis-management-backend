package extras

import "errors"

const (
    ErrCodeInvalidPayload    = 6010
    ErrCodeRepositoryFailure = 6011
    ErrCodeNotFound          = 6012
)

const (
    ErrNilExtra       = "extras: extra entity is nil"
    ErrEmptyID        = "extras: id cannot be empty"
    ErrEmptyUserID    = "extras: user id cannot be empty"
    ErrUnableToPersist = "extras: unable to persist record"
    ErrUnableToFetch  = "extras: unable to fetch records"
    ErrUnableToUpdate = "extras: unable to update record"
    ErrNotFound       = "extras: record not found"
)

type DomainError struct {
    Code    int
    Message string
}

func (e *DomainError) Error() string {
    return e.Message
}

func NewDomainError(code int, message string) *DomainError {
    return &DomainError{
        Code:    code,
        Message: message,
    }
}

func AsDomainError(err error) (*DomainError, bool) {
    var domainErr *DomainError
    if errors.As(err, &domainErr) {
        return domainErr, true
    }
    return nil, false
}

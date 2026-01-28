package certifications

import "errors"

const (
    ErrCodeInvalidPayload    = 6000
    ErrCodeRepositoryFailure = 6001
    ErrCodeNotFound          = 6002
)

const (
    ErrNilCertification  = "certifications: certification entity is nil"
    ErrEmptyID           = "certifications: id cannot be empty"
    ErrEmptyUserID       = "certifications: user id cannot be empty"
    ErrUnableToPersist   = "certifications: unable to persist record"
    ErrUnableToFetch     = "certifications: unable to fetch records"
    ErrUnableToUpdate    = "certifications: unable to update record"
    ErrNotFound          = "certifications: record not found"
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

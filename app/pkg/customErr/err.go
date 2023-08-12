package customerr

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomError struct {
	desc     string
	grpcCode uint32
	httpCode int
}

func NewCustomError(text string, gprcCode codes.Code, httpCode int) error {
	return &CustomError{
		desc:     text,
		grpcCode: uint32(gprcCode),
		httpCode: httpCode,
	}
}

func (e *CustomError) Error() string {
	return e.desc
}

func (e *CustomError) BuildGrpcError() error {
	return status.Error(codes.Code(e.grpcCode), e.desc)
}

func (e *CustomError) GetGrpcCode() uint32 {
	return e.grpcCode
}

func (e *CustomError) GetHttpCode() int {
	return e.httpCode
}

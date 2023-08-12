package utils

import (
	customerr "aifory-pay-admin-bot/pkg/customErr"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println(err)
	}
	// вопросы?
}

func GetMoscowLoc() *time.Location {
	return loc
}

func GetProxyHost(proxy string) (host string, err error) {
	proxySplited := strings.Split(strings.Replace(proxy, "@", ":", 1), ":")
	if len(proxySplited) != 4 {
		return host, errors.New("wrong proxy length")
	}
	return proxySplited[2], nil
}

func GetProxyHostPort(proxy string) (host string, err error) {
	proxySplited := strings.Split(strings.Replace(proxy, "@", ":", 1), ":")
	if len(proxySplited) != 4 {
		return host, errors.New("wrong proxy length")
	}
	return fmt.Sprintf("%s:%s", proxySplited[2], proxySplited[3]), nil
}

func MakeDayUnix(timeValue time.Time) int64 {
	return time.Date(
		timeValue.Year(),
		timeValue.Month(),
		timeValue.Day(),
		0,
		0,
		0,
		0,
		timeValue.Location(),
	).Unix()
}

func InSlice(target int, src []int) bool {
	for _, r := range src {
		if target == r {
			return true
		}
	}
	return false
}

func ParseUrlEncodedRequest(data string) (result map[string]string, err error) {
	values, err := url.ParseQuery(data)
	if err != nil {
		return
	}

	result = map[string]string{}
	for k, v := range values {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return
}

func FormatErr(err error) error {
	if customErr, ok := err.(*customerr.CustomError); ok {
		return customErr.BuildGrpcError()
	}
	return status.Error(codes.Internal, err.Error())
}

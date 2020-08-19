package util

import (
	"fmt"
	"github.com/getsentry/sentry-go"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		sentry.CaptureException(err)
	}
}

package life

import "fmt"

func Message(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func ErrorMessagef(err error, format string, args ...any) string {
	if err == nil {
		return ""
	}

	return fmt.Sprintf("%s: %s", err.Error(), fmt.Sprintf(format, args...))
}

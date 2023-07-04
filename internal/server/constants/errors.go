package constants

type Error uint8

const (
	NoError Error = iota
	InvalidRequestError
	NoMetricNameError
	InvalidMetricTypeError
	InvalidValueError
)

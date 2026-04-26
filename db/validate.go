package db

import "regexp"

var emailRegexp = regexp.MustCompile(`^[a-zA-Z][\w_\-\.]+\w@[a-zA-Z][\w_\-\.]+[a-zA-Z]\.[a-zA-Z]{2,9}$`)

const (
	MAX_NAME_LEN   = 48
	MAX_EMAIL_LEN  = 72
	MIN_SECRET_LEN = 8
	MAX_SECRET_LEN = 32
)

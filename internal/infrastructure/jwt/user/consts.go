package user

import "time"

const (
	idClaims    = "uid"
	emailClaims = "email"
	expClaims   = "exp"
	expDuration = time.Hour * 48
)

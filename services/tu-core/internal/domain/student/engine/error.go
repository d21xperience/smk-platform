package engine

import "errors"

var errEnrollmentNotFound = errors.New("enrollment not found")
var errInvalidGradeTransition = errors.New("invalid grade transition for move")

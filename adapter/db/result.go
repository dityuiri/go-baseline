package db

import "database/sql"

//go:generate mockgen -destination=mock/result.go -package=mock . IResult

type IResult interface {
	sql.Result
}

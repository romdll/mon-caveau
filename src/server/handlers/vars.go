package handlers

import (
	"moncaveau/utils"
)

type AdjustQuantityRequest struct {
	Change int `json:"change"`
}

var (
	logger = utils.CreateLogger("handlers")
)

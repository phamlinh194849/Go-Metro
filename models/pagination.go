package models

import "go-metro/consts"

// PaginationResponse represents paginated response structure
type PaginationResponse struct {
  Data       interface{} `json:"data"`
  Total      int64       `json:"total"`
  Page       int         `json:"page"`
  Limit      int         `json:"limit"`
  TotalPages int         `json:"total_pages"`
}

// UserQueryParams represents query parameters for user listing
type UserQueryParams struct {
  Page   int         `form:"page" binding:"omitempty,min=1"`
  Limit  int         `form:"limit" binding:"omitempty,min=1,max=100"`
  Role   consts.Role `form:"role" binding:"omitempty,min=1,max=3"`
  Name   string      `form:"name" binding:"omitempty"`
  Status string      `form:"status" binding:"omitempty,oneof=active inactive"`
}

package structs

import "time"

//
type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude float64 `json:"latitude"`
	Height float64 `json:"height,omitempty"`
}

type Patrol struct {
	//
	ID int64 `json:"id,omitempty"`
	Line []*Point `json:"line,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Memo string `json:"memo"`
	CreateAt time.Time `json:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at,omitempty"`
}

//
type GetPatrolListResponse struct {
	List []*Patrol `json:"list"`
	Count int64 `json:"count"`
}

// 提交
type PostRequest struct {
	RequestID string `json:"request_id,omitempty"`
	Memo string `json:"memo"`
}
//
type PageParam struct {
	Page int64
	PageSize int64
}

//
type GetPatrolListRequest struct {
	UserID int64
	PageParam
}
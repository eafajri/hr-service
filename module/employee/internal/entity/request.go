package entity

type SubmitAttendanceRequest struct {
	UserID       int64  `json:"user_id"`
	Date         string `json:"date"`
	CheckInTime  string `json:"check_in_time"`
	CheckOutTime string `json:"check_out_time"`
}

type SubmitOvertimeRequest struct {
	UserID    int64  `json:"user_id"`
	Date      string `json:"date"`
	Durations int64  `json:"durations"`
}

type SubmitReimbursementRequest struct {
	UserID      int64   `json:"user_id"`
	Date        string  `json:"date"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

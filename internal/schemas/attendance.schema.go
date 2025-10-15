package schemas

import "time"

type Attendance struct {
	ID         int64     `db:"id" json:"id"`
	UserId     int64     `db:"user_id" json:"user_id"`
	CheckInAt  time.Time `db:"check_in_at" json:"check_in_at"`
	CheckOutAt time.Time `db:"check_out_at" json:"check_out_at"`
	Latitude   string    `db:"latitude" json:"latitude"`
	Longitude  string    `db:"longitude" json:"longitude"`
	Status     string    `db:"status" json:"status"`
	Note       string    `db:"note" json:"note"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

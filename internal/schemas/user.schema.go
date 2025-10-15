package schemas

type User struct {
	ID         int64  `db:"id" json:"id"`
	EmployeeID string `db:"employee_id" json:"employee_id"`
	Email      string `db:"email" json:"email"`
	Password   string `db:"password_hash" json:"password"`
	Name       string `db:"name" json:"name"`
	Timezone   string `db:"timezone" json:"timezone"`
	Role       string `db:"role" json:"role"`
	IsCheckIn  bool   `db:"is_check_in" json:"is_check_in"`
}

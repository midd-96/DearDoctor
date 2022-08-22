package repo

import (
	"database/sql"
	"dearDoctor/model"
	"dearDoctor/utils"
	"log"
)

type AdminRepository interface {
	FindAdmin(username string) (model.AdminResponse, error)
	AddDept(department model.Departments) error
	UpdateApproveFee(approvel model.ApproveAndFee, emailid string) error
	ViewAllAppointments(pagenation utils.Filter, doc_id int, day string) ([]model.AppointmentByDoctor, utils.Metadata, error)
}

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) AdminRepository {
	return &adminRepo{
		db: db,
	}
}

func (c *adminRepo) ViewAllAppointments(pagenation utils.Filter, doc_id int, day string) ([]model.AppointmentByDoctor, utils.Metadata, error) {

	var allappointments []model.AppointmentByDoctor

	query := `SELECT 
				COUNT(*) OVER(),
				time_consult,
				payment_mode,
				payment_status,
				email FROM confirmeds c
				WHERE doctor_id = $1 AND day_consult = $2
				LIMIT $3 OFFSET $4;`

	rows, err := c.db.Query(query, doc_id, day, pagenation.Limit(), pagenation.Offset())
	log.Println(query)
	log.Println(doc_id)
	log.Println(day)
	if err != nil {
		log.Println(err)
		return nil, utils.Metadata{}, err
	}
	var totalRecords int

	defer rows.Close()

	for rows.Next() {
		var appointment model.AppointmentByDoctor

		err = rows.Scan(
			&totalRecords,
			&appointment.Time_consult,
			&appointment.Payment_mode,
			&appointment.Payment_status,
			&appointment.Email)

		if err != nil {
			return allappointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
		}
		allappointments = append(allappointments, appointment)

	}
	if err := rows.Err(); err != nil {
		return allappointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
	}
	log.Println(allappointments)
	log.Println(utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize))
	return allappointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), nil

}

func (c *adminRepo) FindAdmin(username string) (model.AdminResponse, error) {

	var admin model.AdminResponse

	query := `SELECT 
			id,
			username,
			password,
			role
			FROM admins WHERE username = $1;`

	err := c.db.QueryRow(query,
		username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Password,
		&admin.Role)

	return admin, err
}

func (c *adminRepo) UpdateApproveFee(approvel model.ApproveAndFee, emailid string) error {

	var query string

	var err error

	query = `
				UPDATE
				   doctors 
				SET
				   approvel = $1, fee=$2
				WHERE
				   email = $3 ;`

	err = c.db.QueryRow(query, approvel.Approve, approvel.Fee, emailid).Err()

	return err

}

func (c *adminRepo) AddDept(department model.Departments) error {

	query := `INSERT INTO
				departments (dep_id, name) 
				VALUES
				   (
				      $1, $2
				   );`

	err := c.db.QueryRow(
		query,
		department.Dep_Id,
		department.Name,
	).Err()
	return err
}

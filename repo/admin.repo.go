package repo

import (
	"database/sql"
	"dearDoctor/model"
	"dearDoctor/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type AdminRepository interface {
	CreateAdmin(admin model.Admin) error
	FindAdmin(username string) (model.AdminResponse, error)
	AddDept(department model.Departments) error
	UpdateApproveFee(approvel model.ApproveAndFee, emailid string) error
	ViewAllAppointments(pagenation utils.Filter, filters model.Filter) ([]model.AppointmentByDoctor, utils.Metadata, error)
	CalculatePayout(doc_Id int) (string, error)
	ViewSingleUser(user_Id int) (model.UserResponse, error)
	ViewSingleDoctor(doc_Id int) (model.DoctorResponse, error)
	ApprovePayout(email string) (float64, error)
}

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) AdminRepository {
	return &adminRepo{
		db: db,
	}
}

func (c *adminRepo) ApprovePayout(email string) (float64, error) {

	query := `UPDATE payouts
				SET approvel = $1,
					approved_time = $2,
					wallet = wallet - last_requested_amount
					WHERE username = $3 AND approvel = $4;`

	Approved_At, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	log.Println(query)
	log.Println(Approved_At)
	log.Println(email)

	err := c.db.QueryRow(query, true, Approved_At, email, false).Err()

	log.Println(err)

	if err != nil {
		return 0.0, errors.New("Approvel Failed")
	}

	var walletBalance float64

	query1 := `SELECT wallet FROM 
				payouts WHERE 
				username = $1;`

	log.Println(query1)

	err1 := c.db.QueryRow(query1, email).Scan(&walletBalance)

	if err1 != nil {
		return walletBalance, errors.New("Approved but failed to fetch new wallet balance")
	}

	return walletBalance, nil
}

func (c *adminRepo) CalculatePayout(doc_Id int) (string, error) {

	var count int
	query := `SELECT COUNT(*)
				FROM confirmeds
				WHERE doctor_id = $1 AND 
				payment_mode != 'cod';`

	err := c.db.QueryRow(query, doc_Id).Scan(&count)

	if err != nil {
		return "", err
	}

	query = `SELECT fee
				FROM doctors
				WHERE id = $1;`

	res := strconv.Itoa(count * 150)

	return res, err

}

func (c *adminRepo) ViewSingleUser(user_Id int) (model.UserResponse, error) {

	var user model.UserResponse

	query := `SELECT 
			id,
			first_name,
			last_name,
			email,
			phone,
			last_appointment
			FROM users WHERE id = $1;`

	err := c.db.QueryRow(query,
		user_Id).Scan(
		&user.ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Phone,
		&user.Last_appointment)

	return user, err
}

func (c *adminRepo) ViewSingleDoctor(doc_Id int) (model.DoctorResponse, error) {

	var doctor model.DoctorResponse

	query := `SELECT 
			id,
			first_name,
			last_name,
			email,
			phone,
			approvel
			FROM doctors WHERE id = $1;`

	err := c.db.QueryRow(query,
		doc_Id).Scan(
		&doctor.ID,
		&doctor.First_Name,
		&doctor.Last_Name,
		&doctor.Email,
		&doctor.Phone,
		&doctor.Approvel)

	return doctor, err
}

func (c *adminRepo) ViewAllAppointments(pagenation utils.Filter, filters model.Filter) ([]model.AppointmentByDoctor, utils.Metadata, error) {

	var allappointments []model.AppointmentByDoctor

	query := `SELECT 
				COUNT(*) OVER(),
				time_consult,
				payment_mode,
				payment_status,
				email FROM confirmeds c 
				WHERE id IS NOT NULL `

	var totalRecords int

	i := 1
	var arg []interface{}

	//arg = append(arg,id)

	if len(filters.Day) != 0 {

		query = query + `AND (`

		for j, day := range filters.Day {
			query = query + fmt.Sprintf("c.day_consult ILIKE %d", i)
			if j != len(filters.Day)-1 {
				query = query + " OR "
			}
			day = fmt.Sprint(day, "%")
			arg = append(arg, day)
			i++
		}
		query = query + ")"
	}

	if len(filters.DoctorId) != 0 {
		query = query + `AND (`
		for j, id := range filters.DoctorId {
			query = query + fmt.Sprintf("c.doctor_id ILIKE %d", i)
			if j != len(filters.DoctorId)-1 {
				query = query + " OR "
			}
			// id, _ = strconv.Atoi(fmt.Sprint(id))
			id = fmt.Sprintf(id)
			arg = append(arg, id)
			i++
		}
		query = query + ")"
	}

	if len(filters.Sort) != 0 {
		query = query + fmt.Sprintf(`AND ORDER BY
									c.created_at ASC;`)

	}

	query = query + fmt.Sprintf(`LIMIT $%d OFFSET $%d;`, i, i+1)
	arg = append(arg, pagenation.Limit())
	arg = append(arg, pagenation.Offset())
	log.Println(query)

	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Println("Query preparation failed ", err)
		return nil, utils.Metadata{}, err
	}

	res, err := stmt.Query(arg...)
	if err != nil {
		log.Println("Query execution failed ", err)
		return nil, utils.Metadata{}, err
	}

	defer res.Close()
	// rows, err := c.db.Query(query, doc_id, day, pagenation.Limit(), pagenation.Offset())

	// log.Println(doc_id)
	// log.Println(day)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, utils.Metadata{}, err
	// }

	// defer rows.Close()

	for res.Next() {
		var appointment model.AppointmentByDoctor

		err = res.Scan(
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
	if err := res.Err(); err != nil {
		return allappointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
	}
	log.Println(allappointments)
	log.Println(utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize))
	return allappointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), nil

}

func (c *adminRepo) FindAdmin(username string) (model.AdminResponse, error) {

	log.Println("username of admin:", username)
	var admin model.AdminResponse

	query := `SELECT
			id, 
			username,
			password
			FROM admins WHERE username = $1;`

	err := c.db.QueryRow(query,
		username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.Password)

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

func (c *adminRepo) CreateAdmin(admin model.Admin) error {

	query := `INSERT INTO
				admins (username,password)
				VALUES
				(
					$1, $2
					);`
	err := c.db.QueryRow(
		query, admin.Username,
		admin.Password,
	).Err()
	return err
}

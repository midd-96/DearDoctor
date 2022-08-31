package repo

import (
	"database/sql"
	"dearDoctor/model"
	"dearDoctor/utils"
	"errors"
	"log"
	"time"
)

type DoctorRepository interface {
	AddSlotes(slote model.Slotes) (int, error)
	FindDoctor(email string) (model.DoctorResponse, error)
	InsertDoctor(doctor model.Doctor) (int, error)
	AllDoctors(pagenation utils.Filter) ([]model.DoctorResponse, utils.Metadata, error)
	ListAppointments(pagenation utils.Filter, docId int) ([]model.Appointments, utils.Metadata, error)
	StoreVerificationDetails(email string, code int) error
	VerifyAccount(email string, code int) error
	RequestForPayout(email string, requestAmount float64) (float64, error)
}

type doctorRepo struct {
	db *sql.DB
}

func NewDoctorRepo(db *sql.DB) DoctorRepository {
	return &doctorRepo{
		db: db,
	}
}

func (c *doctorRepo) RequestForPayout(email string, requestAmount float64) (float64, error) {

	var walletBalance float64
	log.Println("mail id recieved at repo :", email)
	query := `SELECT wallet FROM 
				payouts WHERE 
				username = $1;`
	err := c.db.QueryRow(query, email).Scan(&walletBalance)
	log.Println("Error from repo of doctor while scanning data from payouts", err, walletBalance)

	if requestAmount <= walletBalance {
		query = `UPDATE payouts
					SET last_requested_amount = $1,
						requested_time = $2,
						approvel = $3
					WHERE username= $4`

		Requested_At, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		err = c.db.QueryRow(query, requestAmount, Requested_At, false, email).Err()

		return requestAmount, nil

	}
	err = errors.New("Requested amount higher than wallet balance")
	return requestAmount, err

}

func (c *doctorRepo) VerifyAccount(email string, code int) error {

	var id int

	query := `SELECT id FROM 
				verifications WHERE 
				email = $1 AND code = $2;`
	err := c.db.QueryRow(query, email, code).Scan(&id)

	if err == sql.ErrNoRows {
		return errors.New("Invalid verification code.")
	}

	if err != nil {
		return err
	}

	query = `UPDATE doctors 
				SET
				 verification = $1
				WHERE
				 email = $2 ;`
	err = c.db.QueryRow(query, true, email).Err()
	log.Println("Updating verification: ", err)
	if err != nil {
		return err
	}

	return nil
}

func (c *doctorRepo) StoreVerificationDetails(email string, code int) error {

	query := `INSERT INTO 
				verifications(email, code)
				VALUES( $1, $2);`

	err := c.db.QueryRow(query, email, code).Err()

	return err

}

func (c *doctorRepo) ListAppointments(pagenation utils.Filter, docId int) ([]model.Appointments, utils.Metadata, error) {

	var appointments []model.Appointments
	doctorId := docId

	query := `SELECT 
				COUNT(*) OVER(),
				day_consult,
				time_consult,
				payment_mode,
				payment_status,
				email
				FROM confirmeds WHERE doctor_id = $1 
				LIMIT $2 OFFSET $3`

	rows, err := c.db.Query(query, doctorId, pagenation.Limit(), pagenation.Offset())
	if err != nil {
		return nil, utils.Metadata{}, err
	}

	var totalRecords int

	defer rows.Close()

	for rows.Next() {
		var appointment model.Appointments

		err = rows.Scan(
			&totalRecords,
			&appointment.Day_consult,
			&appointment.Time_consult,
			&appointment.Payment_mode,
			&appointment.Payment_status,
			&appointment.Email,
		)

		if err != nil {
			return appointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
		}
		appointments = append(appointments, appointment)
	}

	if err := rows.Err(); err != nil {
		return appointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
	}
	log.Println(appointments)
	log.Println(utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize))
	return appointments, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), nil
}

func (c *doctorRepo) AddSlotes(slote model.Slotes) (int, error) {

	var id int

	query := `INSERT INTO slotes(
		doctor_id,
		available_day,
		time_from,
		time_to) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;`

	err := c.db.QueryRow(query,
		slote.Docter_id,
		slote.Available_day,
		slote.Time_from,
		slote.Time_to).Scan(&id)

	return id, err
}

func (c *doctorRepo) FindDoctor(email string) (model.DoctorResponse, error) {

	var doctor model.DoctorResponse

	query := `SELECT 
				id,
				first_name,
				last_name,
				email,
				password,
				phone,
				approvel
				FROM doctors 
				WHERE email = $1;`

	err := c.db.QueryRow(query,
		email).Scan(
		&doctor.ID,
		&doctor.First_Name,
		&doctor.Last_Name,
		&doctor.Email,
		&doctor.Password,
		&doctor.Phone,
		&doctor.Approvel,
	)

	return doctor, err
}

func (c *doctorRepo) InsertDoctor(doctor model.Doctor) (int, error) {

	var id int

	query := `INSERT INTO doctors(
			first_name,
			last_name,
			email,
			phone,
			password,
			department,
			specialization
			)
			VALUES
			($1, $2, $3, $4, $5, $6, $7)
			RETURNING id;`

	err := c.db.QueryRow(query,
		doctor.First_name,
		doctor.Last_name,
		doctor.Email,
		doctor.Phone,
		doctor.Password,
		doctor.Department,
		doctor.Specialization).Scan(
		&id,
	)
	return id, err
}

func (c *doctorRepo) AllDoctors(pagenation utils.Filter) ([]model.DoctorResponse, utils.Metadata, error) {

	var doctors []model.DoctorResponse

	query := `SELECT 
				COUNT(*) OVER(),
				id,
				first_name,
				last_name,
				email,
				password,
				phone,
				approvel
				FROM doctors 
				LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, pagenation.Limit(), pagenation.Offset())
	if err != nil {
		log.Println("Error", "Query prepare failed: ", err)
		return nil, utils.Metadata{}, err
	}

	var totalRecords int

	defer rows.Close()

	for rows.Next() {
		var doctor model.DoctorResponse

		err = rows.Scan(
			&totalRecords,
			&doctor.ID,
			&doctor.First_Name,
			&doctor.Last_Name,
			&doctor.Email,
			&doctor.Password,
			&doctor.Phone,
			&doctor.Approvel,
		)

		if err != nil {
			return doctors, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
		}
		doctors = append(doctors, doctor)
	}

	if err := rows.Err(); err != nil {
		return doctors, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
	}
	log.Println(doctors)
	log.Println(utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize))
	return doctors, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), nil

}

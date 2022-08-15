package repo

import (
	"database/sql"
	"dearDoctor/model"
)

type DoctorRepository interface {
	AddSlotes(slote model.Slotes) (int, error)
	FindDoctor(email string) (model.DoctorResponse, error)
	InsertDoctor(doctor model.Doctor) (int, error)
}

type doctorRepo struct {
	db *sql.DB
}

func NewDoctorRepo(db *sql.DB) DoctorRepository {
	return &doctorRepo{
		db: db,
	}
}

func (c *doctorRepo) AddSlotes(slote model.Slotes) (int, error) {

	var id int

	query := `INSERT INTO slotes(
		id,
		doctor_id,
		available_day,
		time_from,
		time_to,
		status
		) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNIG id;`

	err := c.db.QueryRow(query,
		slote.Id,
		slote.Docter_id,
		slote.Available_day,
		slote.Time_from,
		slote.Time_to,
		slote.Status).Scan(&id)

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

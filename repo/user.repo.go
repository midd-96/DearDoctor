package repo

import (
	"database/sql"
	"dearDoctor/model"
	"dearDoctor/utils"
	"errors"
	"fmt"
	"log"
)

type UserRepository interface {
	AllUsers(pagenation utils.Filter) ([]model.UserResponse, utils.Metadata, error)
	FindUser(email string) (model.UserResponse, error)
	InsertUser(user model.User) (int, error)
	AddAppointment(confirm model.Confirmed) (int, error)
	ManageUsers(email string) error
	UpdateUser(data model.User) error
	FindAppointmentById(id int) (float64, error)
	StoreVerificationDetails(email string, code int) error
	VerifyAccount(email string, code int) error
	Payment(payment model.PaymentDetails) error
	FindUserByAppointmentId(id int) (model.UserResponse, error)
	FindUserById(id int) (model.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (c *userRepo) Payment(data model.PaymentDetails) error {

	var arg []interface{}

	insertQuery := `
				INSERT INTO
				   payments ( appointment_id, user_id, payment_type, created_at, payment_id, razor_order_id, payment_signature )
				VALUES
				   (
				      $1, $2, $3, $4, $5, $6, $7
				   )
				RETURNING id;`

	updateQuery := `
					UPDATE
					   confirmeds
					SET
					   payment_status = true, payment_mode = $1 ,updated_at = $2
					
	`

	err := c.db.QueryRow(
		insertQuery,
		data.Appointment_ID,
		data.User_ID,
		data.PaymentType,
		data.Created_At,
		data.Razorpay_payment_id,
		data.Razorpay_order_id,
		data.Razorpay_signature).Scan(&data.ID)

	if err != nil {
		return err
	}

	arg = append(arg, data.PaymentType)
	arg = append(arg, data.Updated_At)

	i := 3

	updateQuery = updateQuery + `
								WHERE
								id = $` + fmt.Sprintf(`%d;`, i)

	arg = append(arg, data.Appointment_ID)

	stmt, err := c.db.Prepare(updateQuery)

	if err != nil {
		log.Println("Error", err)
		log.Println("Error", "Query prepare failed")
		return err
	}

	_, err = stmt.Query(arg...)

	if err != nil {
		log.Println("Error", err)
		log.Println("Error", "Query Exec failed")
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) VerifyAccount(email string, code int) error {

	var id int

	query := `SELECT id FROM 
				verifications WHERE 
				email = $1 AND code = $2;`
	err := c.db.QueryRow(query, email, code).Scan(&id)

	if err == sql.ErrNoRows {
		return errors.New("Invalid verification code/Email")
	}

	if err != nil {
		return err
	}

	query = `UPDATE users 
				SET
				 verification = $1
				WHERE
				 email = $2 ;`
	err = c.db.QueryRow(query, true, email).Err()
	log.Println("Updating User verification: ", err)
	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) StoreVerificationDetails(email string, code int) error {

	query := `INSERT INTO 
				verifications(email, code)
				VALUES( $1, $2);`

	err := c.db.QueryRow(query, email, code).Err()

	return err

}

func (c *userRepo) FindAppointmentById(id int) (float64, error) {

	var fee float64
	query := `SELECT 
				fee
				FROM confirmeds 
				WHERE id = $1 AND payment_status = false;`

	err := c.db.QueryRow(query,
		id).Scan(&fee)

	return fee, err
}

func (c *userRepo) AllUsers(pagenation utils.Filter) ([]model.UserResponse, utils.Metadata, error) {

	var users []model.UserResponse

	query := `SELECT 
				COUNT(*) OVER(),
				id,
				first_name,
				last_name,
				email,
				password,
				phone,
				last_appointment
				FROM users 
				LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, pagenation.Limit(), pagenation.Offset())
	if err != nil {
		return nil, utils.Metadata{}, err
	}

	var totalRecords int

	defer rows.Close()

	for rows.Next() {
		var user model.UserResponse

		err = rows.Scan(
			&totalRecords,
			&user.ID,
			&user.First_Name,
			&user.Last_Name,
			&user.Email,
			&user.Password,
			&user.Phone,
			&user.Last_appointment,
		)

		if err != nil {
			return users, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return users, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), err
	}
	log.Println(users)
	log.Println(utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize))
	return users, utils.ComputeMetaData(totalRecords, pagenation.Page, pagenation.PageSize), nil

}

func (c *userRepo) FindUserByAppointmentId(id int) (model.UserResponse, error) {

	var email string
	query := `SELECT email 
				FROM confirmeds
				WHERE id = $1;`

	err := c.db.QueryRow(query, id).Scan(&email)
	var user model.UserResponse

	query = `SELECT 
				id,
				first_name,
				last_name,
				email,
				password,
				phone,
				last_appointment
				FROM users 
				WHERE email = $1;`

	err = c.db.QueryRow(query,
		email).Scan(
		&user.ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Last_appointment,
	)

	return user, err
}
func (c *userRepo) FindUserById(id int) (model.User, error) {

	var user model.User

	query := `SELECT 
				id,
				first_name,
				last_name,
				email,
				password,
				phone,
				last_appointment
				FROM users 
				WHERE id = $1;`

	err := c.db.QueryRow(query,
		id).Scan(
		&user.ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Last_appointment,
	)

	return user, err
}

func (c *userRepo) FindUser(email string) (model.UserResponse, error) {

	var user model.UserResponse

	query := `SELECT 
				id,
				first_name,
				last_name,
				email,
				password,
				phone,
				last_appointment
				FROM users 
				WHERE email = $1;`

	err := c.db.QueryRow(query,
		email).Scan(
		&user.ID,
		&user.First_Name,
		&user.Last_Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Last_appointment,
	)

	return user, err
}

func (c *userRepo) InsertUser(user model.User) (int, error) {

	var id int

	query := `INSERT INTO users(
			first_name,
			last_name,
			email,
			phone,
			password,
			last_appointment
			)
			VALUES
			($1, $2, $3, $4, $5, $6)
			RETURNING id;`

	err := c.db.QueryRow(query,
		user.First_Name,
		user.Last_Name,
		user.Email,
		user.Phone,
		user.Password,
		user.Last_appointment).Scan(
		&id,
	)
	return id, err
}

func (c *userRepo) AddAppointment(confirm model.Confirmed) (int, error) {
	var ID int
	var doc_email string
	query := `INSERT INTO confirmeds(
		day_consult,
		time_consult,
		payment_mode,
		payment_status,
		fee,
		email,
		doctor_id
		)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING ID;`

	err1 := c.db.QueryRow(query,
		confirm.Day_consult,
		confirm.Time_consult,
		confirm.Payment_mode,
		confirm.Payment_status,
		confirm.Fee,
		confirm.Email,
		confirm.Doctor_id).Scan(
		&ID,
	)

	//to store email id of the doctor who got an appointment
	query = `SELECT email FROM 
				doctors WHERE 
				id = $1;`
	err := c.db.QueryRow(query, confirm.Doctor_id).Scan(&doc_email)

	//checks weather it is his first appointment or not
	var status bool
	query = `SELECT * FROM 
				payouts WHERE 
				username = $1;`
	err = c.db.QueryRow(query, doc_email).Scan(&status)
	if err == sql.ErrNoRows && confirm.Payment_mode == "cod" {

		//if it is his first appointments create new row.
		query = `INSERT INTO payouts(
					username,
					wallet)
					VALUES (
						$1, $2
					) RETURNING username;`

		err = c.db.QueryRow(query, doc_email, (float64(confirm.Fee)*0.75)-float64(confirm.Fee)).Scan(&doc_email)

	} else if err != sql.ErrNoRows && confirm.Payment_mode == "cod" {

		query = `UPDATE payouts SET
			wallet = wallet - $1;`

		err = c.db.QueryRow(query, (float64(confirm.Fee) - (float64(confirm.Fee) * 0.75))).Err()

	}

	if err == sql.ErrNoRows && confirm.Payment_mode != "cod" {

		//if it is his first appointments create new row.
		query = `INSERT INTO payouts(
					username,
					wallet)
					VALUES (
						$1, $2
					) RETURNING username;`

		err = c.db.QueryRow(query, doc_email, float64(confirm.Fee)*0.75).Scan(&doc_email)

	} else if err != sql.ErrNoRows && confirm.Payment_mode != "cod" {

		query = `UPDATE payouts SET
			wallet = wallet + $1;`

		err = c.db.QueryRow(query, float64(confirm.Fee)*0.75).Err()

	}

	return ID, err1
}

func (c *userRepo) ManageUsers(email string) error {
	//Query
	query := `UPDATE users 
			SET is_active = $1 
			WHERE email = $2 ;`

	err := c.db.QueryRow(query,
		email).Err()

	return err
}

func (c *userRepo) UpdateUser(data model.User) error {

	query := `
				UPDATE
					users 
				 SET`
	i := 1
	var arg []interface{}

	if data.First_Name != "" {
		query = query + `first_name = $` + fmt.Sprintf(`%d`, i)
		i++
		arg = append(arg, data.First_Name)
	}

	if data.Last_Name != "" {
		if i > 1 {
			query = query + `, `
		}
		query = query + `last_name = $` + fmt.Sprintf(`%d`, i)
		i++
		arg = append(arg, data.Last_Name)
	}

	if data.Email != "" {
		if i > 1 {
			query = query + `, `
		}
		query = query + `email = $` + fmt.Sprintf(`%d`, i)
		arg = append(arg, data.Email)
		i++
	}

	if data.Password != "" {
		if i > 1 {
			query = query + `, `
		}
		query = query + `password = $` + fmt.Sprintf(`%d`, i)
		arg = append(arg, data.Password)
		i++
	}

	statement, err := c.db.Prepare(query)

	if err != nil {
		log.Println("Error ", "error in preparing query: ", err)
		return err
	}

	_, err = statement.Query(arg...)

	if err != nil {
		log.Println("Error ", "error in query execution: ", err)
		return err
	}

	return nil
}

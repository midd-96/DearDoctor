package repo

import (
	"database/sql"
	"dearDoctor/model"
	"dearDoctor/utils"
	"fmt"
	"log"
	//"github.com/lib/pq"
)

type UserRepository interface {
	AllUsers(pagenation utils.Filter) ([]model.UserResponse, utils.Metadata, error)
	FindUser(email string) (model.UserResponse, error)
	InsertUser(user model.User) (int, error)
	AddAppointment(confirm model.Confirmed) (int, error)
	ManageUsers(email string) error
	UpdateUser(data model.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
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
	log.Println(confirm)
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

	err := c.db.QueryRow(query,
		confirm.Day_consult,
		confirm.Time_consult,
		confirm.Payment_mode,
		confirm.Payment_status,
		confirm.Fee,
		confirm.Email,
		confirm.Doctor_id).Scan(
		&ID,
	)
	return ID, err
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

	// if data.Phone_Number != 0 {
	// 	if i > 1 {
	// 		query = query + `, `
	// 	}
	// 	query = query + `phone_number = $` + fmt.Sprintf(`%d`, i)
	// 	arg = append(arg, data.Phone_Number)
	// 	i++
	// }

	// if data.IsVerified {
	// 	if i > 1 {
	// 		query = query + `, `
	// 	}
	// 	query = query + `is_verified = $` + fmt.Sprintf(`%d`, i)
	// 	arg = append(arg, data.IsVerified)
	// 	i++
	// }

	// if i > 1 {
	// 	query = query + `, `
	// }
	// query = query + `updated_at = $` + fmt.Sprintf(`%d
	// 										WHERE id = $%d;`, i, i+1)
	// arg = append(arg, data.Updated_At)
	// arg = append(arg, data.ID)

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

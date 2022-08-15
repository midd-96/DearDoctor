package repo

import (
	"database/sql"
	"dearDoctor/model"
	"fmt"
	"log"
)

type UserRepository interface {
	AllUsers() ([]model.UserResponse, error)
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

func (c *userRepo) AllUsers() ([]model.UserResponse, error) {

	var users []model.UserResponse

	//stores related query to a variable
	query := `SELECT 
				id,
				first_name,
				last_name,
				password,
				email,
				phone,
				last_appointment
				FROM users 
				WHERE email = $1;`

	rows, err := c.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	//Loop through each users (raw wise)
	for rows.Next() {
		var user model.UserResponse
		err := rows.Scan(
			&user.ID,
			&user.First_Name,
			&user.Last_Name,
			&user.Email,
			&user.Phone,
			&user.Last_appointment)

		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (c *userRepo) FindUser(email string) (model.UserResponse, error) {

	var user model.UserResponse

	query := `SELECT 
				id,
				first_name,
				last_name,
				password,
				email,
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
	var id int
	query := `INSERT INTO confirmeds(
		id,
		day_consult,
		time_consult,
		payment_mode,
		payment_status,
		fee,
		email,
		doctor_id
		)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;`

	err := c.db.QueryRow(query,
		confirm.Id,
		confirm.Day_consult,
		confirm.Time_consult,
		confirm.Payment_mode,
		confirm.Payment_status,
		confirm.Fee,
		confirm.Email,
		confirm.Doctor_id).Scan(
		&id,
	)
	return id, err
}

func (c *userRepo) ManageUsers(email string) error {
	//Query
	query := `UPDATE users 
			SET is_active = $1 
			WHERE email = $2 ;`

	err := c.db.QueryRow(query,
		email).Err()
	// err := c.db.QueryRow(query,
	// 	Role,
	// 	email).Err()

	return err
}

func (c *userRepo) AddAddress(address model.Address) (int, error) {

	var id int

	query := `INSERT INTO address(
				type,
				user_id,
				house_name,
				street_name,
				landmark,
				district,
				state,
				country,
				pincode,
				created_at)
				VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)	
				RETURNING id;`

	err := c.db.QueryRow(query,
		address.AddressType,
		address.User_id,
		address.HouseName,
		address.StreetName,
		address.Landmark,
		address.District,
		address.State,
		address.Country,
		address.PinCode,
		address.Created_At).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, err
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

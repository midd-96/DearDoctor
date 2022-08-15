package repo

import (
	"database/sql"
	"dearDoctor/model"
)

type AdminRepository interface {
	FindAdmin(username string) (model.AdminResponse, error)
	AddDept(department model.Departments) error
	UpdateApproveFee(approvel model.ApproveAndFee) error
}

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepo(db *sql.DB) AdminRepository {
	return &adminRepo{
		db: db,
	}
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

func (c *adminRepo) UpdateApproveFee(approvel model.ApproveAndFee) error {

	var query string

	var err error

	query = `
				UPDATE
				   doctors 
				SET
				   approvel = $1, fee=$2
				WHERE
				   category_id = $3 ;`

	err = c.db.QueryRow(query, approvel.Approve, approvel.Fee, approvel.Doctor_id).Err()

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

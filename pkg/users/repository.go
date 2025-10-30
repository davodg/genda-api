package users

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type UserRepo struct {
	postgresDB *sql.DB
}

func NewUserRepository(postgresDB *sql.DB) *UserRepo {
	return &UserRepo{postgresDB: postgresDB}
}

func (i *UserRepo) CreateUser(user User) (*User, error) {
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	if user.Uuid == "" {
		user.Uuid = uuid.New().String()
	}

	const sqlStmt = `
		INSERT INTO users
			(id, name, email, is_admin, type, birth_date, social_number, phone_number, gender, photo_url)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`
	_, err := i.postgresDB.Exec(sqlStmt,
		user.Id,
		user.Name,
		user.Email,
		user.IsAdmin,
		user.Type,
		user.BirthDate,
		user.SocialNumber,
		user.PhoneNumber,
		user.Gender,
		user.PhotoUrl,
	)
	if err != nil {
		log.Println("An error occurred while creating user", err)
		return nil, err
	}
	return &user, nil
}

func (i *UserRepo) GetUsers(page int, limit int, name string, email string, userId string) (*GetUserResponse, error) {
	res := GetUserResponse{
		Page:  page,
		Limit: limit,
		Users: []User{},
	}

	// ----- filtros dinâmicos (parametrizados) -----
	var where []string
	var args []any
	arg := 1

	if userId != "" {
		where = append(where, "id = $"+strconv.Itoa(arg))
		args = append(args, userId)
		arg++
	}
	if name != "" {
		where = append(where, "name ILIKE $"+strconv.Itoa(arg))
		args = append(args, "%"+name+"%")
		arg++
	}
	if email != "" {
		where = append(where, "email = $"+strconv.Itoa(arg))
		args = append(args, email)
		arg++
	}

	filter := ""
	if len(where) > 0 {
		filter = " WHERE " + strings.Join(where, " AND ")
	}

	// total
	countSQL := "SELECT count(*) FROM users" + filter
	if err := i.postgresDB.QueryRow(countSQL, args...).Scan(&res.Total); err != nil {
		log.Println("An error occurred while counting users", err)
		return nil, err
	}

	// página
	listSQL := `
		SELECT
			id,
			name,
			email,
			is_admin,
			type,
			birth_date,
			social_number,
			phone_number,
			gender,
			photo_url
		FROM users
	` + filter + `
		ORDER BY name ASC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1) + `;
	`
	offset := (page * limit) - limit
	args = append(args, limit, offset)

	rows, err := i.postgresDB.Query(listSQL, args...)
	if err != nil {
		log.Println("An error occurred while getting users", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user, err := i.formatUser(rows)
		if err != nil {
			log.Println("An error occurred while formatting user", err)
			return nil, err
		}
		res.Users = append(res.Users, *user)
	}
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error while getting users", err)
		return nil, err
	}

	return &res, nil
}

func (i *UserRepo) UpdateUser(id string, user User) (*User, error) {
	const sqlStmt = `
		UPDATE users
		SET name = $1, email = $2, is_admin = $3, photo_url = $4, type = $5,
		    birth_date = $6, social_number = $7, phone_number = $8, gender = $9
		WHERE id = $10
	`
	_, err := i.postgresDB.Exec(sqlStmt,
		user.Name, user.Email, user.IsAdmin, user.PhotoUrl, user.Type,
		user.BirthDate, user.SocialNumber, user.PhoneNumber, user.Gender,
		id,
	)
	if err != nil {
		log.Println("An error occurred while updating user", err)
		return nil, err
	}
	return &user, nil
}

func (i *UserRepo) DeleteUser(id string) error {
	const sqlStmt = `DELETE FROM users WHERE id = $1`
	if _, err := i.postgresDB.Exec(sqlStmt, id); err != nil {
		log.Println("An error occurred while deleting user", err)
		return err
	}
	return nil
}

func (i *UserRepo) formatUser(row *sql.Rows) (*User, error) {
	u := User{}
	if err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.IsAdmin,
		&u.Type,
		&u.BirthDate,
		&u.SocialNumber,
		&u.PhoneNumber,
		&u.Gender,
		&u.PhotoUrl,
	); err != nil {
		log.Println("An error occurred while scanning user", err)
		return nil, err
	}
	return &u, nil
}

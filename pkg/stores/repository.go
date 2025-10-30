package stores

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type StoreRepo struct {
	postgresDB *sql.DB
}

func NewStoreRepository(postgresDB *sql.DB) *StoreRepo {
	return &StoreRepo{postgresDB: postgresDB}
}

func (i *StoreRepo) CreateStore(store Store) (*Store, error) {
	if store.Id == "" {
		store.Id = uuid.New().String()
	}

	const insertSQL = `
		INSERT INTO stores
			(id, name, owner_id, "type", location)
		VALUES
			($1,$2,$3,$4,$5::json)
	`
	_, err := i.postgresDB.Exec(insertSQL,
		store.Id,
		store.Name,
		store.OwnerId,
		store.Type,
		store.Location, // string JSON; ::json valida o formato
	)
	if err != nil {
		log.Println("An error occurred while creating store", err)
		return nil, err
	}
	return &store, nil
}

func (i *StoreRepo) GetStores(page int, limit int, name string, storeType string, storeId string) (*GetStoreResponse, error) {
	res := GetStoreResponse{
		Page:   page,
		Limit:  limit,
		Stores: []Store{},
	}

	var where []string
	var args []any
	arg := 1

	if storeId != "" {
		where = append(where, "id = $"+strconv.Itoa(arg))
		args = append(args, storeId)
		arg++
	}
	if name != "" {
		where = append(where, "name ILIKE $"+strconv.Itoa(arg))
		args = append(args, "%"+name+"%")
		arg++
	}
	if storeType != "" {
		where = append(where, `"type" = $`+strconv.Itoa(arg))
		args = append(args, storeType)
		arg++
	}

	filter := ""
	if len(where) > 0 {
		filter = " WHERE " + strings.Join(where, " AND ")
	}

	// total
	countSQL := "SELECT count(*) FROM stores" + filter
	if err := i.postgresDB.QueryRow(countSQL, args...).Scan(&res.Total); err != nil {
		log.Println("An error occurred while counting stores", err)
		return nil, err
	}

	// página
	listSQL := `
		SELECT
			id,
			name,
			owner_id,
			COALESCE("type",'') AS type,
			COALESCE(location::text,'') AS location
		FROM stores
	` + filter + `
		ORDER BY name ASC
		LIMIT $` + strconv.Itoa(arg) + ` OFFSET $` + strconv.Itoa(arg+1) + `;
	`
	offset := (page * limit) - limit
	args = append(args, limit, offset)

	rows, err := i.postgresDB.Query(listSQL, args...)
	if err != nil {
		log.Println("An error occurred while getting stores", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		store, err := i.formatStore(rows)
		if err != nil {
			log.Println("An error occurred while formatting store", err)
			return nil, err
		}
		res.Stores = append(res.Stores, *store)
	}
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error while getting stores", err)
		return nil, err
	}

	return &res, nil
}

func (i *StoreRepo) GetStore(id string) (*Store, error) {
	res, err := i.GetStores(1, 1, "", "", id)
	if err != nil {
		return nil, err
	}
	return &res.Stores[0], nil
}

func (i *StoreRepo) UpdateStore(id string, store Store) (*Store, error) {
	const sqlStmt = `
		UPDATE stores
		SET name = $1, "type" = $2, owner_id = $3, location = $4::json
		WHERE id = $5
	`
	_, err := i.postgresDB.Exec(sqlStmt,
		store.Name,
		store.Type,
		store.OwnerId,
		store.Location, // string JSON
		id,
	)
	if err != nil {
		log.Println("An error occurred while updating store", err)
		return nil, err
	}
	return &store, nil
}

func (i *StoreRepo) DeleteStore(id string) error {
	const sqlStmt = `DELETE FROM stores WHERE id = $1`
	if _, err := i.postgresDB.Exec(sqlStmt, id); err != nil {
		log.Println("An error occurred while deleting store", err)
		return err
	}
	return nil
}

func (i *StoreRepo) formatStore(row *sql.Rows) (*Store, error) {
	s := Store{}
	if err := row.Scan(
		&s.Id,
		&s.Name,
		&s.OwnerId,
		&s.Type,
		&s.Location, // string (JSON)
	); err != nil {
		log.Println("An error occurred while scanning store", err)
		return nil, err
	}
	return &s, nil
}

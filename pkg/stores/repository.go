package stores

import (
	"database/sql"
	"encoding/json"
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
		store.Location,
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
			type,
			location
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

func (i *StoreRepo) GetStore(id string) (*GetStoreByIdResponse, error) {
	queryStmt, err := i.getStoreStatements(id)
	if err != nil {
		log.Println("An error occurred while building get store query", err)
		return nil, err
	}

	row, err := i.postgresDB.Query(queryStmt)
	if err != nil {
		log.Println("An error occurred while getting store", err)
		return nil, err
	}
	defer row.Close()

	var fullStore FullStore
	fullStore.Ratings = []StoreRating{}
	fullStore.Plans = []StorePlan{}

	for row.Next() {
		s, err := i.formatFullStore(row)
		if err != nil {
			log.Println("An error occurred while formatting store", err)
			return nil, err
		}
		fullStore.Store = *s

		// Additional formatting for Availability, Ratings, Plans would go here
		// depending on the selected fields in the query.

	}
	if err := row.Err(); err != nil {
		log.Println("Row iteration error while getting store", err)
		return nil, err
	}

	res := GetStoreByIdResponse{
		Total: 1,
		Store: fullStore,
	}

	return &res, nil
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

func (i *StoreRepo) CreateStoreRating(rating StoreRating) (*StoreRating, error) {
	if rating.Id == "" {
		rating.Id = uuid.New().String()
	}

	const insertSQL = `
		INSERT INTO store_ratings
			(id, store_id, user_id, rating, message)
		VALUES
			($1,$2,$3,$4,$5)
	`
	_, err := i.postgresDB.Exec(insertSQL,
		rating.Id,
		rating.StoreId,
		rating.UserId,
		rating.Rating,
		rating.Message,
	)
	if err != nil {
		log.Println("An error occurred while creating store rating", err)
		return nil, err
	}
	return &rating, nil
}

func (i *StoreRepo) GetStoreRatings(storeId string, page int, limit int) (*GetStoreRatingsResponse, error) {
	res := GetStoreRatingsResponse{
		Page:    page,
		Limit:   limit,
		Ratings: []StoreRating{},
	}

	// total
	countSQL := "SELECT count(*) FROM store_ratings WHERE store_id = $1"
	if err := i.postgresDB.QueryRow(countSQL, storeId).Scan(&res.Total); err != nil {
		log.Println("An error occurred while counting store ratings", err)
		return nil, err
	}

	// página
	listSQL := `
		SELECT
			id,
			store_id,
			user_id,
			rating,
			message,
			created_at
		FROM store_ratings
		WHERE store_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	offset := (page * limit) - limit

	rows, err := i.postgresDB.Query(listSQL, storeId, limit, offset)
	if err != nil {
		log.Println("An error occurred while getting store ratings", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rating StoreRating
		err := rows.Scan(
			&rating.Id,
			&rating.StoreId,
			&rating.UserId,
			&rating.Rating,
			&rating.Message,
			&rating.CreatedAt,
		)
		if err != nil {
			log.Println("An error occurred while scanning store rating", err)
			return nil, err
		}
		res.Ratings = append(res.Ratings, rating)
	}
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error while getting store ratings", err)
		return nil, err
	}

	return &res, nil
}

func (i *StoreRepo) UpdateStoreRating(id string, rating StoreRating) (*StoreRating, error) {
	const sqlStmt = `
		UPDATE store_ratings
		SET rating = $1, message = $2
		WHERE id = $3
	`
	_, err := i.postgresDB.Exec(sqlStmt,
		rating.Rating,
		rating.Message,
		id,
	)
	if err != nil {
		log.Println("An error occurred while updating store rating", err)
		return nil, err
	}
	return &rating, nil
}

func (i *StoreRepo) DeleteStoreRating(id string) error {
	const sqlStmt = `DELETE FROM store_ratings WHERE id = $1`
	if _, err := i.postgresDB.Exec(sqlStmt, id); err != nil {
		log.Println("An error occurred while deleting store rating", err)
		return err
	}
	return nil
}

func (i *StoreRepo) CreateStorePlan(plan StorePlan) (*StorePlan, error) {
	if plan.Id == "" {
		plan.Id = uuid.New().String()
	}

	const insertSQL = `
		INSERT INTO store_plans
			(id, store_id, name, price, currency, plan_type, frequency)
		VALUES
			($1,$2,$3,$4,$5,$6,$7)
	`
	_, err := i.postgresDB.Exec(insertSQL,
		plan.Id,
		plan.StoreId,
		plan.Name,
		plan.Price,
		plan.Currency,
		plan.PlanType,
		plan.Frequency,
	)
	if err != nil {
		log.Println("An error occurred while creating store plan", err)
		return nil, err
	}
	return &plan, nil
}

func (i *StoreRepo) UpdateStorePlan(id string, plan StorePlan) (*StorePlan, error) {
	const sqlStmt = `
		UPDATE store_plans
		SET name = $1, price = $2, currency = $3, plan_type = $4, frequency = $5
		WHERE id = $6
	`
	_, err := i.postgresDB.Exec(sqlStmt,
		plan.Name,
		plan.Price,
		plan.Currency,
		plan.PlanType,
		plan.Frequency,
		id,
	)
	if err != nil {
		log.Println("An error occurred while updating store plan", err)
		return nil, err
	}
	return &plan, nil
}

func (i *StoreRepo) DeleteStorePlan(id string) error {
	const sqlStmt = `DELETE FROM store_plans WHERE id = $1`
	if _, err := i.postgresDB.Exec(sqlStmt, id); err != nil {
		log.Println("An error occurred while deleting store plan", err)
		return err
	}
	return nil
}

func (i *StoreRepo) GetStorePlans(storeId string) (*[]StorePlan, error) {
	const sqlStmt = `
		SELECT
			id,
			store_id,
			name,
			price,
			currency,
			plan_type,
			frequency,
			created_at,
			updated_at
		FROM store_plans
		WHERE store_id = $1;
	`
	rows, err := i.postgresDB.Query(sqlStmt, storeId)
	if err != nil {
		log.Println("An error occurred while getting store plans", err)
		return nil, err
	}
	defer rows.Close()

	var plans []StorePlan
	for rows.Next() {
		var plan StorePlan
		err := rows.Scan(
			&plan.Id,
			&plan.StoreId,
			&plan.Name,
			&plan.Price,
			&plan.Currency,
			&plan.PlanType,
			&plan.Frequency,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			log.Println("An error occurred while scanning store plan", err)
			return nil, err
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error while getting store plans", err)
		return nil, err
	}

	return plans, nil
}

func (i *StoreRepo) CreateStoreAvailability(availability StoreAvailability) (*StoreAvailability, error) {
	if availability.Id == "" {
		availability.Id = uuid.New().String()
	}

	const insertSQL = `
		INSERT INTO store_availability
			(id, store_id, availability)
		VALUES
			($1,$2,$3)
	`
	_, err := i.postgresDB.Exec(insertSQL,
		availability.Id,
		availability.StoreId,
		availability.Availability,
	)
	if err != nil {
		log.Println("An error occurred while creating store availability", err)
		return nil, err
	}
	return &availability, nil
}

func (i *StoreRepo) UpdateStoreAvailability(id string, availability StoreAvailability) (*StoreAvailability, error) {
	const sqlStmt = `
		UPDATE store_availability
		SET availability = $1
		WHERE id = $2
	`
	_, err := i.postgresDB.Exec(sqlStmt,
		availability.Availability,
		id,
	)
	if err != nil {
		log.Println("An error occurred while updating store availability", err)
		return nil, err
	}
	return &availability, nil
}

func (i *StoreRepo) DeleteStoreAvailability(id string) error {
	const sqlStmt = `DELETE FROM store_availability WHERE id = $1`
	if _, err := i.postgresDB.Exec(sqlStmt, id); err != nil {
		log.Println("An error occurred while deleting store availability", err)
		return err
	}
	return nil
}

func (i *StoreRepo) GetStoreAvailability(storeId string) (*StoreAvailability, error) {
	const sqlStmt = `
		SELECT
			id,
			store_id,
			availability,
			created_at,
			updated_at
		FROM store_availability
		WHERE store_id = $1;
	`
	var availability StoreAvailability
	err := i.postgresDB.QueryRow(sqlStmt, storeId).Scan(
		&availability.Id,
		&availability.StoreId,
		&availability.Availability,
		&availability.CreatedAt,
		&availability.UpdatedAt,
	)
	if err != nil {
		log.Println("An error occurred while getting store availability", err)
		return nil, err
	}
	return &availability, nil
}

func (i *StoreRepo) CreateStoreAppointment(appointment StoreAppointment) (*StoreAppointment, error) {
	if appointment.Id == "" {
		appointment.Id = uuid.New().String()
	}

	const insertSQL = `
		INSERT INTO store_appointments
			(id, store_id, user_id, start_at, end_at, status, hold_expires_at, price, currency, fee_platform, payment_id, notes)
		VALUES
			($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`
	_, err := i.postgresDB.Exec(insertSQL,
		appointment.Id,
		appointment.StoreId,
		appointment.UserId,
		appointment.StartAt,
		appointment.EndAt,
		appointment.Status,
		appointment.HoldExpiresAt,
		appointment.Price,
		appointment.Currency,
		appointment.FeePlatform,
		appointment.PaymentId,
		appointment.Notes,
	)
	if err != nil {
		log.Println("An error occurred while creating store appointment", err)
		return nil, err
	}
	return &appointment, nil
}

func (i *StoreRepo) GetStoreAppointments(storeId string, page int, limit int) (*GetStoreAppointmentsResponse, error) {
	res := GetStoreAppointmentsResponse{
		Page:         page,
		Limit:        limit,
		Appointments: []StoreAppointment{},
	}

	// total
	countSQL := "SELECT count(*) FROM store_appointments WHERE store_id = $1"
	if err := i.postgresDB.QueryRow(countSQL, storeId).Scan(&res.Total); err != nil {
		log.Println("An error occurred while counting store appointments", err)
		return nil, err
	}

	// página
	listSQL := `
		SELECT
			id,
			store_id,
			user_id,
			start_at,
			end_at,
			status,
			hold_expires_at,
			price,
			currency,
			fee_platform,
			payment_id,
			notes,
			created_at,
			updated_at
		FROM store_appointments
		WHERE store_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	offset := (page * limit) - limit

	rows, err := i.postgresDB.Query(listSQL, storeId, limit, offset)
	if err != nil {
		log.Println("An error occurred while getting store appointments", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var appointment StoreAppointment
		err := rows.Scan(
			&appointment.Id,
			&appointment.StoreId,
			&appointment.UserId,
			&appointment.StartAt,
			&appointment.EndAt,
			&appointment.Status,
			&appointment.HoldExpiresAt,
			&appointment.Price,
			&appointment.Currency,
			&appointment.FeePlatform,
			&appointment.PaymentId,
			&appointment.Notes,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
		if err != nil {
			log.Println("An error occurred while scanning store appointment", err)
			return nil, err
		}
		res.Appointments = append(res.Appointments, appointment)
	}
	if err := rows.Err(); err != nil {
		log.Println("Row iteration error while getting store appointments", err)
		return nil, err
	}

	return &res, nil
}

func (i *StoreRepo) UpdateStoreAppointment(id string, appointment StoreAppointment) (*StoreAppointment, error) {
	const sqlStmt = `
		UPDATE store_appointments
		SET start_at = $1, end_at = $2, status = $3, hold_expires_at = $4, price = $5, currency = $6, fee_platform = $7, payment_id = $8, notes = $9
		WHERE id = $10
	`
	_, err := i.postgresDB.Exec(sqlStmt,
		appointment.StartAt,
		appointment.EndAt,
		appointment.Status,
		appointment.HoldExpiresAt,
		appointment.Price,
		appointment.Currency,
		appointment.FeePlatform,
		appointment.PaymentId,
		appointment.Notes,
		id,
	)
	if err != nil {
		log.Println("An error occurred while updating store appointment", err)
		return nil, err
	}
	return &appointment, nil
}

func (i *StoreRepo) DeleteStoreAppointment(id string) error {
	const sqlStmt = `DELETE FROM store_appointments WHERE id = $1`
	if _, err := i.postgresDB.Exec(sqlStmt, id); err != nil {
		log.Println("An error occurred while deleting store appointment", err)
		return err
	}
	return nil
}

func (i *StoreRepo) getStoreStatements(storeId string) (string, error) {
	const sqlStmt = `
		SELECT
			stores.id,
			stores.name,
			stores.owner_id,
			stores.type,
			stores.location,
			store_availability.availability,
			store_ratings.user_id,
			store_ratings.rating,
			store_ratings.message,
			store_ratings.created_at
			store_plans.id as plan_id,
			store_plans.name as plan_name,
			store_plans.price as plan_price,
			store_plans.currency as plan_currency,
			store_plans.plan_type,
			store_plans.frequency
		FROM stores
		LEFT JOIN store_availability ON stores.id = store_availability.store_id
		LEFT JOIN store_ratings ON stores.id = store_ratings.store_id
		LEFT JOIN store_plans ON stores.id = store_plans.store_id
		WHERE stores.id = '` + storeId + `';
	`
	return sqlStmt, nil
}

func (i *StoreRepo) formatStore(row *sql.Rows) (*Store, error) {
	s := Store{}

	var location []byte
	err := row.Scan(
		&s.Id,
		&s.Name,
		&s.OwnerId,
		&s.Type,
		&location,
	)

	if err != nil {
		log.Println("An error occurred while scanning store", err)
		return nil, err
	}

	if location != nil {
		if err := json.Unmarshal(location, &s.Location); err != nil {
			log.Println("An error occurred while getting stores", err)
			return nil, err
		}
	}

	return &s, nil
}

func (i *StoreRepo) formatFullStore(row *sql.Rows) (*FullStore, error) {
	s := FullStore{
		Ratings: []StoreRating{},
		Plans:   []StorePlan{},
	}

	var location []byte
	err := row.Scan(
		&s.Store.Id,
		&s.Store.Name,
		&s.Store.OwnerId,
		&s.Store.Type,
		&location,
		&s.Availability.Availability,
		// Additional fields for Ratings and Plans would go here
	)

	if err != nil {
		log.Println("An error occurred while scanning full store", err)
		return nil, err
	}

	if location != nil {
		if err := json.Unmarshal(location, &s.Store.Location); err != nil {
			log.Println("An error occurred while unmarshalling store location", err)
			return nil, err
		}
	}

	return &s, nil
}

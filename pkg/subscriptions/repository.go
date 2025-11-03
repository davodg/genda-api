package subscriptions

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type SubscriptionRepo struct {
	postgresDB *sql.DB
}

func NewSubscriptionRepository(postgresDB *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{postgresDB: postgresDB}
}

func (i *SubscriptionRepo) CreateSubscription(subscription Subscription) (*Subscription, error) {
	if subscription.Id == "" {
		subscription.Id = uuid.New().String()
	}

	const sqlStmt = `
		INSERT INTO subscriptions (
			id,
			user_id,
			store_id,
			store_plan_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at,
			canceled_at,
			trial_end,
			provider,
			provider_sub_id,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
	`

	_, err := i.postgresDB.Exec(
		sqlStmt,
		subscription.Id,
		subscription.UserId,
		subscription.StoreId,
		subscription.StorePlanId,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAt,
		subscription.CanceledAt,
		subscription.TrialEnd,
		subscription.Provider,
		subscription.ProviderSubId,
	)
	if err != nil {
		log.Println("An error occurred while creating subscription", err)
		return nil, err
	}

	return &subscription, nil
}

func (i *SubscriptionRepo) GetSubscriptions(userId string, page int, limit int) (*GetSubscriptionsResponse, error) {
	const countStmt = `SELECT COUNT(*) FROM subscriptions WHERE user_id = $1`
	var total int
	err := i.postgresDB.QueryRow(countStmt, userId).Scan(&total)
	if err != nil {
		log.Println("An error occurred while counting subscriptions", err)
		return nil, err
	}

	const sqlStmt = `
		SELECT
			id,
			user_id,
			store_id,
			store_plan_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at,
			canceled_at,
			trial_end,
			provider,
			provider_sub_id,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := i.postgresDB.Query(sqlStmt, userId, limit, (page-1)*limit)
	if err != nil {
		log.Println("An error occurred while retrieving subscriptions", err)
		return nil, err
	}
	defer rows.Close()

	subscriptions := []Subscription{}
	for rows.Next() {
		var subscription Subscription
		err := rows.Scan(
			&subscription.Id,
			&subscription.UserId,
			&subscription.StoreId,
			&subscription.StorePlanId,
			&subscription.Status,
			&subscription.CurrentPeriodStart,
			&subscription.CurrentPeriodEnd,
			&subscription.CancelAt,
			&subscription.CanceledAt,
			&subscription.TrialEnd,
			&subscription.Provider,
			&subscription.ProviderSubId,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			log.Println("An error occurred while scanning subscription", err)
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	response := &GetSubscriptionsResponse{
		Total:         total,
		Limit:         limit,
		Subscriptions: subscriptions,
	}

	return response, nil
}

func (i *SubscriptionRepo) UpdateSubscription(id string, subscription Subscription) (*Subscription, error) {
	const sqlStmt = `
		UPDATE subscriptions
		SET
			user_id = $1,
			store_id = $2,
			store_plan_id = $3,
			status = $4,
			current_period_start = $5,
			current_period_end = $6,
			cancel_at = $7,
			canceled_at = $8,
			trial_end = $9,
			provider = $10,
			provider_sub_id = $11,
			updated_at = NOW()
		WHERE id = $12
	`

	_, err := i.postgresDB.Exec(
		sqlStmt,
		subscription.UserId,
		subscription.StoreId,
		subscription.StorePlanId,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAt,
		subscription.CanceledAt,
		subscription.TrialEnd,
		subscription.Provider,
		subscription.ProviderSubId,
		id,
	)
	if err != nil {
		log.Println("An error occurred while updating subscription", err)
		return nil, err
	}

	return &subscription, nil
}

func (i *SubscriptionRepo) DeleteSubscription(id string) error {
	const sqlStmt = `DELETE FROM subscriptions WHERE id = $1`
	if _, err := i.postgresDB.Exec(sqlStmt, id); err != nil {
		log.Println("An error occurred while deleting subscription", err)
		return err
	}
	return nil
}

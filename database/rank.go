package database

import (
	"fmt"
	"time"
)

func AddRankHistory(userID, rank string) error {
	db := GetDB()
	_, err := db.Exec(`
        INSERT INTO rank_history (user_id, rank) VALUES ($1, $2);
    `, userID, rank)
	if err != nil {
		return fmt.Errorf("error adding rank history: %v", err)
	}
	return nil
}

func GetRankHistory(userID string) ([]map[string]interface{}, error) {
	db := GetDB()
	rows, err := db.Query("SELECT rank, achieved_at FROM rank_history WHERE user_id = $1 ORDER BY achieved_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching rank history: %v", err)
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var rank string
		var achievedAt time.Time
		err := rows.Scan(&rank, &achievedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		history = append(history, map[string]interface{}{
			"rank":       rank,
			"achievedAt": achievedAt,
		})
	}

	return history, nil
}

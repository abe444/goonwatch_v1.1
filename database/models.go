package database

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID          string
	Streak      int
	CurrentRank string
}

type GoonerStatus struct {
	UserID        string
	LastStreak    int
	BestStreak    int
	RelapseCount  int
	LastRelapseAt time.Time
}

func GetUserStreak(userID string) (int, error) {
	db := GetDB()
	var streak int
	err := db.QueryRow("SELECT streak FROM users WHERE id = $1", userID).Scan(&streak)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error fetching user streak: %v", err)
	}
	return streak, nil
}

func UpdateStreak(userID string, streak int) error {
	db := GetDB()
	_, err := db.Exec(`
        INSERT INTO users (id, streak) VALUES ($1, $2)
        ON CONFLICT (id) DO UPDATE SET streak = $2;
    `, userID, streak)
	if err != nil {
		return fmt.Errorf("error updating streak: %v", err)
	}
	return nil
}

func GetLeaderboard() ([]User, error) {
	db := GetDB()
	rows, err := db.Query("SELECT id, streak FROM users ORDER BY streak DESC LIMIT 10")
	if err != nil {
		return nil, fmt.Errorf("error fetching leaderboard: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Streak)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func AddStreakHistory(userID string, streak int) error {
	db := GetDB()
	_, err := db.Exec(`
        INSERT INTO streak_history (user_id, streak) VALUES ($1, $2);
    `, userID, streak)
	if err != nil {
		return fmt.Errorf("error adding streak history: %v", err)
	}
	return nil
}

func GetStreakHistory(userID string) ([]map[string]interface{}, error) {
	db := GetDB()
	rows, err := db.Query("SELECT streak, recorded_at FROM streak_history WHERE user_id = $1 ORDER BY recorded_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching streak history: %v", err)
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var streak int
		var recordedAt time.Time
		err := rows.Scan(&streak, &recordedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		history = append(history, map[string]interface{}{
			"streak":     streak,
			"recordedAt": recordedAt,
		})
	}

	return history, nil
}

func UpdateUserRank(userID, rank string) error {
	db := GetDB()
	_, err := db.Exec(`
        UPDATE users SET current_rank = $1 
        WHERE id = $2
    `, rank, userID)
	if err != nil {
		return fmt.Errorf("error updating user rank: %v", err)
	}
	return nil
}

func AddToHallOfGooners(userID string, streak int) error {
	db := GetDB()
	_, err := db.Exec(`
        INSERT INTO hall_of_gooners (user_id, last_streak, best_streak)
        VALUES ($1, $2, $2)
        ON CONFLICT (user_id) DO UPDATE SET
            last_streak = $2,
            best_streak = GREATEST(hall_of_gooners.best_streak, $2),
            relapse_count = hall_of_gooners.relapse_count + 1,
            last_relapse_at = CURRENT_TIMESTAMP
    `, userID, streak)
	return err
}

func RemoveFromHallOfGooners(userID string) error {
	db := GetDB()
	_, err := db.Exec("DELETE FROM hall_of_gooners WHERE user_id = $1", userID)
	return err
}

func GetGoonerStatus(userID string) (*GoonerStatus, error) {
	db := GetDB()
	status := &GoonerStatus{}
	err := db.QueryRow(`
        SELECT user_id, last_streak, best_streak, relapse_count, last_relapse_at 
        FROM hall_of_gooners 
        WHERE user_id = $1
    `, userID).Scan(&status.UserID, &status.LastStreak, &status.BestStreak,
		&status.RelapseCount, &status.LastRelapseAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return status, err
}

func UpdateLastLogTime(userID string) error {
	db := GetDB()
	_, err := db.Exec(`
        UPDATE users 
        SET last_log_time = CURRENT_TIMESTAMP 
        WHERE id = $1
    `, userID)
	return err
}

func CanLogAgain(userID string) (bool, error) {
	db := GetDB()
	var lastLogTime sql.NullTime
	err := db.QueryRow(`
        SELECT last_log_time 
        FROM users 
        WHERE id = $1
    `, userID).Scan(&lastLogTime)

	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	if !lastLogTime.Valid {
		return true, nil
	}

	cooldownDuration := 24 * time.Hour
	timeRemaining := time.Until(lastLogTime.Time.Add(cooldownDuration))

	return timeRemaining <= 0, nil
}

func GetTimeUntilNextLog(userID string) (time.Duration, error) {
	db := GetDB()
	var lastLogTime sql.NullTime
	err := db.QueryRow(`
        SELECT last_log_time 
        FROM users 
        WHERE id = $1
    `, userID).Scan(&lastLogTime)

	if err == sql.ErrNoRows || !lastLogTime.Valid {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	if !lastLogTime.Valid {
		return 0, nil
	}

	cooldownDuration := 24 * time.Hour
	timeRemaining := time.Until(lastLogTime.Time.Add(cooldownDuration))

	return timeRemaining, nil
}

package pg

import (
	"fmt"
)

// --- Structs ---

// FriendUser represents a user's public profile, suitable for JSON responses.
type FriendUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// IncomingRequest represents a pending friend request with sender's info.
type IncomingRequest struct {
	RequestID int64      `json:"request_id"`
	Sender    FriendUser `json:"sender"`
}

// --- Database Functions ---

// CreateFriendRequest creates a new pending friendship request.
func CreateFriendRequest(senderID, receiverID int64) error {
	if senderID == receiverID {
		return fmt.Errorf("user cannot send a friend request to themselves")
	}

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM friendships 
			WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
		)`
	err := DB.QueryRow(query, senderID, receiverID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for existing friendship: %w", err)
	}
	if exists {
		return fmt.Errorf("a friend request or friendship already exists between these users")
	}

	insertQuery := "INSERT INTO friendships (user_id, friend_id, status) VALUES ($1, $2, 'pending')"
	_, err = DB.Exec(insertQuery, senderID, receiverID)
	if err != nil {
		return fmt.Errorf("failed to create friend request: %w", err)
	}

	return nil
}

// GetFriendsByUserID retrieves a list of accepted friends for a given user.
func GetFriendsByUserID(userID int64) ([]FriendUser, error) {
	query := `
		SELECT u.id, u.username
		FROM users u
		JOIN (
			SELECT friend_id AS id FROM friendships
			WHERE user_id = $1 AND status = 'accepted'
			UNION
			SELECT user_id AS id FROM friendships
			WHERE friend_id = $1 AND status = 'accepted'
		) AS friends ON u.id = friends.id
		ORDER BY u.username;
	`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query friends: %w", err)
	}
	defer rows.Close()

	var friendsList []FriendUser
	for rows.Next() {
		var friend FriendUser
		if err := rows.Scan(&friend.ID, &friend.Username); err != nil {
			return nil, fmt.Errorf("failed to scan friend row: %w", err)
		}
		friendsList = append(friendsList, friend)
	}

	return friendsList, rows.Err()
}

// DeleteFriendship removes an accepted friendship between two users.
func DeleteFriendship(userID1, userID2 int64) error {
	query := `
		DELETE FROM friendships
		WHERE 
			(user_id = $1 AND friend_id = $2 AND status = 'accepted') OR
			(user_id = $2 AND friend_id = $1 AND status = 'accepted')
	`
	result, err := DB.Exec(query, userID1, userID2)
	if err != nil {
		return fmt.Errorf("failed to execute delete friendship query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no accepted friendship found between users to delete")
	}

	return nil
}

// GetIncomingFriendRequests retrieves all pending friend requests for a user.
func GetIncomingFriendRequests(receiverID int64) ([]IncomingRequest, error) {
	query := `
		SELECT f.id, u.id, u.username
		FROM friendships f
		JOIN users u ON f.user_id = u.id
		WHERE f.friend_id = $1 AND f.status = 'pending';
	`
	rows, err := DB.Query(query, receiverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query incoming friend requests: %w", err)
	}
	defer rows.Close()

	var requests []IncomingRequest
	for rows.Next() {
		var req IncomingRequest
		if err := rows.Scan(&req.RequestID, &req.Sender.ID, &req.Sender.Username); err != nil {
			return nil, fmt.Errorf("failed to scan incoming friend request: %w", err)
		}
		requests = append(requests, req)
	}

	return requests, rows.Err()
}

// UpdateFriendRequestStatus updates the status of a request ('accepted' or 'rejected').
// It ensures that only the intended recipient of the request can update it.
func UpdateFriendRequestStatus(requestID, receiverID int64, newStatus string) error {
	if newStatus != "accepted" && newStatus != "rejected" {
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	// If accepting, we create a two-way friendship by creating the inverse relationship
	if newStatus == "accepted" {
		tx, err := DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Update the original request
		updateQuery := `UPDATE friendships SET status = 'accepted' WHERE id = $1 AND friend_id = $2 AND status = 'pending'`
		result, err := tx.Exec(updateQuery, requestID, receiverID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to accept friend request: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get affected rows for update: %w", err)
		}
		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("no pending friend request found with the specified ID for this user")
		}

		return tx.Commit()
	}

	// If rejecting, just delete the request
	deleteQuery := `DELETE FROM friendships WHERE id = $1 AND friend_id = $2 AND status = 'pending'`
	result, err := DB.Exec(deleteQuery, requestID, receiverID)
	if err != nil {
		return fmt.Errorf("failed to reject friend request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows for delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no pending friend request found with the specified ID for this user to reject")
	}

	return nil
}
package types

// UserSearchFilters represents search filters for user queries
type UserSearchFilters struct {
	Username string
	Email    string
	Status   *int // pointer to allow nil (no filter)
	RoleID   uint
}

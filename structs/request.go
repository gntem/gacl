package structs

// PaginationQuery struct
type PaginationQuery struct {
	Page   int64  `validate:"gte=0" form:"page,default=1" binding:"required"`
	Limit  int64  `validate:"gte=0" form:"limit,default=10" binding:"required"`
	SortBy string `validate:"oneof=created_at id" form:"sortBy,default=ID" binding:"required"`
	Order  string `validate:"oneof=desc asc" form:"order,default=asc" binding:"required"`
}

// UserCreateRequest struct
type UserCreateRequest struct {
	Name string `form:"name" validate:"min=4,max=255" binding:"required"`
}

// GroupCreateRequest struct
type GroupCreateRequest struct {
	Name string `form:"name" validate:"min=4,max=255" binding:"required"`
}

// UserUpdateRequest struct
type UserUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

// GroupUpdateRequest struct
type GroupUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

// PermissionUpdateRequest struct
type PermissionUpdateRequest struct {
	Name string `form:"name" validate:"min=4,max=255"`
}

// PermissionCreateRequest struct
type PermissionCreateRequest struct {
	Name string `form:"name" validate:"min=4,max=255" binding:"required"`
}

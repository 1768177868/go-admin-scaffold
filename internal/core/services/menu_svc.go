package services

import (
	"context"
	"errors"

	"app/internal/core/models"
)

var (
	ErrMenuNotFound    = errors.New("menu not found")
	ErrMenuHasChildren = errors.New("menu has children, cannot delete")
)

type MenuRepository interface {
	FindByID(ctx context.Context, id uint) (*models.Menu, error)
	FindAll(ctx context.Context) ([]models.Menu, error)
	FindByParentID(ctx context.Context, parentID *uint) ([]models.Menu, error)
	FindTree(ctx context.Context) ([]models.Menu, error)
	FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Menu, error)
	FindVisibleMenus(ctx context.Context) ([]models.Menu, error)
	Create(ctx context.Context, menu *models.Menu) error
	Update(ctx context.Context, menu *models.Menu) error
	Delete(ctx context.Context, id uint) error
	UpdateMenuRoles(ctx context.Context, menuID uint, roleIDs []uint) error
	GetMaxSort(ctx context.Context, parentID *uint) (int, error)
}

type MenuService struct {
	menuRepo MenuRepository
	userRepo UserRepository
}

func NewMenuService(menuRepo MenuRepository, userRepo UserRepository) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
		userRepo: userRepo,
	}
}

type CreateMenuRequest struct {
	Name       string          `json:"name" binding:"required"`
	Title      string          `json:"title" binding:"required"`
	Icon       string          `json:"icon"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	ParentID   *uint           `json:"parent_id"`
	Sort       int             `json:"sort"`
	Type       int             `json:"type"`
	Visible    int             `json:"visible"`
	Status     int             `json:"status"`
	KeepAlive  bool            `json:"keep_alive"`
	External   bool            `json:"external"`
	Permission string          `json:"permission"`
	Meta       models.MenuMeta `json:"meta"`
	RoleIDs    []uint          `json:"role_ids"`
}

type UpdateMenuRequest struct {
	Name       string          `json:"name"`
	Title      string          `json:"title"`
	Icon       string          `json:"icon"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	ParentID   *uint           `json:"parent_id"`
	Sort       int             `json:"sort"`
	Type       int             `json:"type"`
	Visible    int             `json:"visible"`
	Status     int             `json:"status"`
	KeepAlive  bool            `json:"keep_alive"`
	External   bool            `json:"external"`
	Permission string          `json:"permission"`
	Meta       models.MenuMeta `json:"meta"`
	RoleIDs    []uint          `json:"role_ids"`
}

// MenuRouteItem represents a menu item for frontend routing
type MenuRouteItem struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	Path      string          `json:"path"`
	Component string          `json:"component"`
	Meta      models.MenuMeta `json:"meta"`
	Children  []MenuRouteItem `json:"children,omitempty"`
}

// Create creates a new menu
func (s *MenuService) Create(ctx context.Context, req *CreateMenuRequest) (*models.Menu, error) {
	// Set sort value if not provided
	if req.Sort == 0 {
		maxSort, err := s.menuRepo.GetMaxSort(ctx, req.ParentID)
		if err != nil {
			return nil, err
		}
		req.Sort = maxSort + 1
	}

	menu := &models.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		ParentID:   req.ParentID,
		Sort:       req.Sort,
		Type:       req.Type,
		Visible:    req.Visible,
		Status:     req.Status,
		KeepAlive:  req.KeepAlive,
		External:   req.External,
		Permission: req.Permission,
		Meta:       req.Meta,
	}

	if err := s.menuRepo.Create(ctx, menu); err != nil {
		return nil, err
	}

	// Update role associations if provided
	if len(req.RoleIDs) > 0 {
		if err := s.menuRepo.UpdateMenuRoles(ctx, menu.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	return s.menuRepo.FindByID(ctx, menu.ID)
}

// Update updates an existing menu
func (s *MenuService) Update(ctx context.Context, id uint, req *UpdateMenuRequest) (*models.Menu, error) {
	menu, err := s.menuRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrMenuNotFound
	}

	// Update fields
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Title != "" {
		menu.Title = req.Title
	}
	menu.Icon = req.Icon
	menu.Path = req.Path
	menu.Component = req.Component
	menu.ParentID = req.ParentID
	if req.Sort != 0 {
		menu.Sort = req.Sort
	}
	if req.Type != 0 {
		menu.Type = req.Type
	}
	if req.Visible != 0 {
		menu.Visible = req.Visible
	}
	if req.Status != 0 {
		menu.Status = req.Status
	}
	menu.KeepAlive = req.KeepAlive
	menu.External = req.External
	menu.Permission = req.Permission
	menu.Meta = req.Meta

	if err := s.menuRepo.Update(ctx, menu); err != nil {
		return nil, err
	}

	// Update role associations if provided
	if req.RoleIDs != nil {
		if err := s.menuRepo.UpdateMenuRoles(ctx, menu.ID, req.RoleIDs); err != nil {
			return nil, err
		}
	}

	return s.menuRepo.FindByID(ctx, menu.ID)
}

// Delete deletes a menu
func (s *MenuService) Delete(ctx context.Context, id uint) error {
	// Check if menu has children
	children, err := s.menuRepo.FindByParentID(ctx, &id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrMenuHasChildren
	}

	return s.menuRepo.Delete(ctx, id)
}

// GetByID gets a menu by ID
func (s *MenuService) GetByID(ctx context.Context, id uint) (*models.Menu, error) {
	return s.menuRepo.FindByID(ctx, id)
}

// GetTree gets the menu tree
func (s *MenuService) GetTree(ctx context.Context) ([]models.Menu, error) {
	return s.menuRepo.FindTree(ctx)
}

// GetAll gets all menus
func (s *MenuService) GetAll(ctx context.Context) ([]models.Menu, error) {
	return s.menuRepo.FindAll(ctx)
}

// GetUserMenus gets menus accessible by a user
func (s *MenuService) GetUserMenus(ctx context.Context, userID uint) ([]MenuRouteItem, error) {
	// Get user roles
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Extract role IDs
	roleIDs := make([]uint, len(user.Roles))
	for i, role := range user.Roles {
		roleIDs[i] = role.ID
	}

	// Get menus by role IDs
	menus, err := s.menuRepo.FindByRoleIDs(ctx, roleIDs)
	if err != nil {
		return nil, err
	}

	// Filter visible menus
	var visibleMenus []models.Menu
	for _, menu := range menus {
		if menu.IsVisible() && menu.IsMenu() {
			visibleMenus = append(visibleMenus, menu)
		}
	}

	// Build tree and convert to route items
	return s.buildMenuRoutes(visibleMenus), nil
}

// GetVisibleMenuTree gets the visible menu tree for public access
func (s *MenuService) GetVisibleMenuTree(ctx context.Context) ([]MenuRouteItem, error) {
	menus, err := s.menuRepo.FindVisibleMenus(ctx)
	if err != nil {
		return nil, err
	}

	// Filter menu type only
	var menuItems []models.Menu
	for _, menu := range menus {
		if menu.IsMenu() {
			menuItems = append(menuItems, menu)
		}
	}

	return s.buildMenuRoutes(menuItems), nil
}

// buildMenuRoutes builds menu route items from menu models
func (s *MenuService) buildMenuRoutes(menus []models.Menu) []MenuRouteItem {
	menuMap := make(map[uint]*MenuRouteItem)
	var rootMenus []MenuRouteItem

	// Create menu items
	for _, menu := range menus {
		item := &MenuRouteItem{
			ID:        menu.ID,
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Meta:      menu.Meta,
			Children:  make([]MenuRouteItem, 0),
		}
		menuMap[menu.ID] = item
	}

	// Build tree structure
	for _, menu := range menus {
		item := menuMap[menu.ID]
		if menu.ParentID == nil {
			// Root menu
			rootMenus = append(rootMenus, *item)
		} else {
			// Child menu
			parent, exists := menuMap[*menu.ParentID]
			if exists {
				parent.Children = append(parent.Children, *item)
			}
		}
	}

	return rootMenus
}

// UpdateMenuRoles updates the roles associated with a menu
func (s *MenuService) UpdateMenuRoles(ctx context.Context, menuID uint, roleIDs []uint) error {
	// Check if menu exists
	_, err := s.menuRepo.FindByID(ctx, menuID)
	if err != nil {
		return ErrMenuNotFound
	}

	return s.menuRepo.UpdateMenuRoles(ctx, menuID, roleIDs)
}

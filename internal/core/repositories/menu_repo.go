package repositories

import (
	"context"

	"app/internal/core/models"

	"gorm.io/gorm"
)

type MenuRepository struct {
	*BaseRepository
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// FindByID retrieves a menu by ID with its children
func (r *MenuRepository) FindByID(ctx context.Context, id uint) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).
		Preload("Children").
		Where("id = ?", id).
		First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// FindAll retrieves all menus with their relationships
func (r *MenuRepository) FindAll(ctx context.Context) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Order("sort ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// FindByParentID retrieves menus by parent ID
func (r *MenuRepository) FindByParentID(ctx context.Context, parentID *uint) ([]models.Menu, error) {
	var menus []models.Menu
	query := r.db.WithContext(ctx).Order("sort ASC, id ASC")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Find(&menus).Error
	return menus, err
}

// FindTree retrieves menu tree structure
func (r *MenuRepository) FindTree(ctx context.Context) ([]models.Menu, error) {
	// First get all menus
	var allMenus []models.Menu
	err := r.db.WithContext(ctx).
		Order("sort ASC, id ASC").
		Find(&allMenus).Error
	if err != nil {
		return nil, err
	}

	// Build tree structure
	menuMap := make(map[uint]*models.Menu)
	var rootMenus []models.Menu

	// First pass: create map of all menus
	for i := range allMenus {
		menu := &allMenus[i]
		menu.Children = make([]*models.Menu, 0)
		menuMap[menu.ID] = menu
	}

	// Second pass: build tree
	for _, menu := range allMenus {
		if menu.ParentID == nil {
			// Root menu
			rootMenus = append(rootMenus, *menuMap[menu.ID])
		} else {
			// Child menu
			parent, exists := menuMap[*menu.ParentID]
			if exists {
				parent.Children = append(parent.Children, menuMap[menu.ID])
			}
		}
	}

	return rootMenus, nil
}

// FindByRoleIDs retrieves menus accessible by given role IDs
func (r *MenuRepository) FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.WithContext(ctx).
		Joins("LEFT JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id IN ? OR menus.id IN (SELECT id FROM menus WHERE parent_id IS NULL)", roleIDs).
		Group("menus.id").
		Order("sort ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// FindVisibleMenus retrieves all visible and enabled menus
func (r *MenuRepository) FindVisibleMenus(ctx context.Context) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.WithContext(ctx).
		Where("visible = 1 AND status = 1").
		Order("sort ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// Create creates a new menu
func (r *MenuRepository) Create(ctx context.Context, menu *models.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

// Update updates an existing menu
func (r *MenuRepository) Update(ctx context.Context, menu *models.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

// Delete deletes a menu by ID
func (r *MenuRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete role-menu associations
		if err := tx.Where("menu_id = ?", id).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}

		// Delete menu
		return tx.Delete(&models.Menu{}, id).Error
	})
}

// UpdateMenuRoles updates the roles associated with a menu
func (r *MenuRepository) UpdateMenuRoles(ctx context.Context, menuID uint, roleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove existing associations
		if err := tx.Where("menu_id = ?", menuID).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}

		// Add new associations
		if len(roleIDs) > 0 {
			roleMenus := make([]models.RoleMenu, 0, len(roleIDs))
			for _, roleID := range roleIDs {
				roleMenus = append(roleMenus, models.RoleMenu{
					MenuID: menuID,
					RoleID: roleID,
				})
			}
			return tx.Create(&roleMenus).Error
		}

		return nil
	})
}

// GetMaxSort returns the maximum sort value for a given parent
func (r *MenuRepository) GetMaxSort(ctx context.Context, parentID *uint) (int, error) {
	var maxSort int
	query := r.db.WithContext(ctx).Model(&models.Menu{}).Select("COALESCE(MAX(sort), 0)")

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Scan(&maxSort).Error
	return maxSort, err
}

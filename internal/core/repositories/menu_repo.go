package repositories

import (
	"context"
	"fmt"
	"log"

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
		Order("CASE WHEN parent_id IS NULL THEN 0 ELSE 1 END, sort ASC, id ASC").
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
	log.Printf("[DEBUG] ========== FindTree called ==========")

	// First get all menus
	var allMenus []models.Menu
	err := r.db.WithContext(ctx).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort ASC, id ASC")
		}).
		Order("CASE WHEN parent_id IS NULL THEN 0 ELSE 1 END, sort ASC, id ASC").
		Find(&allMenus).Error
	if err != nil {
		log.Printf("[ERROR] FindTree: Failed to query all menus: %v", err)
		return nil, err
	}

	log.Printf("[DEBUG] FindTree: Found %d total menus from database", len(allMenus))
	for i, menu := range allMenus {
		parentIDStr := "NULL"
		if menu.ParentID != nil {
			parentIDStr = fmt.Sprintf("%d", *menu.ParentID)
		}
		log.Printf("[DEBUG] FindTree: Menu %d - ID=%d, Name=%s, Title=%s, ParentID=%s, Visible=%d, Status=%d",
			i+1, menu.ID, menu.Name, menu.Title, parentIDStr, menu.Visible, menu.Status)
	}

	// Build tree structure
	menuMap := make(map[uint]*models.Menu)
	var rootMenus []models.Menu

	// First pass: create map of all menus
	log.Printf("[DEBUG] FindTree: First pass - creating menu map...")
	for i := range allMenus {
		menu := &allMenus[i]
		menuMap[menu.ID] = menu
		log.Printf("[DEBUG] FindTree: Added menu ID=%d (%s) to map", menu.ID, menu.Title)
	}
	log.Printf("[DEBUG] FindTree: Menu map contains %d items", len(menuMap))

	// Second pass: build tree
	log.Printf("[DEBUG] FindTree: Second pass - building tree structure...")
	for _, menu := range allMenus {
		if menu.ParentID == nil {
			// Root menu
			rootMenus = append(rootMenus, menu)
			log.Printf("[DEBUG] FindTree: ✅ Added ID=%d (%s) as ROOT menu", menu.ID, menu.Title)
		}
	}

	log.Printf("[DEBUG] FindTree: Final tree structure has %d root menus", len(rootMenus))
	for i, root := range rootMenus {
		log.Printf("[DEBUG] FindTree: Root menu %d - ID=%d, Title=%s, Children=%d",
			i+1, root.ID, root.Title, len(root.Children))
		for j, child := range root.Children {
			log.Printf("[DEBUG] FindTree:   Child %d.%d - ID=%d, Title=%s",
				i+1, j+1, child.ID, child.Title)
		}
	}
	log.Printf("[DEBUG] ========== FindTree completed ==========")

	return rootMenus, nil
}

// FindByRoleIDs retrieves menus accessible by given role IDs
func (r *MenuRepository) FindByRoleIDs(ctx context.Context, roleIDs []uint) ([]models.Menu, error) {
	log.Printf("[DEBUG] FindByRoleIDs called with role IDs: %v", roleIDs)

	// First check if there are any role-menu associations
	var count int64
	err := r.db.WithContext(ctx).
		Table("role_menus").
		Where("role_id IN ?", roleIDs).
		Count(&count).Error
	if err != nil {
		log.Printf("[ERROR] Failed to count role_menus for role IDs %v: %v", roleIDs, err)
		return nil, err
	}
	log.Printf("[DEBUG] Found %d role_menu associations for role IDs %v", count, roleIDs)

	// Get menu IDs that are associated with the given role IDs
	var menuIDs []uint
	err = r.db.WithContext(ctx).
		Table("role_menus").
		Where("role_id IN ?", roleIDs).
		Distinct("menu_id").
		Pluck("menu_id", &menuIDs).Error
	if err != nil {
		log.Printf("[ERROR] Failed to get menu IDs for role IDs %v: %v", roleIDs, err)
		return nil, err
	}
	log.Printf("[DEBUG] Found menu IDs: %v for role IDs %v", menuIDs, roleIDs)

	if len(menuIDs) == 0 {
		log.Printf("[DEBUG] No menu IDs found for role IDs %v", roleIDs)
		return []models.Menu{}, nil
	}

	// Get all menus (including parents) for the found menu IDs
	var allMenus []models.Menu
	err = r.db.WithContext(ctx).
		Where("id IN ?", menuIDs).
		Order("sort ASC, id ASC").
		Find(&allMenus).Error
	if err != nil {
		log.Printf("[ERROR] Failed to find menus for menu IDs %v: %v", menuIDs, err)
		return nil, err
	}
	log.Printf("[DEBUG] Found %d direct menus for menu IDs %v", len(allMenus), menuIDs)

	// Collect parent menu IDs
	var parentIDs []uint
	parentIDMap := make(map[uint]bool)
	for _, menu := range allMenus {
		if menu.ParentID != nil && !parentIDMap[*menu.ParentID] {
			parentIDs = append(parentIDs, *menu.ParentID)
			parentIDMap[*menu.ParentID] = true
		}
	}
	log.Printf("[DEBUG] Found parent menu IDs: %v", parentIDs)

	// Get parent menus if any
	if len(parentIDs) > 0 {
		var parentMenus []models.Menu
		err = r.db.WithContext(ctx).
			Where("id IN ?", parentIDs).
			Order("sort ASC, id ASC").
			Find(&parentMenus).Error
		if err != nil {
			log.Printf("[ERROR] Failed to find parent menus for IDs %v: %v", parentIDs, err)
			return nil, err
		}
		log.Printf("[DEBUG] Found %d parent menus", len(parentMenus))
		allMenus = append(allMenus, parentMenus...)
	}

	log.Printf("[DEBUG] Total menus (including parents): %d", len(allMenus))
	return allMenus, nil
}

// FindVisibleMenus retrieves all visible and enabled menus
func (r *MenuRepository) FindVisibleMenus(ctx context.Context) ([]models.Menu, error) {
	log.Printf("[DEBUG] FindVisibleMenus called")

	// First, check if menus table has any data
	var totalCount int64
	if err := r.db.WithContext(ctx).Model(&models.Menu{}).Count(&totalCount).Error; err != nil {
		log.Printf("[ERROR] Failed to count total menus: %v", err)
		return nil, err
	}
	log.Printf("[DEBUG] Total menus in database: %d", totalCount)

	var menus []models.Menu
	err := r.db.WithContext(ctx).
		Where("visible = 1 AND status = 1").
		Order("sort ASC, id ASC").
		Find(&menus).Error

	log.Printf("[DEBUG] Found %d visible and enabled menus out of %d total", len(menus), totalCount)
	for i, menu := range menus {
		log.Printf("[DEBUG] Menu %d: ID=%d, Name=%s, Title=%s, Visible=%d, Status=%d",
			i+1, menu.ID, menu.Name, menu.Title, menu.Visible, menu.Status)
	}

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

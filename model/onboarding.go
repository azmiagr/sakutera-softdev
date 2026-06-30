package model

type WorkPlatformItem struct {
	WorkPlatformID string `json:"work_platform_id"`
	Name           string `json:"name"`
}

type WorkCategoryItem struct {
	WorkCategoryID string             `json:"work_category_id"`
	Name           string             `json:"name"`
	Platforms      []WorkPlatformItem `json:"platforms"`
}

type GetWorkCategoriesResponse struct {
	Categories []WorkCategoryItem `json:"categories"`
}

type SelectWorkPlatformRequest struct {
	WorkPlatformID string `json:"work_platform_id" binding:"required"`
}

type SelectWorkPlatformResponse struct {
	Message string `json:"message"`
}

type IncomeSourceItem struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Description string `json:"description"`
	IsAvailable bool   `json:"is_available"`
}

type GetIncomeSourcesResponse struct {
	Sources []IncomeSourceItem `json:"sources"`
}

type SelectIncomeSourceRequest struct {
	IncomeSourceType string `json:"income_source_type" binding:"required,oneof=manual"`
}

type SelectIncomeSourceResponse struct {
	Message string `json:"message"`
}

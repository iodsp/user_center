package params

type RoleResource struct {
	RoleId      int    `json:"RoleId"`
	ResourceId  int    `json:"ResourceId"`
	ResourceUrl string `json:"ResourceUrl"`
}

type RoleItem struct {
	RoleId      int    `json:"RoleId"`
	ResourceId  int    `json:"ResourceId"`
	ResourceUrl string `json:"ResourceUrl"`
	RoleName    string `json:"RoleName"`
}

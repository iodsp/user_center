package params

type UserRoleParam struct {
	RoleId     int    `json:"RoleId"`
	UserId     int    `json:"UserId"`
	RoleName   string `json:"RoleName"`
	DomainId   int    `json:"DomainId"`
	DomainName string `json:"DomainName"`
}

type UserRoleItem struct {
	RoleId     int    `json:"RoleId"`
	UserId     int    `json:"UserId"`
	UserName   string `json:"UserName"`
	Phone      string `json:"Phone"`
	RoleName   string `json:"RoleName"`
	DomainId   int    `json:"DomainId"`
	DomainName string `json:"RDomainName"`
}

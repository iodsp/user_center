package params

type UserResource struct {
	UserId      int    `json:"UserId"`
	ResourceUrl string `json:"ResourceUrl"`
}

type UserResourceItem struct {
	UserId       int    `json:"UserId"`
	ResourceUrl  string `json:"ResourceUrl"`
	ResourceName string `json:"ResourceName"`
	UserName     string `json:"UserName"`
}

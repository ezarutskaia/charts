package issue

type  DTOissue struct {
	Title     string    `json:"title"`
	UserID    uint      `json:"user_id"`
	ProjectID uint      `json:"project_id"`
	Priority  int       `json:"priority"`
	Status    string    `json:"status"`
	Deadline  string 	`json:"deadline"`
	Watchers  []uint    `json:"watchers"`
}
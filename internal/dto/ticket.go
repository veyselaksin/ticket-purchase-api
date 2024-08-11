package dto

type TicketCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

type TicketResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

type TicketPurchaseRequest struct {
	TicketId string `json:"-"`
	UserId   string `json:"user_id"`
	Quantity int    `json:"quantity"`
}

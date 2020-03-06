package payment

type Payment struct {
	Id          int     `json:"id"`
	Description string  `json:"name"`
	Amount      float64 `json:"amount"`
	Sender      string  `json:"sender"`
}

func NewPayment(id int, description, sender string, amount float64) Payment {
	return Payment{
		Id:          id,
		Description: description,
		Amount:      amount,
		Sender:      sender,
	}
}

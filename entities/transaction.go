package entities

type Transaction struct {
	ID              string  `json:"id"`
	UserDocument    string  `json:"cpf"`
	CreditCardToken string  `json:"creditCardToken"`
	Value           float64 `json:"value"`
}

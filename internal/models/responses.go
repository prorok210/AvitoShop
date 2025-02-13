package models

type Error400Response struct {
	Error string `json:"error" example:"Неверный запрос"`
}

type Error401Response struct {
	Error string `json:"error" example:"Неавторизован"`
}

type Error404Response struct {
	Error string `json:"error" example:"Не найдено"`
}

type Error500Response struct {
	Error string `json:"error" example:"Внутренняя ошибка сервера"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Успешный запрос"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type InfoResponse struct {
	Coins       int             `json:"coins" example:"100" description:"Количество доступных монет."`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type" example:"weapon" description:"Тип предмета."`
	Quantity int    `json:"quantity" example:"2" description:"Количество предметов."`
}

type CoinHistory struct {
	Received []ReceivedTx `json:"received"`
	Sent     []SentTx     `json:"sent"`
}

type ReceivedTx struct {
	FromUser string `json:"fromUser" example:"Alice" description:"Имя пользователя, который отправил монеты."`
	Amount   int    `json:"amount" example:"50" description:"Количество полученных монет."`
}

type SentTx struct {
	ToUser string `json:"toUser" example:"Bob" description:"Имя пользователя, которому отправлены монеты."`
	Amount int    `json:"amount" example:"30" description:"Количество отправленных монет."`
}

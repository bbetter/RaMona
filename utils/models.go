package utils

type TChat struct {
	Id int `json:"id"`
}

type TMessage struct {
	Chat TChat `json:"chat"`
}

type TResult struct {
	Message  TMessage `json:"message"`
	UpdateId int      `json:"update_id"`
}

type TUpdate struct {
	Result []TResult `json:"result"`
}

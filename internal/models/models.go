package models

/*
Здесь находятся общие структуры, которые обеспечивают единый формат передачи данных и могут применяться
в разных частях приложения.
*/

type RespSucc struct { // Структура для успешных ответов.
	Status string        `json:"status"`         // Статус ответа = "success"
	Code   int           `json:"code"`           // HTTP-код
	Body   *RespSuccData `json:"body,omitempty"` // Тело ответа
}

type RespSuccData struct {
	Data interface{} `json:"data,omitempty"`
}

type RespErr struct { // Структура для ошибочных ответов.
	Status  string     `json:"status"`            // Статус ответа = "error"
	Code    int        `json:"code"`              // HTTP-код ошибки
	Message RespErrMsg `json:"message,omitempty"` // Сообщение об ошибке
}

type RespErrMsg struct {
	Client string `json:"client,omitempty"` // Ошибка для клиента (на русском языке)
	Detail string `json:"detail,omitempty"` // Ошибка для разработчика
}

type Err struct { // Внутренняя структура ошибок.
	Code      int
	ClientMsg string
	Error     error
}

type User struct { // Общая структура пользователя.
	Id       int
	Username string
	Email    string
	Password string
}

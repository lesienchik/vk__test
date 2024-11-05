package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Проверяет username на корректность написания.
func IsValidUsername(username string) bool {
	if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 24 {
		return false
	}

	regex := `^[a-zA-Zа-яА-Я0-9]+$` // Допускаем только русские/английские символы + цифры
	re := regexp.MustCompile(regex)
	return re.MatchString(username)
}

// Проверяет почту на корректность написания.
func IsValidEmail(email string) bool {
	// Удаляем пробелы по краям.
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Проверяем, что адрес не содержит пробелов.
	if strings.Contains(email, " ") {
		return false
	}

	// Регулярное выражение для валидации email:
	// 1. Должен содержать символ @
	// 2. После @ должна быть хотя бы одна буква и точка
	// 3. Должен быть хотя бы один символ перед @
	// 4. Не допускается подряд идущие точки
	// 5. Разрешены только латинские символы, цифры и точки в имени
	// 6. Разрешен символы перед точкой и после, но не точки в начале и конце
	emailRegex := `^[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return false
	}

	// Проверка на наличие недопустимых символов.
	for _, char := range email {
		if char < '!' || char > '~' { // Все символы должны быть в диапазоне ASCII
			return false
		}
	}
	return true
}

// Проверяет валидность пароля по заданным критериям.
func IsValidPassword(password string) bool {
	if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 24 {
		return false
	}

	// Запрещаем пробелы и табуляции.
	if regexp.MustCompile(`[\s\t]`).MatchString(password) {
		return false
	}

	// Разрешенные символы: русские, английские буквы, цифры, знаки.
	validChars := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9!@#$%^&*()_+=\-\[\]{};:'",.<>?/|\\~` + "`" + `]+$`)
	return validChars.MatchString(password)
}

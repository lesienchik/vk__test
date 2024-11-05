package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidUsername(t *testing.T) {
	// Arrange
	requires := require.New(t)

	testTable := []struct {
		desc     string // Описание теста
		input    string // Входные данные
		expected bool   // Ожидаемый результат выполнения теста
	}{
		{
			desc:     "Success",
			input:    "ThisIsMyUsername", // Допустимы только русские и английские буквы + цифры
			expected: true,
		},
		{
			desc:     "Success",
			input:    "Пользователь",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "User",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "Польз",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "User2001",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "Пользователь2001",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "Use1",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "П0льзователь",
			expected: true,
		},
		{
			desc:     "Fail",
			input:    "Us", // Не короче трех символов
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "По", // Не короче трех символов
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "Us.er", // Никаких точек, запятых или иных знаков - только текст и цифры
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "По,льзователь", // Никаких точек, запятых или иных знаков - только текст и цифры
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "",
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "UserUserUserUserUserUserUserUserUserUserUserUser", // Не больше 24-ех символов
			expected: false,
		}, {
			desc:     "Fail",
			input:    "ПользовательПользовательПользовательПользовательПользовательПользователь", // Не больше 24-ех символов
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "Use@r", // Никаких символов - только текст и цифры
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "Пользов@тель", // Никаких символов - только текст и цифры
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "Tͦ̎h͊̚ȅ̐    Nͧͥ e͂̎  zͩ͗ p̂̔e", // Никаких символов - только текст и цифры
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้", // Никаких символов - только текст и цифры
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "58#̏̏f57#͓͓58#̏̏f", // Никаких символов - только текст и цифры
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "Sha ra ma n ga", // Никаких пробелов
			expected: false,
		},
	}

	// Action
	for number, testCase := range testTable {
		t.Logf("testCase number: %d", number)

		actual := IsValidUsername(testCase.input)
		// Assert
		requires.Equal(testCase.expected, actual)
	}
}

func TestIsValidEmail(t *testing.T) {
	// Arrange
	requires := require.New(t)

	testTable := []struct {
		desc     string // Описание теста
		input    string // Входные данные
		expected bool   // Ожидаемый результат выполнения теста
	}{
		{
			desc:     "Success",
			input:    "user@test.ru",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "usEr@test.com", // Допустимы большие буквы
			expected: true,
		},
		{
			desc:     "Success",
			input:    "useR4442@tes1t.com", // Допустимы большие буквы + цифры
			expected: true,
		},
		{
			desc:     "Success",
			input:    "iam.user@test.yahoo", // Допустимы точки
			expected: true,
		},
		{
			desc:     "Fail",
			input:    "user@testru", // После @ должен быть знак точки
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "usertest.ru", // Должен быть знак @
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "",
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "useruseruser",
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "@user.",
			expected: false,
		},
		{
			desc:     "Success",
			input:    "iam,user@test.yahoo", // Никаких других знаков препинания, кроме точек
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "user@mail.ru.", // Никаких лишник знаков после @
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "user@mail.ru.pitaemsy.slomat.ji.est", // Никаких лишник знаков после @
			expected: false,
		},
		{
			desc:     "Success",
			input:    "iam!user@test.yahoo", // Никаких других знаков препинания, кроме точек
			expected: false,
		},
		{
			desc:     "Success",
			input:    "iam..user@test.yahoo", // Не должно быть подряд идущих точек
			expected: false,
		},
		{
			desc:     "Fail",
			input:    ".user@.",
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "...@@@...",
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้@test.ru", // Никаких левых символов, только английский алфавит
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "пользователь@test.ru", // Никаких левых символов, только английский алфавит
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "us er 06@mail.ru", // Никаких пробелов
			expected: false,
		},
	}

	// Action
	for number, testCase := range testTable {
		t.Logf("testCase number: %d", number)

		actual := IsValidEmail(testCase.input)
		// Assert
		requires.Equal(testCase.expected, actual)
	}
}

func TestIsValidPass(t *testing.T) {
	// Arrange
	requires := require.New(t)

	testTable := []struct {
		desc     string // Описание теста
		input    string // Входные данные
		expected bool   // Ожидаемый результат выполнения теста
	}{
		{
			desc:     "Success",
			input:    "USERPASS13",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "us@.041fWng6",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "11111111111",
			expected: true,
		},
		{
			desc:     "Success",
			input:    "ЭтоПарольПользователя", // Допускаем только русские/английские символы + цифры + знаки
			expected: true,
		},
		{
			desc:     "Success",
			input:    "ThisIsПароль333.04?", // Допускаем только русские/английские символы + цифры + знаки
			expected: true,
		},
		{
			desc:     "Success",
			input:    "...@@@...", // Допускаем такое
			expected: true,
		},
		{
			desc:     "Fail",
			input:    "i am password for user", // Запрещаем любые пробелы
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "i		twotab", // Запрещаем любые отступы (табуляции)
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "e", // Длина не меньше 6 символов
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "          ", // Никаких пробелов без символов
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "",
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "@user.sogTTT,sogTTTsogTTT.sogTTT.sogTTTsogTTT.sogTTT", // Длина не больше 24-ех символов
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้้ด้@", // Никаких левых символов, только русский/английский алфавит + цифры + знаки
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "Tͦ̎h͊̚ȅ̐    Nͧͥ e͂̎  zͩ͗ p̂̔e", // Никаких левых символов, только русский/английский алфавит + цифры + знаки
			expected: false,
		},
		{
			desc:     "Fail",
			input:    "58#̏̏f57#͓͓58#̏̏f", // Никаких левых символов, только русский/английский алфавит + цифры + знаки
			expected: false,
		},
	}

	// Action
	for number, testCase := range testTable {
		t.Logf("testCase number: %d", number)

		actual := IsValidPassword(testCase.input)
		// Assert
		requires.Equal(testCase.expected, actual)
	}
}

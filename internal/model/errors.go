package model

import "errors"

/* model/task.go */
var (
	ErrTaskTitleTooShort      = errors.New("Минимальная длина названия - 3 символа")
	ErrTaskTitleTooLong       = errors.New("Максимальная длина названия - 15 символов")
	ErrTaskDescriptionTooLong = errors.New("Максимальная длина описания - 200 символов")
	ErrTaskInvalidID          = errors.New("Невалидный ID задачи")
)

/* model/user.go */
var (
	ErrUserNameTooShort     = errors.New("Минимальная длина имени - 3 символа")
	ErrUserNameTooLong      = errors.New("Максимальная длина имени - 10 символов")
	ErrUserPasswordTooShort = errors.New("Минимальная длина пароля - 8 символов")
	ErrUserPasswordTooLong  = errors.New("Максимальная длина пароля - 16 символов")
	ErrUserInvalidID        = errors.New("Невалидный ID пользователя")
)

/* repository */
var (
	ErrTaskNotExist     = errors.New("Задача не существует")
	ErrUserNotExist     = errors.New("Пользователь не существует")
	ErrUserAlreadyExist = errors.New("Пользователь с таким name уже существует")
)

/* service */
var (
	ErrTaskIsNil     = errors.New("Задача не может быть ниловой")
	ErrUserInvalidPW = errors.New("Неверный пароль")
)

/* handler */
var (
	ErrUnauthorized = errors.New("Ошибка авторизации")
)

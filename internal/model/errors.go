package model

import "errors"

/* model/task.go */
var (
	ErrTaskTitleTooShort      = errors.New("Минимальная длина названия - 3 символа")
	ErrTaskTitleTooLong       = errors.New("Максимальная длина названия - 15 символов")
	ErrTaskDescriptionTooLong = errors.New("Максимальная длина описания - 200 символов")
	ErrTaskInvalidID          = errors.New("Невалидный ID задачи")
)

/* repository */
var (
	ErrTaskNotExist = errors.New("задача не существует")
)

/* service */
var (
	ErrTaskIsNil = errors.New("задача не может быть ниловой")
)

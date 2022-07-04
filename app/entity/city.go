package entity

import "mysql/app/apperr"

type City struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Population int    `json:"population"`
}

type Cities []*City

func (c City) Validate() error {
	if c.Name == "" {
		return apperr.Errorf(apperr.EINVALID, "Nome invalido")
	}

	if c.Population < 0 {
		return apperr.Errorf(apperr.EINVALID, "Popolazione invalida")
	}

	return nil
}
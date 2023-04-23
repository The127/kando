package tests

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func mergeMaps(maps ...map[string]any) map[string]any {
	result := map[string]any{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func ManufacturerValues(values map[string]any) map[string]any {
	return mergeMaps(map[string]any{
		"name": faker.Word() + " Inc.",
	}, values)
}

func LocationValues(values map[string]any) map[string]any {
	return mergeMaps(map[string]any{
		"name": faker.Word() + " Town",
	}, values)
}

func AssetTypeValues(values map[string]any) map[string]any {
	return mergeMaps(map[string]any{
		"name": faker.Word(),
	}, values)
}

func AssetValues(assetTypeId *uuid.UUID, manufacturerId *uuid.UUID, parentId *uuid.UUID, values map[string]any) map[string]any {
	serialNumber, err := faker.RandomInt(1000, 9999)
	if err != nil {
		panic(err)
	}
	serialNumberString := fmt.Sprintf("S%d", serialNumber)

	batchNumber, err := faker.RandomInt(1000, 9999)
	if err != nil {
		panic(err)
	}
	batchNumberString := fmt.Sprintf("B%d", batchNumber)

	return mergeMaps(map[string]any{
		"name":          faker.Word(),
		"serial_number": &serialNumberString,
		"batch_number":  &batchNumberString,
	}, values)
}

func UserValues(values map[string]any) map[string]any {
	password, ok := values["password"]
	if !ok {
		password = faker.Password()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.(string)), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return mergeMaps(map[string]any{
		"email":           faker.Email(),
		"password":        faker.Password(),
		"hashed_password": hashedPassword,
	}, values)
}

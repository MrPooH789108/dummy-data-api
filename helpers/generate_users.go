package helpers

import (
	"github.com/jaswdr/faker"
	"github.com/omeiirr/dummy-data-api/models"
)

func GenerateUsers(n int) []models.User {
	fake := faker.New()

	var users []models.User
	for i := 0; i < n; i++ {
		user := models.User{
			ID:       fake.Int(),
			Name:     fake.Person().FirstName(),
			Username: fake.Internet().User(),
			Email:    fake.Person().Contact().Email,
		}
		users = append(users, user)
	}

	return users
}

package main

import (
	"github.com/gofiber/fiber/v2"
)

// Handler functions
// getMovies godoc
// @Summary Get all movies
// @Description Get details of all movies
// @Tags movies
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} Movie
// @Router /movies [get]

// Get all movies
func getMovies(c *fiber.Ctx) error {
	return c.JSON(movies)
}

// Get id movie
func getMovie(c *fiber.Ctx) error {
	movieId := c.Params("id")
	for _, movie := range movies {
		if movie.ID == movieId {
			return c.JSON(movie)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

// Create a new movie
func createMovie(c *fiber.Ctx) error {
	movie := new(Movie)
	if err := c.BodyParser(movie); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	movies = append(movies, *movie)
	return c.JSON(movie)
}

// Update a movie
func updateMovie(c *fiber.Ctx) error {
	movieId := c.Params("id")
	movieUpdate := new(Movie)
	if err := c.BodyParser(movieUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, movie := range movies {
		if movie.ID == movieId {
			movies[i] = *movieUpdate
			return c.JSON(movies)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

func deleteMovie(c *fiber.Ctx) error {
	movieId := c.Params("id")
	for i, movie := range movies {
		if movie.ID == movieId {
			/*
				... -> กระจาย slice แต่ละตัวออกจากกันและต่อใหม่
				[1,2,3,4,5] delete 3
				[1,2] + [4,5] = [1,2,4,5]
			*/
			movies = append(movies[:i], movies[i+1:]...)
			return c.JSON(movies)
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}

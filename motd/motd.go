package motd

import "math/rand/v2"

var Messages = []string{
	"use the -a flag to invoke from anywhere",
	"combine -a with -y to invoke from anywhere without confirming",
	"`--env dev` will include .env, .env.dev, .env-dev, and dev.env",
	"pik will start in a viewport if the terminal is too thin",
	"instead of .pik you can use .tasks, or .bin, or any of these with _",
}

func One() string {
	idx := rand.Int() % len(Messages)
	return Messages[idx]
}

package motd

import "math/rand/v2"

var Messages = []string{
	"use the -a flag to invoke from anywhere",
	"combine -a with -y to invoke from anywhere without confirming",
	"`--env dev` will include .env, .env.dev, .env-dev, and dev.env",
	"pik will start in a viewport if the terminal is too thin",
	"instead of .pik you can use .tasks, or .bin, or any of these with _",
	"include .pre. in a pik target filename to make it run as a prerequisite for other targets",
	"a pik target filename with .post. will make it run after another targets completes successfully",
	"you can always run a target after others by adding .final. to its name",
	"you can 'invoke' a directory by having a target with the same name in it",
	"if you're cool, use hjkl to browse me",
	"left and right (h and l) can be used to scroll between sources",
	"add .here. to a target filename to make it run in your shell working directory",
	"the target's first line will be used as description, if it is a comment",
	"create the .pik/.alias file to add additional names to the project for pik",
	"put some utf-8 in .pik/.icon to add an icon to your pik-enabled project",
	"create a target with the same base name but containing .override. to prefer it over the other during invocations",
	"pik crawls both the regular and symlink evaluated locations, if they are different",
	"unsure about what you're doing? use --dry to check what you're running",
}

func One() string {
	idx := rand.Int() % len(Messages)
	return Messages[idx]
}

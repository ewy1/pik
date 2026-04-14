<h1 align="center">
pik
</h1>

> file-based task runner

## the goal

are you ever uncertain what to do after cloning a repository? `pik` aims to fix this by making tasks findable,
predictable and programmable.

running `pik` in a supported repository will index its own file-based task system, as well as 

## getting started

1. create a `.pik` folder in your project
2. put a script in there, for example: `.pik/build.sh` containing `go build .`
3. you can now access this script from almost anywhere by calling `pik build`. If you want to trigger a specific
   projects build, specify `pik project build`, where `project` is the folder name.

## current features

* run targets by their approximate name: `uwu build` will trigger `.uwu/build.sh`, or `.uwu/build.py`, or `make build`
  depending on what's possible.
    * including external targets from `just`
* specify source names to search for the target in that source explicitly: `uwu myproject build` can
  run `../../myproject/.uwu/build.sh`
* unless overridden, targets will run in the directory where the `.uwu` folder resides.
* use `--all` flag to start out-of-tree targets without having to navigate to the directory.
* `--here` to run the target in the current working directory instead of the source directory.
    * `--at` to run the target in an arbitrary location
* target tags in filenames which trigger flag behaviours
* aliases to sources through the `.alias` file
* tui for viewing and running targets
* y/n confirmation with yes as default
    * will be used if we have an uncertain target guess
    * `--yes` to automatically confirm y/n prompts
* autoload .env files
    * both the project root and `.pik` folder will be searched
    * values specified with `--env` will be tried as pre- and suffixes: `--env asdf` will load `.env-asdf`
      and `asdf.env` if they exist.
    * env files are reread for every trigger, meaning you can have a pre-trigger fetch credentials and save it in .env
* create any kind of target: high-level support for shell and python, and arbitrary shells with the shebang.

## planned features

As this program has already gone through a number of iterations and forms, this (hopefully permanent) version will need
time to catch up with all the features it used to have. This list is not exhaustive, but it is ordered by importance I
attach to these features.

* runner for executable files
    * this will also enable arbitrary shells like node by way of the shebang
* indexers for other target types such as `make` and `npm`
* expand tui:
    * adding descriptions to targets based on the first comment in a target
    * support for categories and ordering of targets through the `.order` file
    * search
    * more hotkeys (filter jumping, toggle all, etc.)
* git pre-commit and pre-push triggers
* linking sources together by `.include` and `.wants` files

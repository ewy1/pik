# uwu
> targets should be files, after all

## current features

* run targets by their approximate name: `uwu build` will trigger `.uwu/build.sh`, or `.uwu/build.py`, or `make build` depending on what's possible.
* specify source names to search for the target in that source explicitly: `uwu myproject build` can run `../../myproject/.uwu/build.sh`
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

## planned features

As this program has already gone through a number of iterations and forms, this (hopefully permanent) version will need
time to catch up with all the features it used to have. This list is not exhaustive, but it is ordered by importance I
attach to these features.

* support for .env files
    * the .env files will be reindexed for every script, meaning a `.pre` trigger can prepare the `.env` file for the
      real targets.
* python runner and indexer
    * runner with support for uv and venvs
    * indexer with support for `pyproject.toml` script targets
* runner for executable files
  * this will also enable arbitrary shells like node by way of the shebang
* indexers for other target types such as `make`, `just`, and `npm`
* expand tui:
  * adding descriptions to targets based on the first comment in a target
  * support for categories and ordering of targets through the `.order` file
  * search
  * more hotkeys (filter jumping, toggle all, etc.)
* git pre-commit and pre-push triggers
* linking sources together by `.include` and `.wants` files

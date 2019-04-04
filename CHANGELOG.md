# [1.2.0](https://github.com/auto-staging/builder/compare/1.1.2...1.2.0) (2019-04-04)


### Features

* added new random default variable for AWS ressoource names, this allows environments for longer branch names - fixes [#9](https://github.com/auto-staging/builder/issues/9) ([24ba13b](https://github.com/auto-staging/builder/commit/24ba13b))

## [1.1.2](https://github.com/auto-staging/builder/compare/1.1.1...1.1.2) (2019-04-01)


### Bug Fixes

* replaced test docker image with official auto staging image ([72a3c69](https://github.com/auto-staging/builder/commit/72a3c69))

## [1.1.1](https://github.com/auto-staging/builder/compare/1.1.0...1.1.1) (2019-03-28)


### Bug Fixes

* compile binary for linux ([24a80d9](https://github.com/auto-staging/builder/commit/24a80d9))

# [1.1.0](https://github.com/auto-staging/builder/compare/1.0.2...1.1.0) (2019-02-22)


### Features

* made region dynamic based on lambda region from env var ([52903ed](https://github.com/auto-staging/builder/commit/52903ed))

## [1.0.2](https://github.com/auto-staging/builder/compare/1.0.1...1.0.2) (2019-02-15)


### Bug Fixes

* replaced hardcoded scheduler lambda arn with environment variable - fixes [#2](https://github.com/auto-staging/builder/issues/2) ([3136d8a](https://github.com/auto-staging/builder/commit/3136d8a))

## [1.0.1](https://github.com/auto-staging/builder/compare/1.0.0...1.0.1) (2019-02-15)


### Bug Fixes

* replaced fmt log outputs with lightning log ([9db8256](https://github.com/auto-staging/builder/commit/9db8256))
* replaced go log with lightning log ([b7635b3](https://github.com/auto-staging/builder/commit/b7635b3))

# 1.0.0 (2019-01-26)


### Bug Fixes

* added check for current status before executing codebuild tasks - fixes [#1](https://github.com/auto-staging/builder/issues/1) ([30b97d6](https://github.com/auto-staging/builder/commit/30b97d6))
* added missing destroying failed and updating failed states - fixes [#1](https://github.com/auto-staging/builder/issues/1) ([72dcaa9](https://github.com/auto-staging/builder/commit/72dcaa9))
* added missing error handling for BatchGetProjects ([76627f1](https://github.com/auto-staging/builder/commit/76627f1))
* added missing error handling for get of codebuild project infos ([337a8fa](https://github.com/auto-staging/builder/commit/337a8fa))
* fixed debug output of buildspec ([d1643f2](https://github.com/auto-staging/builder/commit/d1643f2))
* fixed wrong module in logger messages ([ef213e6](https://github.com/auto-staging/builder/commit/ef213e6))
* only allow update and delete controller calls in specific states - fixes [#1](https://github.com/auto-staging/builder/issues/1) ([e9b7978](https://github.com/auto-staging/builder/commit/e9b7978))
* only allow updates from "running" and "updating failed" status - fixes [#2](https://github.com/auto-staging/builder/issues/2) ([2ef9b8b](https://github.com/auto-staging/builder/commit/2ef9b8b))
* set updating failed on update failure instead of initiating failed ([f1e5a43](https://github.com/auto-staging/builder/commit/f1e5a43))
* use RESULT_UPDATE instead of RESULT_CREATE after update - fixes [#3](https://github.com/auto-staging/builder/issues/3) ([f836276](https://github.com/auto-staging/builder/commit/f836276))


### Features

* added build trigger in creation process ([e16e41c](https://github.com/auto-staging/builder/commit/e16e41c))
* added codebuild role arn to event struct and replaced env var map with struct ([f2ec160](https://github.com/auto-staging/builder/commit/f2ec160))
* added create and delete options ([04dd18c](https://github.com/auto-staging/builder/commit/04dd18c))
* added endpoint to delete cloudwatch events for schedules ([050cff7](https://github.com/auto-staging/builder/commit/050cff7))
* added endpoint to update codebuild job with new url and env vars ([377df81](https://github.com/auto-staging/builder/commit/377df81))
* added endpoint to update status after codebuild run ([397ed99](https://github.com/auto-staging/builder/commit/397ed99))
* added finally step to catch status at the end and changed success to int (to use codebuild env var) ([582f351](https://github.com/auto-staging/builder/commit/582f351))
* added function to update / create cloudwatch events for shutdown and startup schedules ([1c9f26f](https://github.com/auto-staging/builder/commit/1c9f26f))
* added functions for "after create success" result ([25e3547](https://github.com/auto-staging/builder/commit/25e3547))
* changed build image to auto-staging docker image and adapted buildspec commands to notify build lambda ([da1b40c](https://github.com/auto-staging/builder/commit/da1b40c))
* implemented delete environment workflow (codebuild adapting, codebuild remove and status update) ([ac8f743](https://github.com/auto-staging/builder/commit/ac8f743))
* made repository url dynamic ([7cea9d2](https://github.com/auto-staging/builder/commit/7cea9d2))
* project init ([30b724c](https://github.com/auto-staging/builder/commit/30b724c))
* remove environment after codebuild destroy ([c74085a](https://github.com/auto-staging/builder/commit/c74085a))
* replaced default logger with lightning log ([f6de4d8](https://github.com/auto-staging/builder/commit/f6de4d8))
* set correct environment status at each error handling step ([6d6badd](https://github.com/auto-staging/builder/commit/6d6badd))
* set default variables for every build job (branch, repository and branch_raw) ([6eda70b](https://github.com/auto-staging/builder/commit/6eda70b))
* set init status at creation start ([cf9ac43](https://github.com/auto-staging/builder/commit/cf9ac43))

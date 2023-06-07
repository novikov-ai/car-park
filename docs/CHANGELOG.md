# [Changelog](https://keepachangelog.com/en/1.0.0/)

## [0.6.0] - 2023-06-07

### Added

- new Manager model
- add basic Auth

## [0.5.0] - 2023-06-03

### Added

- new models: Enterprise & Driver
- new providers for created models
- add migration files for new tables

### Changed

- db init-file expand with new data
- refactor main & vehicle provider
- rename db-migrations 

## [0.4.0] - 2023-05-31

### Added

- add json tags for Vehicle model

### Changed

- errors handling
- change api-path and clean-up

### Fixed

- fix storage Query() with args

## [0.3.0] - 2023-05-26

### Added

- endpoint '/api/v1/vehicles/admin'
- vehicles controller with CRUD: Create, Update, Delete

### Changed

- refactor storage with interface and MVC logic

## [0.2.0] - 2023-05-22

### Added

- models:
  - model
  - types enum
- add models fetcher & view


## [0.1.0] - 2023-05-21

### Added

- models:
  - vehicle
  - colors
- vehicle fetcher & web-view
- storage & migrations
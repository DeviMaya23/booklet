## 1. Filter Structs & Repository Interfaces

- [x] 1.1 Add `ListArtistFilters` struct (with `query` tag on `Q *string`) to `usecase/artist_repository.go`; update `ArtistRepository.List` signature to accept it
- [x] 1.2 Add `ListCharacterFilters` struct (with `query` tag on `Q *string`) to `usecase/character_repository.go`; update `CharacterRepository.List` signature to accept it
- [x] 1.3 Add `ListImageFilters` struct (with `query` tags on `Q *string`, `CharacterIDs []uuid.UUID`, `ArtistIDs []uuid.UUID`) to `usecase/image_repository.go`; update `ImageRepository.List` signature to accept it

## 2. Repository Implementations

- [x] 2.1 Update `artistRepository.List` to conditionally append `LOWER(name) LIKE ?` when `Q` is set
- [x] 2.2 Update `characterRepository.List` to conditionally append `LOWER(name) LIKE ?` when `Q` is set
- [x] 2.3 Update `imageRepository.List` to conditionally append `LOWER(title) LIKE ?` when `Q` is set
- [x] 2.4 Update `imageRepository.List` to conditionally join `image_characters` and apply `character_id IN ?` when `CharacterIDs` is non-empty; add `.Distinct()` to avoid duplicate rows
- [x] 2.5 Update `imageRepository.List` to conditionally append `artist_id IN ?` when `ArtistIDs` is non-empty

## 3. Usecase Layer

- [x] 3.1 Update `artistUsecase.List` signature to accept `ListArtistFilters`; pass through to repo
- [x] 3.2 Update `characterUsecase.List` signature to accept `ListCharacterFilters`; pass through to repo
- [x] 3.3 Update `imageUsecase.List` signature to accept `ListImageFilters`; pass through to repo

## 4. Handler Layer

- [x] 4.1 Update `ArtistUsecase` interface in `artist_handler.go` to reflect new `List` signature
- [x] 4.2 Update `CharacterUsecase` interface in `character_handler.go` to reflect new `List` signature
- [x] 4.3 Update `ImageUsecase` interface in `image_handler.go` to reflect new `List` signature
- [x] 4.4 Update `ListArtists` handler: bind `ListArtistFilters` via `c.Bind`, pass to usecase
- [x] 4.5 Update `ListCharacters` handler: bind `ListCharacterFilters` via `c.Bind`, pass to usecase
- [x] 4.6 Update `ListImages` handler: bind `ListImageFilters` via `c.Bind`; parse `CharacterIDs` and `ArtistIDs` from `[]string` to `[]uuid.UUID`, return 400 on malformed UUID; pass to usecase

## 5. Tests

- [x] 5.1 Update `spyArtistUsecase.List` signature in `artist_handler_test.go`; add handler test for `q` filter binding
- [x] 5.2 Update `spyCharacterUsecase.List` (or equivalent) in `character_handler_test.go`; add handler test for `q` filter binding
- [x] 5.3 Update image handler test spy `List` signature; add handler tests for `q`, `character_ids`, `artist_ids` binding and UUID validation (malformed UUID → 400)
- [x] 5.4 Update `fakeArtistRepository.List` in `usecase/fakes_test.go` (or inline); add usecase test for `List` with `Q` filter
- [x] 5.5 Update `fakeCharacterRepository.List`; add usecase test for `List` with `Q` filter
- [x] 5.6 Update `fakeImageRepository.List`; add usecase test for `List` with `Q`, `CharacterIDs`, `ArtistIDs` filters
- [x] 5.7 Update artist repository integration test for `List` to cover `Q` filter
- [x] 5.8 Update character repository integration test for `List` to cover `Q` filter
- [x] 5.9 Update image repository integration test for `List` to cover `Q`, `CharacterIDs`, and `ArtistIDs` filters; verify no duplicate rows on multi-character match

## 6. Bruno Files

- [x] 6.1 Update the `GET /artists` Bruno file to include an example request with `q` query param
- [x] 6.2 Update the `GET /characters` Bruno file to include an example request with `q` query param
- [x] 6.3 Update the `GET /images` Bruno file to include example requests with `q`, `character_ids`, and `artist_ids` query params

## 7. Lint

- [x] 7.1 Run `golangci-lint run ./...` and fix any issues

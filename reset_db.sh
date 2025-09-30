# Used to reset db for testing purposes

cd sql/schema
goose postgres postgres://postgres:postgres@localhost:5432/gator?sslmode=disable down
goose postgres postgres://postgres:postgres@localhost:5432/gator?sslmode=disable up
cd ../..
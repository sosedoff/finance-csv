build:
	gox -osarch="linux/amd64 darwin/amd64" -output="./bin/finance-csv_{{.OS}}_{{.Arch}}"
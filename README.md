# mysql2go

convert mysql table into go struct.

# Example

```bash
$ mysql2go --help

Hi~, Welcome to mysql2go.
        
mysql2go is used to convert mysql table structure to go struct.

Usage:
  mysql2go [flags]

Flags:
  -h, --help             help for mysql2go
  -i, --infile string    input file path，eg: your sql file.
  -o, --outfile string   output file path. go file.


$ mysql2go table.sql

or

$ mysql2go table.sql table.go

or

$ mysql2go --infile=table.sql --outfile=table.go

```
package migration

import _ "embed"

//go:embed schema.sql
var DDL string

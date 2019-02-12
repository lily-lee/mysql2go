package convert

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	start = `type {TYPE} struct {
`
	comment = `	// {COMMENT}
`
	tab = `	{FIELD} {FIELD_TYPE} ` + "`json:" + `"{JSON}"` + "`" + `
`
	end = `
}

`
	unknown = "unknown"
)

var s = `package {PACKAGE} 

`

func Convert(infile, outfile string) error {
	// read file
	f, err := os.Open(infile)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if outfile == "" {
		outfile = infile + ".go"
	}

	return convert(f, outfile)
}

func convert(f *os.File, outfile string) error {
	packageName := filepath.Dir(outfile)
	if packageName == "." {
		packageName = "main"
	}

	s = strings.Replace(s, "{PACKAGE}", packageName, 1)

	otf, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	defer otf.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	otf.Write([]byte(s))

	input := bufio.NewScanner(f)
	sql := false

	for input.Scan() {
		line := strings.Trim(input.Text(), " ")
		field0 := regexp.MustCompile("`(.*)`").FindString(line)

		if strings.Contains(strings.ToLower(line), "create table") {
			sql = true
			tableName := getName(field0)
			//s += strings.Replace(start, "{TYPE}", tableName, 1)
			otf.Write([]byte(strings.Replace(start, "{TYPE}", tableName, 1)))
			continue
		}

		if !sql {
			continue
		}

		fields := strings.Fields(strings.Replace(line, field0, strings.Replace(field0, " ", "_", -1), -1))
		if len(fields) <= 1 || strings.ToLower(fields[0]) == "key" || strings.ToLower(fields[0]) == "primary" {
			continue
		}

		if fields[0] == ")" {
			//s += end
			otf.Write([]byte(end))
			sql = false
			continue
		}

		t := ""
		name := getName(fields[0])
		for i := range fields {
			c := comment
			if strings.ToLower(fields[i]) == "comment" {
				//s += strings.Replace(c, "{COMMENT}", name+" "+strings.Trim(fields[i+1], "',"), 1)
				otf.Write([]byte(strings.Replace(c, "{COMMENT}", name+" "+strings.TrimLeft(strings.TrimRight(strings.Join(fields[i+1:], " "), "',"), "'"), 1)))
			}
		}

		t = strings.Replace(tab, "{FIELD}", name, 1)
		t = strings.Replace(t, "{JSON}", strings.Trim(fields[0], "`"), 1)
		t = strings.Replace(t, "{FIELD_TYPE}", getType(fields), 1)

		otf.Write([]byte(t))
		//s += t
	}

	return nil
}

var upperWords = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func getName(field string) string {
	s := strings.Split(strings.Replace(strings.Trim(field, "`"), "_", " ", -1), " ")
	for i := range s {
		if _, ok := upperWords[strings.ToUpper(s[i])]; ok {
			s[i] = strings.ToUpper(s[i])
		} else {
			s[i] = strings.Title(s[i])
		}
	}

	return strings.Join(s, "")
}

func getType(fields []string) string {
	if len(fields) < 2 {
		return unknown
	}

	field1 := strings.ToLower(fields[1])
	unsign := ""
	if len(fields) >= 3 {
		if strings.Contains(strings.ToLower(fields[2]), "unsign") {
			unsign = "u"
		}
	}

	name := regexp.MustCompile(`[a-z]{1,}`).FindString(field1)
	switch name {
	case "tinyint":
		return unsign + "int8"

	case "smallint":
		return unsign + "int16"

	case "int", "mediumint":
		return unsign + "int32"

	case "bigint":
		return unsign + "int64"

	case "float", "decimal":
		return "float32"

	case "double":
		return "float64"

	case "date", "char", "varchar", "blob", "text", "tinyblob", "tinytext", "mediumblob", "mediumtext", "longblob", "longtext", "enum":
		return "string"

	case "timestamp", "datetime":
		return "time.Time"

	}

	return unknown
}

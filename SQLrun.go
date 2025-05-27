package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/base64"
	"database/sql"
	"strings"
	"unsafe"
	_ "github.com/go-sql-driver/mysql"
)

//export SQLrun
func SQLrun(conexion *C.char, query *C.char, args **C.char, argCount C.int) *C.char {
	goConexion := C.GoString(conexion)
	goQuery := C.GoString(query)
	
	var goArgs []interface{}
	if argCount > 0 {
		argSlice := (*[1 << 30]*C.char)(unsafe.Pointer(args))[:argCount:argCount]
		for _, arg := range argSlice {
			goArgs = append(goArgs, C.GoString(arg))
		}
	}

	_, result := sqlRunInternal(goConexion, goQuery, goArgs...)
	
	return C.CString(result)
}

func sqlRunInternal(conexion, query string, args ...any) (int, string) {
	respuesta := "["
	db, err := sql.Open("mysql", conexion)
	if err != nil {
		db.Close()
		return 1, "No se logro aperturar conexion a: " + conexion + "\n"
	}
	defer db.Close()
	rows, err := db.Query(query, args...)
	if err != nil {
		db.Close()
		rows.Close()
		return 1, "No se logro ejecutar la query: " + query + "\n"
	}
	defer rows.Close()

	llaves := 0

	for {
		colTypes, err := rows.ColumnTypes()
		if err != nil {
			db.Close()
			rows.Close()
			return 1, "No se logro obtener el tipado al ejecutar la query: " + query + "\n"
		}		
		columns, err := rows.Columns()
		if err != nil {
			db.Close()
			rows.Close()
			return 1, "No se logro obtener respuesta al ejecutar la query: " + query + "\n"
		}
		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				db.Close()
				rows.Close()
				return 1, "No se logro obtener datos al ejecutar la query: " + query + "\n"
			}
			if llaves == 0 {
				respuesta = respuesta + "\n\t{"
				llaves = llaves + 1
			} else {
				respuesta = respuesta + ",\n\t{"
			}
			campos := 0
			var value string
			for i, col := range values {
				if col == nil {
					value = ""
				} else {
					if(strings.Contains(colTypes[i].DatabaseTypeName(), "BLOB")){
						value = base64.StdEncoding.EncodeToString(col)
					} else {
						value = string(col)
						value = strings.ReplaceAll(value, "\"", "")
						value = strings.ReplaceAll(value, "\n", "\\n")
					}
				}
				if campos == 0 {
					respuesta = respuesta + "\n\t\t\"" + strings.ReplaceAll(strings.ReplaceAll(columns[i], "\"", ""), "\n", "\\n") + "\": " + "\"" + value + "\""
					campos = campos + 1
				} else {
					respuesta = respuesta + ", \n\t\t\"" + strings.ReplaceAll(strings.ReplaceAll(columns[i], "\"", ""), "\n", "\\n") + "\": " + "\"" + value + "\""
				}
			}
			respuesta = respuesta + "\n\t}"
		}

		if !rows.NextResultSet() {
			break
		}
	}
	db.Close()
	rows.Close()
	respuesta = respuesta + "\n]"

	if respuesta == "[\n]" && strings.HasPrefix(query, "CALL ") {
		return 0, "{\n\t\"EJECUCION\": \"OK\"\n}"
	} else {
		return 0, respuesta
	}
}

//export FreeString
func FreeString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

func main() {}
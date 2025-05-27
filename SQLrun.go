package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strconv"
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
			argStr := C.GoString(arg)

			switch {
			case strings.HasPrefix(argStr, "int::"):
				intVal, err := strconv.ParseInt(argStr[5:], 10, 64)
				if err != nil {
					return C.CString(fmt.Sprintf(`{"error":"Error parseando entero: %s"}`, argStr[5:]))
				}
				goArgs = append(goArgs, intVal)

			case strings.HasPrefix(argStr, "float::"), strings.HasPrefix(argStr, "double::"):
				prefixLen := 7
				if strings.HasPrefix(argStr, "double::") {
					prefixLen = 8
				}
				floatVal, err := strconv.ParseFloat(argStr[prefixLen:], 64)
				if err != nil {
					return C.CString(fmt.Sprintf(`{"error":"Error parseando float: %s"}`, argStr[prefixLen:]))
				}
				goArgs = append(goArgs, floatVal)

			case strings.HasPrefix(argStr, "bool::"):
				boolVal, err := strconv.ParseBool(argStr[6:])
				if err != nil {
					return C.CString(fmt.Sprintf(`{"error":"Error parseando booleano: %s"}`, argStr[6:]))
				}
				goArgs = append(goArgs, boolVal)

			case strings.HasPrefix(argStr, "null::"):
				goArgs = append(goArgs, nil)

			case strings.HasPrefix(argStr, "blob::"):
				data, err := base64.StdEncoding.DecodeString(argStr[6:])
				if err != nil {
					return C.CString(fmt.Sprintf(`{"error":"Error decodificando blob: %v"}`, err))
				}
				goArgs = append(goArgs, data)

			default:
				goArgs = append(goArgs, argStr)
			}
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

	if respuesta == "[\n]" && (strings.HasPrefix(query, "CALL ") || strings.HasPrefix(query, "INSERT ") || strings.HasPrefix(query, "UPDATE ") || strings.HasPrefix(query, "DELETE ") || strings.HasPrefix(query, "DROP ") ) {
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

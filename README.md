# SQL2JSON
Libreria C para consultas a Bases de Datos MariaDB/MySQL con Resultado JSON, Libreria basada en el proyecto: https://gitlab.com/RicardoValladares/api-mysql.git se recompilo usando: go build -o SQLrun.dll -buildmode=c-shared SQLrun.go


### Descargar Libreria:
| Linux | Windows |
| --- | --- |
| `wget https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.so` | `Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.dll -OutFile ./SQLrun.dll` |
| `wget https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.h` | `Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.h -OutFile ./SQLrun.h` |


### Compilar:
| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./SQLrun.so` | `gcc -o main.exe main.c ./SQLrun.dll` |


Ejemplo:
```C
#include <stdio.h>
#include <stdlib.h>
#include "SQLrun.h"

int main() {
    // Configuración de conexión
    char* conexion = "root:123456@tcp(127.0.0.1:3306)/test";
    
    // Consulta SQL con parámetros
    char* query = "select now();";
        
    // Llamar a la función
    char* result = SQLrun(conexion, query, 0, 0);
    printf("Resultado: %s\n", result);
    
    // Liberar memoria
    FreeString(result);
    
    return 0;
}
```


Ejemplo:
```C
#include <stdio.h>
#include <stdlib.h>
#include "SQLrun.h"

int main() {
    char* conn = "root:123456@tcp(127.0.0.1:3306)/chat";
    char* query = "INSERT INTO chat.usuario(nickname, picture) VALUES (?, ?);";
    
    // Preparar argumentos
    char* args[2];
    
    // String (normal)
    args[0] = "Ricardo";
    
    // Blob (imagen en base64) //int::  float::  double::  bool::  null::  blob::
    args[1] = "blob::iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAArSURBVBhXY/iPA0AlGBgwGFAKlwQmAKrAIgcVRZODCsI5cAAVgVDo4P9/AHe4m2U/OJCWAAAAAElFTkSuQmCC";
    
    // Convertir a arreglo de char*
    char** args_ptr = (char**)malloc(2 * sizeof(char*));
    for (int i = 0; i < 2; i++) {
        args_ptr[i] = strdup(args[i]);
    }
    
    // Ejecutar consulta
    char* result = SQLrun(conn, query, args_ptr, 2);
    printf("Resultado: %s\n", result);
    
    // Liberar memoria
    FreeString(result);
    for (int i = 0; i < 2; i++) {
        free(args_ptr[i]);
    }
    free(args_ptr);
    
    return 0;
}
```

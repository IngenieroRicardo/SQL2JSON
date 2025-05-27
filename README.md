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

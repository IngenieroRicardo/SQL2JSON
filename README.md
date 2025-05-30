# SQL2JSON

Librería en C para realizar consultas a bases de datos MariaDB/MySQL y obtener los resultados en formato JSON.  
Esta librería está basada en el proyecto original: https://gitlab.com/RicardoValladares/api-mysql.git  
Fue recompilada usando el siguiente comando: go build -o SQL2JSON.dll -buildmode=c-shared SQL2JSON.go

---

### 📥 Descargar la librería

| Linux | Windows |
| --- | --- |
| `wget https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQL2JSON.so` | `Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQL2JSON.dll -OutFile ./SQL2JSON.dll` |
| `wget https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQL2JSON.h` | `Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQL2JSON.h -OutFile ./SQL2JSON.h` |

---

### 🛠️ Compilar

| Linux | Windows |
| --- | --- |
| `gcc -o main.bin main.c ./SQL2JSON.so` | `gcc -o main.exe main.c ./SQL2JSON.dll` |
| `x86_64-w64-mingw32-gcc -o main.exe main.c ./SQL2JSON.dll` |  |

---

### 🧪 Ejemplo básico

```C
#include <stdio.h>
#include "SQL2JSON.h"

int main() {
    // Ejemplo de conexión y consulta
    char* conexion = "root:123456@tcp(127.0.0.1:3306)/test";
    char* query = "SELECT now();";
    
    SQLResult resultado = SQLrun(conexion, query, NULL, 0);
    
    if (resultado.is_error) {
        printf("Error: %s\n", resultado.json);
    } else if (resultado.is_empty) {
        printf("Consulta ejecutada pero no retornó datos\n");
        printf("JSON: %s\n", resultado.json); // Mostrará {"status":"OK"} o []
    } else {
        printf("Datos obtenidos:\n%s\n", resultado.json);
    }
    
    // Liberar memoria
    FreeSQLResult(&resultado);
    
    return 0;
}
```

---

### 🧪 Ejemplo con parámetros

```C
#include <stdio.h>
#include "SQL2JSON.h"

int main() {
    // Ejemplo de conexión e inserción
    char* conexion = "root:123456@tcp(127.0.0.1:3306)/test";
    
    // Ejemplo 1: Consulta INSERT con parámetros
    char* consulta_insert = "INSERT INTO chat.usuario(nickname, picture) VALUES (?, ?);";
    
    // Preparar los argumentos para el INSERT
    char* argumentos_insert[2];
    argumentos_insert[0] = "Ricardo";  // Parámetro de tipo cadena (nickname)
    // Parámetro de tipo blob (imagen codificada en base64)
    argumentos_insert[1] = "blob::iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAArSURBVBhXY/iPA0AlGBgwGFAKlwQmAKrAIgcVRZODCsI5cAAVgVDo4P9/AHe4m2U/OJCWAAAAAElFTkSuQmCC";
    
    // Convertir a un arreglo de char** (necesario para la función SQLrun)
    char** ptr_argumentos_insert = (char**)malloc(2 * sizeof(char*));
    for (int i = 0; i < 2; i++) {
        ptr_argumentos_insert[i] = strdup(argumentos_insert[i]); // Copiar cada argumento
    }
    
    // Ejecutar la consulta INSERT
    SQLResult resultado_insert = SQLrun(conexion, consulta_insert, ptr_argumentos_insert, 2);
    
    // Mostrar los resultados
    printf("Resultado del INSERT:\n");
    printf("JSON: %s\n", resultado_insert.json);         // Respuesta en formato JSON
    printf("Es error: %d\n", resultado_insert.is_error); // 1 si hubo error, 0 si éxito
    printf("Está vacío: %d\n\n", resultado_insert.is_empty); // 1 para consultas que no retornan datos
    
    // Liberar los recursos utilizados
    FreeSQLResult(&resultado_insert); // Liberar la memoria del resultado
    
    // Liberar los argumentos copiados
    for (int i = 0; i < 2; i++) {
        free(ptr_argumentos_insert[i]);
    }
    free(ptr_argumentos_insert); // Liberar el arreglo de argumentos
    
    return 0;
}
```



📝 Los tipos de datos soportados en los argumentos son:
- `string` (por defecto)
- `int::123`
- `float::3.14`
- `double::2.718`
- `bool::true` / `bool::false`
- `null::`
- `blob::<base64>`

---


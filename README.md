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
#include <stdlib.h>
#include "SQL2JSON.h"

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

---

### 🧪 Ejemplo con parámetros

```C
#include <stdio.h>
#include <stdlib.h>
#include "SQL2JSON.h"

int main() {
    char* conn = "root:123456@tcp(127.0.0.1:3306)/chat";
    char* query = "INSERT INTO chat.usuario(nickname, picture) VALUES (?, ?);";
    
    // Preparar argumentos
    char* args[2];
    args[0] = "Ricardo";  // String simple
    args[1] = "blob::iVBORw0KGgoAAAANSUhEUgAAAAgAAAAICAIAAABLbSncAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAArSURBVBhXY/iPA0AlGBgwGFAKlwQmAKrAIgcVRZODCsI5cAAVgVDo4P9/AHe4m2U/OJCWAAAAAElFTkSuQmCC";  // Imagen en base64

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



📝 Los tipos de datos soportados en los argumentos son:
- `string` (por defecto)
- `int::123`
- `float::3.14`
- `double::2.718`
- `bool::true` / `bool::false`
- `null::`
- `blob::<base64>`

---


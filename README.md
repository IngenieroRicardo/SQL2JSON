# SQL2JSON
Libreria C para consultas a Bases de Datos MariaDB/MySQL con Resultado JSON, Libreria basada en el proyecto: https://gitlab.com/RicardoValladares/api-mysql/

Descargar libreria dentro del proyecto en Linux:
```bash
wget https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.so
wget https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.h
```

Descargar libreria dentro del proyecto en Windows:
```powershell
Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.dll
Invoke-WebRequest https://raw.githubusercontent.com/IngenieroRicardo/SQL2JSON/refs/heads/main/SQLrun.h
```




Compilar usando la libreria
```bash
gcc main.c ./SQLrun.so
```
Recompilar libreria
```bash
cd SQL2JSON
make
```

Ejemplos para main.c
```bash
gcc main.c ./SQL2JSON/SQLrun.so
```

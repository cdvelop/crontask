#!/bin/bash

# Crear directorio bin si no existe
mkdir -p bin

# Descargar dependencias
go mod tidy

# Verificar formato y linting básico
go fmt ./...
go vet ./...

# Ejecutar tests
echo "Ejecutando tests..."
go test ./... || { echo "Los tests han fallado"; exit 1; }

# Crear un archivo temporal para la nueva compilación
TEMP_FILE=$(mktemp)
FINAL_FILE="bin/crontask.exe"

# Configurar entorno para compilación cruzada Windows
export GOOS=windows GOARCH=amd64 CGO_ENABLED=0

# Compilar el binario con optimizaciones al archivo temporal
echo "Compilando versión temporal..."
go build -ldflags="-s -w" -o "$TEMP_FILE" ./cmd/crontask

# Comprobar si el binario existe
if [ -f "$FINAL_FILE" ]; then
    # Comparar los archivos
    if cmp -s "$TEMP_FILE" "$FINAL_FILE"; then
        echo "No hay cambios en el binario. Manteniendo versión existente."
    else
        echo "Se detectaron cambios. Actualizando binario..."
        mv "$TEMP_FILE" "$FINAL_FILE"
        echo "Binario actualizado: $FINAL_FILE"
    fi
else
    echo "Creando nuevo binario..."
    mv "$TEMP_FILE" "$FINAL_FILE"
    echo "Binario creado: $FINAL_FILE"
fi

# Limpiar archivo temporal si todavía existe
if [ -f "$TEMP_FILE" ]; then
    rm "$TEMP_FILE"
fi

echo "Proceso de compilación completado."

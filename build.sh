#!/bin/bash

# Check for uncommitted changes in Go files first
echo "Checking for changes in Go files..."
# Use git status to check for modified or staged .go files. Exclude untracked files.
if [[ -z $(git status --porcelain --untracked-files=no -- *.go) ]]; then
    echo "No changes detected in Go files. Skipping build."
    exit 0
else
    echo "Changes detected in Go files. Proceeding with build..."
fi

# Crear directorio bin si no existe
mkdir -p bin

# Descargar dependencias
go mod tidy

# Verificar formato y linting básico
go fmt ./...
go vet ./...

# Ejecutar tests en subdirectorios específicos para evitar recursión con build_test.go
echo "Ejecutando tests en subdirectorios..."
# Modifica las rutas (ej: ./cmd/..., ./pkg/...) según la estructura de tu proyecto
# donde residen los tests que NO deben ejecutar este script.
# Si tus tests están solo bajo cmd/, usa solo ./cmd/...
# Si no tienes otros tests, puedes comentar o ajustar esta línea.
# Forzamos la ejecución en cmd/ para el ejemplo, asumiendo que ahí hay código/tests.
# Si no hay tests en cmd/, esto pasará silenciosamente.
go test ./cmd/... || { echo "Los tests en ./cmd/... han fallado"; exit 1; }
# Añade aquí otras rutas si tienes tests en más subdirectorios, por ejemplo:
# go test ./internal/... || { echo "Los tests en ./internal/... han fallado"; exit 1; }

# Definir el archivo final
FINAL_FILE="bin/crontask.exe"

# Configurar entorno para compilación cruzada Windows
export GOOS=windows GOARCH=amd64 CGO_ENABLED=0

# Compilar el binario con optimizaciones directamente al archivo final
echo "Compilando..."
go build -ldflags="-s -w" -o "$FINAL_FILE" ./cmd/crontask

# Check if build was successful (go build exits with 0 on success)
if [ $? -eq 0 ]; then
    echo "Binario compilado/actualizado: $FINAL_FILE"
else
    echo "La compilación ha fallado."
    exit 1
fi

echo "Proceso de compilación completado."

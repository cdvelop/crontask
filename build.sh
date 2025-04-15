#!/bin/bash

# Crear directorio bin si no existe
mkdir -p bin

# Limpiar binarios anteriores
rm -f bin/crontask.exe

# Descargar dependencias
go mod tidy

# Verificar formato y linting básico
go fmt ./...
go vet ./...

# Ejecutar tests
echo "Ejecutando tests..."
go test ./... || { echo "Los tests han fallado"; exit 1; }

# Configurar entorno para compilación cruzada Windows
export GOOS=windows GOARCH=amd64 CGO_ENABLED=0

# Compilar el binario con optimizaciones
go build -ldflags="-s -w" -o bin/crontask.exe ./cmd/crontask

echo "Compilación completada: bin/crontask.exe"

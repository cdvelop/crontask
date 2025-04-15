# CronTask

CronTask es un gestor de tareas programadas en Go con soporte para configuración tanto en código como en archivos YAML. Permite programar tareas con la sintaxis crontab de forma simple y potente.

## Instalación

```bash
go get github.com/[tu-usuario]/crontask
```

## Características principales

- Sintaxis crontab familiar y potente ("* * * * *")
- Configuración mediante código Go o archivos YAML
- Soporte para entornos nativos y WASM
- Ejecución de comandos del sistema o funciones Go
- API simple y fácil de usar

## Guía rápida

### 1. Crear un motor de tareas programadas

```go
// Crear un motor con configuración predeterminada (carga crontasks.yml automáticamente)
engine := crontask.NewCronTaskEngine()

// O con configuración personalizada
engine := crontask.NewCronTaskEngine(crontask.Config{
    TasksPath: "ruta/a/mis/tareas.yml",
    NoAutoSchedule: false, // true para deshabilitar la programación automática
})
```

### 2. Programar tareas desde código

```go
// Añadir una tarea programada (una función sin argumentos)
err := engine.AddTask("* * * * *", miFuncion)

// Añadir una tarea con argumentos
err := engine.AddTask("0 12 * * *", miFuncionConArgs, "argumento1", 123)

// Ejecutar un comando del sistema
task := crontask.Task{
    Name:     "Backup diario",
    Schedule: "0 0 * * *",       // A medianoche todos los días
    Command:  "/ruta/al/script.sh",
    Args:     "arg1 arg2",
}
err := engine.ScheduleTask(task)
```

### 3. Configuración mediante YAML

Archivo `crontasks.yml`:
```yaml
- name: "Backup Windows"
  schedule: "0 0 * * *"  # Todos los días a medianoche
  command: "c:\\Program Files\\Backup\\backup.exe"
  args: "%USERPROFILE%\\Backups\\config.ini"

- name: "Actualizar datos"
  schedule: "0 */4 * * *"  # Cada 4 horas
  command: "/usr/bin/update-data.sh"
  args: "--force"
```

## Sintaxis Crontab

La sintaxis crontab sigue el formato estándar de 5 campos:

```
*    *    *    *    *
┬    ┬    ┬    ┬    ┬
│    │    │    │    └─  Día de la semana (0-6) (Domingo=0)
│    │    │    └──────  Mes (1-12)
│    │    └───────────  Día del mes (1-31)
│    └────────────────  Hora (0-23)
└─────────────────────  Minuto (0-59)
```

### Ejemplos de Programación

- `* * * * *` - Cada minuto
- `0 12 * * *` - Todos los días a las 12:00
- `0 0 * * 1-5` - A medianoche de lunes a viernes
- `*/15 * * * *` - Cada 15 minutos
- `0 7 15 1 *` - A las 07:00 del 15 de enero
- `0 7 * * 1,4` - A las 07:00 de lunes y jueves
- `5 14,19 * * *` - Dos veces al día, a las 14:05 y a las 19:05
- `0 */5 * * *` - Cada 5 horas (a las 0:00, 5:00, 10:00, 15:00, 20:00)


Para una referencia visual, puedes usar la herramienta: [crontab.guru](https://crontab.guru/)

## Conversión de días de la semana

- Lunes: 1
- Martes: 2
- Miércoles: 3
- Jueves: 4
- Viernes: 5
- Sábado: 6
- Domingo: 0

## Agradecimientos

Este proyecto está basado en el trabajo de:
- [github.com/mileusna/crontab](https://github.com/mileusna/crontab)



# Wails Shell

El host Wails del Studio vive en `apps/studio/main.go` y `apps/studio/app.go`.

El frontend shell vive en `apps/studio/frontend` y se enlaza mediante `apps/studio/wails.json`.

El directorio se conserva para documentación y futuras capacidades específicas del host. El dominio Go continúa en el módulo raíz y no se duplica aquí.

Para compilar desde `apps/studio`, usa la CLI Wails v2 fijada en el entorno de desarrollo:

```text
wails build
```

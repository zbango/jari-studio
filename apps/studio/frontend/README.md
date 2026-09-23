# Studio Frontend

UI React/TypeScript embebida por el host Wails del Studio.

Estado actual:

- Build Vite configurado y verificable.
- Bindings Wails generados para proyectos, revisiones, adapters y terminales.
- El navegador sólo sirve como preview visual; las capacidades de filesystem,
  SQLite y terminal requieren el host nativo.

Para ejecutar el Studio completo, inicia Wails desde `apps/studio`.

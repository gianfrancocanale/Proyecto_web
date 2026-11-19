# Mapa de sitio y capturas — PredictOne

**Fecha:** 2026-11-19  
**Fuente:** recorrido real con Chrome headless sobre la API levantada con `make run`.

## Diagrama de navegación

```mermaid
flowchart TD
    A["Inicio<br/>/"] --> B["Escuderías<br/>/escuderias"]
    A --> C["Pilotos<br/>/pilotos"]
    A --> D["Grandes Premios<br/>/gran_premios"]
    A --> E["Ingresar<br/>/ingresar"]
    A --> F["Registrar<br/>/registrar"]
    X["Ruta inexistente /*"] --> G["Página 404"]
    G -->|"Volver al inicio"| A
    F -->|"alta OK"| E

    E -->|"login usuario"| A
    E -->|"login admin"| H["Admin<br/>/admin"]
    E -->|"sesión activa"| A
    E -.->|"anónimo → /apuestas o /admin"| E

    I["Apuestas<br/>/apuestas"] -->|"logo"| A
    I --> M["Modificar apuesta<br/>/apuestas/modificar"]
    I -.->|"htmx: POST / DELETE"| I
    S["Salir<br/>/salir"] --> A

    H -->|"menú"| B
    H -->|"menú"| C
    H -->|"menú"| D
    H --> J["Programar carrera<br/>/admin/carrera<br/>Cargar resultado<br/>/admin/resultado"]
    H --> K["CRUD<br/>pilotos / escuderías / GP<br/>(GET, POST, DELETE)"]

    B -.->|"htmx ?orden=asc/desc"| B
    C -.->|"htmx ?orden=asc/desc"| C

    U["Usuario normal → /admin"] -->|"403 No autorizado"| U
    L["API REST (JSON)<br/>/usuarios<br/>/pilotos_historicos<br/>/gran_premios_historicos"]
    T["Estáticos<br/>/front/"]

    classDef publico fill:#e8f1ff,stroke:#2b6cb0,color:#1a365d;
    classDef sesion fill:#e6ffed,stroke:#2f855a,color:#1c4532;
    classDef admin fill:#fff0f0,stroke:#c53030,color:#742a2a;
    classDef error fill:#fffaf0,stroke:#dd6b20,color:#7b341e;
    classDef api fill:#f4f5f7,stroke:#718096,color:#2d3748;

    class A,B,C,D,E,F publico;
    class I,M,S sesion;
    class H,J,K admin;
    class G,X,U error;
    class L,T api;
```

## Leyenda de accesos

| Condición | Destino | Resultado observado |
|---|---|---|
| Visitante anónimo | `/apuestas`, `/admin` | Redirección a `/ingresar` |
| Login de usuario | `/` | Menú con “Apuestas” y “Cerrar sesión” |
| Login de admin | `/admin` | Menú de administración |
| Usuario normal | `/admin` | `403 No autorizado` (texto plano) |
| Sesión activa | `/ingresar` | Redirección a `/` |
| Cualquier ruta inexistente | Página 404 | Código HTTP 404 + link al inicio |


Los metadatos de cada captura (URL, título y altura) están en [`capturas/manifest.json`](capturas/manifest.json).

> Las capturas 10 y 13 evidencian dos hallazgos del [informe de navegación](INFORME_NAVEGACION.md): el 403 sin estilo y la pantalla de edición de Grandes Premios con el layout público en lugar del layout de administración.

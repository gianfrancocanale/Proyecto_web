# Informe de navegación — PredictOne

**Fecha:** 2026-11-19  
**Entorno:** `make run` (Docker: PostgreSQL 18 + API Go en `:8080`)  
**Método:** Chrome headless (`--headless=new`) controlado por CDP. Recorrido por clics reales, envío de formularios y verificación de códigos HTTP, redirecciones, consola del navegador y recursos cargados.

## Resumen ejecutivo

La navegación general del sitio funciona correctamente: las páginas públicas, el flujo de autenticación, el panel de administración y las operaciones con htmx responden como se espera. Se verificaron **64 enlaces internos únicos**, todos con estado **200**, sin excepciones JavaScript, imágenes rotas ni hojas de estilo faltantes.

Se detectaron **5 puntos a mejorar**, ninguno bloqueante: una inconsistencia de layout en la edición de Grandes Premios, una página 403 sin estilo, un favicon ausente en el 404 y dos detalles menores de redirección/menú.

## Alcance probado

| Bloque | Casos verificados |
|---|---|
| Público | `/`, `/escuderias`, `/pilotos`, `/gran_premios`, `/ingresar`, `/registrar`, ruta inexistente (404) |
| API REST | `/usuarios`, `/pilotos_historicos`, `/gran_premios_historicos` |
| Sesión | Registro, login válido e inválido, `/apuestas`, `/salir`, `/ingresar` con sesión activa |
| Admin | `/admin`, menú completo, formularios de programación/resultados, alta y edición de entidades |
| htmx | Orden de escuderías y pilotos, crear/modificar/eliminar apuesta |
| Enlaces | 64 `href` internos relevados desde las seis páginas principales, incluidos los parámetros `?id=` dinámicos |

## Resultados por perfil

| Perfil | Comportamiento |
|---|---|
| **Visitante anónimo** | Home, listados y formularios disponibles. `/apuestas` y `/admin` redirigen a `/ingresar`. Ruta inexistente devuelve 404 con página propia. |
| **Usuario registrado** | Tras el login accede a `/` con opciones “Apuestas” y “Cerrar sesión”. No ve controles de administración. `/admin` responde 403. |
| **Administrador** | El login redirige a `/admin`. El menú “Completar resultados / Pilotos / Escuderías / Gran Premio / Cerrar sesión” navega correctamente. Los formularios de edición y el flujo de apuestas funcionan. |

Comprobación adicional: el orden de escuderías cambia de `Alpine, Aston Martin, Audi…` a `Williams, Red Bull, Racing Bulls…` al activar el control htmx, sin recarga completa.

## Hallazgos y recomendaciones

| # | Severidad | Hallazgo | Ubicación | Recomendación |
|---|---|---|---|---|
| 1 | Media | `ModificarGranPremioPage` usa el layout público: pierde el menú admin y el footer apunta a `/`. Es la única pantalla de edición con este comportamiento. | `views/gran_premios.templ:131` | Reemplazar `@Layout(...)` por `@LayoutAdmin(...)`, igual que en `pilotos.templ` y `escuderias.templ`. |
| 2 | Media | Un usuario autenticado sin rol admin que entra a `/admin` recibe `No autorizado` en texto plano, sin layout ni forma de volver. | `handlers/admin.go` | Renderizar una página 403 con el layout estándar o redirigir al inicio con un aviso. |
| 3 | Baja | La página 404 no declara favicon, por lo que se genera una petición extra a `/favicon.ico` que también falla (error en consola). | `views/noEncontrado.templ` | Agregar `<link rel="icon" type="image/png" href="/front/imagenes/LogoPredictOne.png"/>`. |
| 4 | Baja | `LayoutAdmin` no incluye acceso a “Apuestas”; el admin llega a `/admin` y solo puede acceder a sus apuestas escribiendo la URL. | `views/layoutAdmin.templ` | Agregar el enlace a `/apuestas` si se considera parte del flujo del administrador. |
| 5 | Baja | Con sesión activa, visitar `/ingresar` redirige a `/` en lugar de `/admin` o `/apuestas`. | `handlers/usuario.go:41` | Redirigir según el rol del usuario, para coherencia con el login. |

## Aspectos positivos

- Los 64 enlaces internos relevados responden 200; no hay enlaces muertos.
- No se registraron errores de JavaScript ni fallos de carga de CSS, htmx o imágenes.
- Los controles de administración están correctamente ocultos para usuarios normales.
- El cierre de sesión invalida el acceso a rutas protegidas.
- Los formularios de login, registro, alta/edición y apuestas completan su flujo y redirigen correctamente.

## Notas

- Durante las pruebas se crearon un usuario temporal y una apuesta; ambos fueron eliminados al finalizar.
- La API quedó ejecutándose en `http://localhost:8080` al cierre del análisis.
- Evidencia cruda de las corridas: `/tmp/opencode/` (`public.json`, `flow.json`, `user.json`, `deep.json`, `links.json`, `make-run.log`).

## Entregables gráficos

- [Mapa de sitio y capturas](MAPA_SITIO.md) · [diagrama SVG](mapa_sitio.svg) · [diagrama PNG](mapa_sitio.png)
- [Capturas de pantalla](capturas/) — 13 PNG de 1440 px de ancho.

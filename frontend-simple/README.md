# Frontend Simple - Andrei Mes Manur API

## Descripción

Este es un frontend simple desarrollado en un solo archivo HTML que permite interactuar con todos los endpoints de la API de Andrei Mes Manur. Está diseñado para ser funcional y fácil de usar sin necesidad de frameworks complejos.

## Características

- **Un solo archivo**: Todo el frontend está contenido en `index.html`
- **Sin dependencias**: No requiere Node.js, npm, ni framework alguno
- **Interfaz completa**: Acceso a todos los endpoints de la API
- **Diseño responsive**: Funciona en desktop y móvil
- **Autenticación JWT**: Manejo completo del sistema de autenticación
- **Roles diferenciados**: Interfaces específicas para cada rol de usuario

## Estructura del Frontend

### 🎨 **Diseño Visual**
- Tema oscuro con colores inspirados en "imperio digital"
- Gradientes y efectos visuales modernos
- Layout responsive con grid y flexbox
- Iconos emoji para mejor UX

### 🔐 **Sistema de Autenticación**
- Login con email y contraseña
- Registro de nuevos usuarios
- Almacenamiento seguro del JWT token
- Información del usuario autenticado
- Función de logout

### 📋 **Organización por Pestañas**
1. **Públicos**: Endpoints accesibles sin autenticación
2. **Andrei**: Funcionalidades de administrador supremo
3. **Demon**: Herramientas para demonios
4. **Network Admin**: Opciones para administradores de red

## Funcionalidades por Rol

### 🌐 **Endpoints Públicos**
- Ver página de resistencia con todos los posts públicos

### 👑 **Andrei (Administrador Supremo)**
- **Gestión de Usuarios**:
  - Ver todos los usuarios
  - Ver usuario específico por ID
  - Eliminar usuarios
  
- **Estadísticas y Rankings**:
  - Ver estadísticas de la plataforma
  - Ver ranking de demonios
  
- **Gestión de Posts**:
  - Ver todos los posts
  - Eliminar posts
  - Crear posts de Andrei
  
- **Sistema de Recompensas**:
  - Crear recompensas y castigos para demonios
  - Asignar puntos positivos o negativos

### 👹 **Demon**
- **Gestión de Víctimas**:
  - Ver network admins disponibles
  - Ver víctimas asignadas
  - Registrar nuevas víctimas
  
- **Sistema de Reportes**:
  - Ver reportes creados
  - Crear nuevos reportes de actividad
  - Actualizar estado de reportes
  
- **Estadísticas**:
  - Ver estadísticas personales
  
- **Posts**:
  - Crear posts de propaganda

### 👨‍💻 **Network Admin**
- **Posts Anónimos**:
  - Crear posts anónimos para la resistencia

## Uso del Frontend

### 1. **Preparación**
```bash
# Asegúrate de que el backend esté corriendo en puerto 8080
go run main.go

# Abre el archivo HTML en tu navegador
# Puedes usar cualquier navegador moderno
```

### 2. **Autenticación**
Usa una de las cuentas de prueba incluidas:

| Rol | Email | Password |
|-----|-------|----------|
| 🔮 Andrei | `andrei@evil.com` | `password123` |
| 👹 Demon | `shadow@evil.com` | `password123` |
| 👨‍💻 Network Admin | `john.admin@company.com` | `password123` |

### 3. **Navegación**
1. Inicia sesión con una cuenta
2. Navega entre las pestañas según tu rol
3. Usa los formularios para interactuar con los endpoints
4. Las respuestas aparecen en el área de respuesta de cada sección

### 4. **Características de UX**
- **Respuestas JSON**: Todas las respuestas se muestran formateadas
- **Alertas**: Notificaciones de éxito y error
- **Validación**: Validación de formularios del lado cliente
- **Estado Visual**: Indicadores de usuario autenticado

## Endpoints Implementados

### Públicos
- `GET /api/v1/resistance` - Ver página de resistencia

### Autenticación
- `POST /api/v1/login` - Iniciar sesión
- `POST /api/v1/register` - Registrar usuario

### Andrei
- `GET /api/v1/admin/users` - Listar usuarios
- `GET /api/v1/admin/users/:id` - Ver usuario
- `DELETE /api/v1/admin/users/:id` - Eliminar usuario
- `GET /api/v1/admin/stats` - Estadísticas
- `GET /api/v1/admin/demons/ranking` - Ranking demonios
- `GET /api/v1/admin/posts` - Listar posts
- `DELETE /api/v1/admin/posts/:id` - Eliminar post
- `POST /api/v1/admin/posts` - Crear post
- `POST /api/v1/admin/rewards` - Crear recompensa

### Demon
- `GET /api/v1/demons/available-network-admins` - Network admins disponibles
- `GET /api/v1/demons/victims` - Ver víctimas
- `POST /api/v1/demons/victims` - Registrar víctima
- `GET /api/v1/demons/reports` - Ver reportes
- `POST /api/v1/demons/reports` - Crear reporte
- `PUT /api/v1/demons/reports/:id` - Actualizar reporte
- `GET /api/v1/demons/stats` - Estadísticas
- `POST /api/v1/demons/posts` - Crear post

### Network Admin
- `POST /api/v1/network-admins/posts/anonymous` - Crear post anónimo

## Arquitectura Técnica

### JavaScript Nativo
```javascript
// Cliente API centralizado
async function apiCall(method, endpoint, data = null) {
    // Configuración de headers con JWT
    // Manejo de errores
    // Formateo de respuestas
}
```

### Manejo de Estado
```javascript
let authToken = '';        // JWT token
let currentUser = null;    // Información del usuario
```

### Sistema de Respuestas
```javascript
function showResponse(elementId, data) {
    // Formateo JSON con sintaxis highlighting
    // Manejo de errores y éxitos
}
```

## Personalización

### Colores y Tema
```css
/* Variables CSS para fácil personalización */
:root {
    --primary-color: #f56565;
    --bg-primary: #1e1e2e;
    --bg-secondary: #2d3748;
}
```

### Añadir Nuevos Endpoints
1. Crear la función JavaScript correspondiente
2. Añadir el formulario HTML si es necesario
3. Conectar con el evento correspondiente
4. Manejar la respuesta en el área designada

## Ventajas de Este Approach

✅ **Simplicidad**: Un solo archivo, fácil de modificar
✅ **Portabilidad**: Funciona en cualquier servidor web
✅ **Debugging**: Fácil inspección en DevTools
✅ **Personalización**: CSS y JS modificables directamente
✅ **Sin Dependencies**: No requiere build process
✅ **Completo**: Cubre todos los endpoints de la API

## Limitaciones

❌ **Escalabilidad**: No ideal para aplicaciones muy grandes
❌ **Organización**: Todo en un archivo puede ser difícil de mantener
❌ **Testing**: Limitado para testing automatizado
❌ **SEO**: SPA básico sin optimizaciones SEO
❌ **Offline**: No tiene capacidades offline

## Mejoras Futuras Posibles

1. **Persistencia**: LocalStorage para formularios
2. **Validación Avanzada**: Validación más robusta
3. **UI/UX**: Animaciones y transiciones
4. **Accessibility**: Mejores prácticas de accesibilidad
5. **Mobile**: Optimizaciones específicas para móvil

## Soporte

Para usar este frontend:

1. Asegúrate de que el backend esté corriendo en `localhost:8080`
2. Abre `index.html` en cualquier navegador moderno
3. Usa las cuentas de prueba proporcionadas
4. Reporta cualquier issue o mejora necesaria

Este frontend simple te permite probar todos los endpoints de la API de manera rápida y eficiente, con una interfaz visual clara y funcional.
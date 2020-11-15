# Chat-Go


## Objetivo

Desarrollar una sala de chat donde los usuarios *(clientes)* que se conecten al servidor puedan enviar/recibir mensajes de texto así como archivos de los otros usuarios.

## Descripción

Implementar un cliente-servidor con las siguientes características:

- El servidor será el encargado (intermediario) de coordinar el envío de mensajes o archivos.
- El servidor tendrá las siguientes opciones:
    1. Mostrar los mensajes/nombre de los archivos enviados.
    2. Opción para respaldar en un archivo de texto los mensajes/nombre de los archivos enviados.
    3. Terminar servidor.

Los clientes tendrán las siguientes opciones:

- Enviar un mensaje de texto al servidor.
- Enviar un archivo (con opción para escribir la ubicación del archivo a enviar).
- Mostrar los mensajes recibidos del servidor (como si fuera una ventana de chat).

Hay que tener en cuenta lo siguiente:

- El servidor deberá de mantener una *lista* de clientes, los cuales en cualquier momento pueden enviar un mensaje de texto o un archivo. Con lo anterior, el servidor siempre estará *escuchando* por mensajes de entrada de los clientes.
- Cuando un cliente envíe una señal para desconectarse, este deberá ser liberado de la *lista.*

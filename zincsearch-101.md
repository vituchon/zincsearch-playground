# Primeros pasos con zincsearch

* Zincsearch esta diseñado para almacenar documentos de texto y que luego uno pueda buscar en ellos determinados textos conforme a determinadas reglas de búsqueda. 
* Zincsearch es un programa, que al ejecutarse queda "escuchando" en un puerto arbitrario (se suele usar el 4080) esperando peticiones http, exponiendo su funcionalidad a travéz de un API que sigue la convención REST. A su vez tmb dicho programa nos prove de un cliente web que podemos usar abriendo el navegador y yendo a http://localhost:4080/ui/search, esto nos abre una interfaz gráfica donde podemos hacer las API call para "jugar" con el programa (otra alternativa es usar programa como curl o bien programarse un cliente)
  *  Esto no es una tarea trivial, basta ponerse a pensar un escenario donde buscar a personas por su nombre y por ejemplo introduciendo "Marco Aurelio" deberia saltar como coincidencia una persona que se llame "Marco Garcia Aurelio", bueno posiblemente sea bueno que hagas el ejercicio para entender la motivación que hay detras de todo esto.
* Para instalarlo directamente te podés bajar el programa e instalarlo, siguiendo las instrucciones [acá](https://zincsearch-docs.zinc.dev/quickstart/).
  * Si están windows y están usando el git bash recuerden que el `set variable=valor` no funciona sino que tienen que usar `export variable=valor`, lo digo pues van a tener que establecer dos variables de entorno para poder ejecutar el programa, a saber: ZINC_FIRST_ADMIN_USER y ZINC_FIRST_ADMIN_PASSWORD, que son basicamente las credenciales de acceso.
* Los conceptos fundamentales son dos: Indices y Documentos (más transfondo al respecto [aquí](https://zincsearch-docs.zinc.dev/concepts)). 
  * Un documento es un objeto JSON que contienen datos por los cuales después se pueden buscar. Como por ejemplo ´{ "id": 1, "name": "Moises and the papas rock band" }´, en donde se podria luego por "1" del campo "id" y/o alguna combinación de las palabras que hay en el campo "name". Recordemos que los campos no tienen que ser de texto (como ocurre con el campo "id") pues todo tipo de dato puede ser llevado a un represtanción de texto (tal como Java afirmo hace tiempo en 1995/1996 al agregar toString al objeto Object).
  * Un indice almacena datos adicionales para que luego uno pueda buscar datos de los documentos que fuiron indexados con dicho indice.
  * El indice terminando siendo en alguna medida un agrupador de documentos, siendo algo que "envuelve" a los documentos dentro de un listado y además posibilitando que se hagan busquedas sobre dicho listado.
  * Esto en los ejemplos, quedea un poco confuso como se mezclan indices y documentos, pues existe una forma de agregar muchos documentos (bulk) con cierto formato en donde se aclara explicitamente el indice.
* Basicamente zincsearch puede tener varios listados de documentos. Cada listado podemos verlo como una tabla de postgres.
* La idea es que todos los documentos de un mismo índice sigan un mismo formato
  * Se puede crear primero el indice (configurandolo de forma fina para que la indexación resulte "mejor" o se adapte a ciertos parametros)
    * Sobre esto de crear el índice: No sabría decir si es una práctica recomemdada o no
  * Se puede directametne insertar documentos creando el indice con la primera inserción de forma "automática (zincsearch lo hace por nosotros)
    * En muchos ejemplos que hay ví que hacen esto.
      * Por ejemplo con hacer `curl http://localhost:4080/api/algun_indice/_doc -u admin:Complexpass#123  --data-binary '{ "id": 1, "name": "Moises and the papas rock band" }'` crea el indice "algun_indice" en caso de no existir. Y agrega el documento que pasé como argumento "data-binary".
      
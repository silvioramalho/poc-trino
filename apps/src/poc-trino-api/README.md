# Trino API Proxy

## TODO:

* Adicionar healthcheck

Adicionar endpoints de healthcheck


* Adicionar rodapé nos resultados das consultas:

```
pageNumber
pageSize 
totalPages
totalRecords
```

* Adicionar restrição na query a partir da segmentação do cliente

Nesse caso, os usuários que tiverem utilizando a API deverá ter em seu token as restrições de acesso aos registros do banco, forçando sempre a inclusão do seu ID ou ID acessíveis nas consultas.

* Adicionar filtro a partir dos campos disponíveis no payload

Possibilitar que o usuário faça envio de filtros que irão compor a restrição da query

Exemplo:

```
{
    "phone": {
        "value": ["33-464-151-3439"],
        "restriction": "AND"
    },
    "name": {
        "value": ["Customer#000000011"],
        "restriction": "OR"
    }
}
```
`É importante garantir que não haja SQL INJECTION`



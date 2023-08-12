# Dashboard Service

This service acts as a data access layer translating API requests into MongoDB pipelines. It manages caching, queuing and database connection. Based on the data model and incoming query the service generates the appropriate MongoDB aggregation pipeline. It queries the database and then sends the result back to the client.
This tool's idea and documentation is heavily inspired from https://cube.dev/. However, the whole source code has been written from scratch by me in golang.

# Block definition

Blocks are defined by YAML files. They represent a table of data

```json
cubes:
  - name: Users
    sql: SELECT * FROM USERS
    joins:
      - name: Organizations
        relationship: belongs_to
        local_key: organization_id
				foreign_key: id
    measures:
      - name: count
        type: count
        sql: id
    dimensions:
      - name: organization_id
        sql: organization_id
        type: number
        primary_key: true
      - name: created_at
        type: time
        sql: created_at
      - name: country
        type: string
        sql: country
```

The \***\*BLOCK\*\*** references the current block

Every block in schema can have its own **\*\***data_source**\*\***

## Query Properties

A Query has the following properties:

- `measures`: An array of measures.
- `dimensions`: An array of dimensions.
- `filters`: An array of objects, describing filters.
- `timeDimensions`: A convenient way to specify a time dimension with a filter.
  It is an array of objects in time dimension format
- `segments`: An array of segments. A segment is a named filter, created in the
  Data Schema.
- `limit`: A row limit for your query. The default value is `10000`. The maximum
  allowed limit is `50000`
- `offset`: The number of initial rows to be skipped for your query. The default
  value is `0`.
- `order`: An object, where the keys are measures or dimensions to order by and
  their corresponding values are either `asc` or `desc`. The order of the fields
  to order on is based on the order of the keys in the object.

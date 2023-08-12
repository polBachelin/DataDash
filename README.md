# Dashboard Service

This service acts as a data access layer translating API requests into MongoDB pipelines. It manages caching, queuing and database connection. Based on the data model and incoming query the service generates the appropriate MongoDB aggregation pipeline. It queries the database and then sends the result back to the client.
This tool's idea and documentation is heavily inspired from https://cube.dev/. However, the whole source code has been written from scratch by me in golang.

# Block definition

Blocks are defined by YAML files. They represent a table of data

```json
blocks:
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

The **BLOCK** references the current block

Every block in schema can have its own **data_source**

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

### Filter Operators

- `equals`: use when you need an exact match. If you use it on a dimension of type boolean the values will be ignored. Supports multiple values
- `notEquals`: The opposite of equals. Supports multiple values
- `contains`: Acts as a wildcard case-insensitive LIKE operator. Supports multiple values
- `notContains`: Opposite of contains. Supports multiple values
- `startsWith`: To check if a string starts with a substring. Supports multiple values
- `endsWith`: To check if a string ends with a substring. Supports multiple values
- `gt`: The greater than operator
- `gte`: The great than or equal operator
- `lt`: The less than operator
- `lte`: The less than or equal operator
- `set`: This operator checks if the value is not null
- `notSet`: This operator checks if the value is null

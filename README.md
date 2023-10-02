# Dashboard Service

This service acts as a data access layer translating API requests into MongoDB pipelines. It manages caching, queuing and database connection. Based on the data model and incoming query the service generates the appropriate MongoDB aggregation pipeline. It queries the database and then sends the result back to the client.

Afin que le service fonctionne correctement il faut specifier certaines variables d’environment

```json
SCHEMA_PATH : Le chemin vers la dossier qui contient les schemas en fichier YAML
DB_HOST : Le host de la base de données auquel le service se connectera
DB_PORT : Le port de la base de données auquel le service se connectera
DB_USERNAME : Le nom d'utilisateur de la base de données auquel le service se connectera
DB_PASS : Le mot de passe de la base de données auquel le service se connectera
DB_NAME : Le nom de la de la base de données auquel le service se connectera
API_PORT : Le port de l'api valeur par defaut : 8080
```

****\*\*****\*\*\*\*****\*\*****\*\*\*\*****\*\*****\*\*\*\*****\*\*****Connecter à une base de données****\*\*****\*\*\*\*****\*\*****\*\*\*\*****\*\*****\*\*\*\*****\*\*****

Endpoint: `POST /connect`

Body: JSON

- `db_host` _string_
- `db_port` _string_
- `db_username` _string_
- `db_pass` _string_
- `db_name` _string_

Response:

- `status`: 200 for success

**Envoyer une query pour récuperer de la données**

Endpoint: `POST /query`

Body: JSON

- `measures` Tableau de [measures](https://www.notion.so/Dashboard-Service-5a5b1a5b477e4ea5911b96bc23d07e9f?pvs=21)
- `dimension` Tableau de [dimensions](https://www.notion.so/Dashboard-Service-5a5b1a5b477e4ea5911b96bc23d07e9f?pvs=21)
- `filters` Tableau de [filtres](https://www.notion.so/Dashboard-Service-5a5b1a5b477e4ea5911b96bc23d07e9f?pvs=21)
- `timeDimensions` Tableau de [timeDimensions](https://www.notion.so/Dashboard-Service-5a5b1a5b477e4ea5911b96bc23d07e9f?pvs=21)
- `limit` Limite le nombre de réponses
- `offset` Point de départ au sein de la données à partir duquel les resultats sont récupérés. Par défaut à `0`
- `order` Un objet, où les clés sont des mesures ou des dimensions par lesquelles trier et leurs valeurs correspondantes sont soit `asc` (croissant) soit `desc` (décroissant). L'ordre des champs à trier est basé sur l'ordre des clés dans l'objet.
-

## Measures

The measures parameter contains a set of measures and each measure is an aggregation over a certain column in your database table

### Measure types

- string
  - This measure type allows defining measures as a `string` value.
- time
  - This measure type allows defining measures as a `time` value
- boolean
  - condense data into a single boolean value, returns “true” - “false” is **all** match the SQL
- number
  - Can take any valid SQL expression that results in a number or integer
- count
  - Performs a table count, similar to SQL’s `COUNT` function.
- sum
  - Adds up the values in a given field. It is similar to SQL’s `SUM` function.
- avg
  - Averages the values in a given field. It is similar to SQL’s AVG function.
- min
- max

## Dimensions

Refers to attributes or categorical data fields that are used for grouping and categorizing data in a multi-dimensional dataset or cube. These dimensions help organize and analyze data along specific categories, providing insights into different aspects of the dataset

## Filters

A filter is an object with the following properties:

- `member`: Dimension or measure to be used in the filter, for example:
  `stories.isDraft`. See below on difference between filtering dimensions vs
  filtering measures.
- `operator`: An operator to be used in the filter. Only some operators are
  available for measures. For dimensions the available operators depend on the
  type of the dimension. Please see the reference below for the full list of
  available operators.
- `values`: An array of values for the filter. Values must be of type String. If
  you need to pass a date, pass it as a string in `YYYY-MM-DD` format.

The following are the possible operators for filters

- equals : equals
- notEquals
- contains
- notContains
- startsWith
- endsWith
- gt : greater than
- gte : greater than or equals
- lt : less than
- lte : less than or equals
- set : is not null
- notSet : is null
- inDateRange
- notInDateRange
- beforeDate
- afterDate

```json
"member": "Sale.amount"
"operator": "gt"
"values": ["5000"]
```

## Time Dimensions

Since grouping and filtering by a time dimension is quite a common case,

- `dimension`: Time dimension name.
- `dateRange`: An array of dates with the following format `YYYY-MM-DD` or in `YYYY-MM-DDTHH:mm:ss.SSS` format. Values should always be local and in query `timezone`. Dates in `YYYY-MM-DD` format are also accepted. Such dates are padded to the start and end of the day if used in start and end of date range interval accordingly. Please note that for timestamp comparison, `>=` and `<=` operators are used. It requires, for example, that the end date range date `2020-01-01` is padded to `2020-01-01T23:59:59.999`. If only one date is specified it's equivalent to passing two of the same dates as a date range.
- `compareDateRange`: An array of date ranges to compare a measure change over previous period
- `granularity`: A granularity for a time dimension. It supports the following values `second`, `minute`, `hour`, `day`, `week`, `month`, `quarter`, `year`.

```json
{
  "measures": ["stories.count"],
  "timeDimensions": [
    {
      "dimension": "stories.time",
      "dateRange": ["2015-01-01", "2015-12-31"],
      "granularity": "month"
    }
  ]
}
```

## Example query

```json
//Input query
{
  "measures": ["Stories.count"],
  "dimensions": ["Stories.category"],
  "filters": [
    {
      "member": "Stories.isDraft",
      "operator": "equals",
      "values": ["No"]
    }
  ],
  "time_dimensions": [
    {
      "dimension": "Stories.time",
      "dateRange": ["2015-01-01", "2015-12-31"],
      "granularity": "month"
    }
  ],
  "limit": 100,
  "offset": 0,
  "order": {
    "Stories.time": "asc",
    "Stories.count": "desc"
  }
}
```

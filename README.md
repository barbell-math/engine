# powerlifting-engine

## Application Setup
This application requires PostgreSQL to be installed. Once installed, database(s) will need to be created based on what you are doing.
1. If you are just a casual user (If you are confused only create this database):
    1. ```createdb production```
1. If you intend to modify the application for development purposes, two additional databases will need to be created. Having two separate databases allows for the test cases to run in parallel.
    1. ```createdb dbTest```
    1. ```createdb modelTest```

Before running the application, two environmnent variables need to be set: ```DB_USER``` and ```DB_PSWD```. Both are expected to be strings.

### Todo
1. Add escape characters to CSV file splitter
1. Make CSVSplitter accept strings for delimiters
1. Make Join interface
1. Add determinant function to matrix
1. Look for ways to reduce error from matrix inversion/other calculations

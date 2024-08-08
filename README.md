# Quote Scraper

A simple Go application that scrapes quotes from [Quotes to Scrape](https://quotes.toscrape.com) and stores them into an AWS RDS MySQL database.

## Features

- Scrapes quotes, authors, and tags from multiple pages
- Inserts data into a MySQL database
- Configurable data source through `config.json`

## Setup

1. **Clone the Repository:**

    ```sh
    git clone https://github.com/kuldeepsingh666/GoScraper.git
    cd GoScraper
    ```

2. **Modify the `config.json` File:**

    ```json
    {
      "data_source": "user:password@tcp(rds_endpoint:3306)/quotes"
    }
    ```

3. **Initialize the Database:**

   Ensure the database and table are set up.Assuming the database you have is called quotes.

    ```mysql
    USE quotes;

   CREATE TABLE IF NOT EXISTS quotes (
   id INT AUTO_INCREMENT PRIMARY KEY,
   text TEXT NOT NULL,
   author VARCHAR(255) NOT NULL,
   tags TEXT
   );

    ```

4. **Run the Application:**

    ```sh
    go run main.go
    ```

## Dependencies

- Go (1.18+)
- `github.com/PuerkitoBio/goquery`
- `github.com/go-sql-driver/mysql`


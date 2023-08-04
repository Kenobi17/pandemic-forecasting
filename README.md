# Pandemic Forecasting 

Pandemic Forecasting is a project that uses web scraping and the ARIMA forecasting model to predict the number of new daily COVID-19 cases.

## Live Demo
For a live demo, you can visit this [page](https://pandemic-forecasting-production.up.railway.app/).

## Local Setup
Before you start, make sure you have Docker, Go, and Python installed on your system.

There are two options for setting up the project locally: Docker and manual setup.

### Docker
To run the project using Docker, execute the following commands:
```bash
docker build -t pandemic-forecasting .
docker run -p 8080:8080 pandemic-forecasting
```

### Manual
For a manual setup, follow these steps:

1. Set up a Python virtual environment. If you don't have the `virtualenv` package installed, you can install it with:

    ```bash  
    pip install virtualenv
    ```

Then, to create a new virtual environment named `env`, use:

    ```bash
    python3 -m venv env
    ```

Activate the virtual environment:
- On Unix or MacOS, use: `source env/bin/activate`
- On Windows, use: `.\env\Scripts\activate`

2. Install the necessary Python and Go packages:

    ```bash
    pip install numpy matplotlib statsmodels pandas
    go mod download
    ```
3. Build and run the application:

    ```bash
    go build -o main
    ./main
    ```
## Usage
After setting up and running the application, the COVID-19 statistics will be scraped from a webpage and parsed into a JSON file. This file is then read by a Python script, which generates a CSV file and a plot image. 

The application exposes several endpoints:

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | /api/forecast | Returns the forecast data in JSON format |
| GET | /api/forecast/fake | Returns a "fake" forecast for demonstration purposes in JSON format |
| GET | /api/forecast/download | Downloads the forecast CSV file |
| GET | /api/forecast/fake/download | Downloads the fake forecast CSV file |
| GET | / | Serves an HTML file that displays the forecast data on a chart |

## Example Responses
### GET - /api/forecast
```json
{
  "daily_cases": [
    0,
    0,
    0,
    0,
    0,
    0,
    0
  ],
  "dates": [
    "2023-08-04",
    "2023-08-05",
    "2023-08-06",
    "2023-08-07",
    "2023-08-08",
    "2023-08-09",
    "2023-08-10"
  ]
}
```

### GET - /api/forecast/fake
The fake forecast data was created to simulate a scenario where there have been recent reports of new COVID-19 cases, leading to a forecast that isn't just a constant zero. This also helps to generate a more interesting chart for demonstration purposes.
```json
{
  "daily_cases": [
    11,
    10,
    9,
    9,
    10,
    11,
    11
  ],
  "dates": [
    "2023-08-03",
    "2023-08-04",
    "2023-08-05",
    "2023-08-06",
    "2023-08-07",
    "2023-08-08",
    "2023-08-09"
  ]
}
```

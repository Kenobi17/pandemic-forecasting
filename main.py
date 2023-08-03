from utils.forecast_utils import get_data, get_forecast, csv_forecast, plot_forecast 

data_path = 'data.json'
file_name = 'forecast'
daily_cases, last_reported_date = get_data(data_path)
forecast = get_forecast(daily_cases)
csv_forecast(forecast, last_reported_date, file_name)
plot_forecast(forecast, last_reported_date, file_name)
print(forecast)

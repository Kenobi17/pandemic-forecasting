import csv
import json
import numpy as np
import matplotlib.pyplot as plt
import matplotlib.dates as mdates
from datetime import datetime, timedelta
from statsmodels.tsa.arima.model import ARIMA

def get_data(file_path):
    with open(file_path, 'r') as f:
        data = json.load(f)
    daily_cases = [x for x in data['daily_cases'] if x is not None]
    daily_cases = np.array(daily_cases)
    last_reported_date = datetime.strptime(data['date'], "%Y-%m-%d")
    return daily_cases, last_reported_date

def get_forecast(daily_cases):
    model = ARIMA(daily_cases, order=(2,1,2))
    fit = model.fit()
    forecast = fit.forecast(steps=7)
    forecast = [round(int(x)) for x in forecast]
    return forecast

def csv_forecast(forecast, last_reported_date, csv_name):
    start_date = last_reported_date + timedelta(days=1)
    dates = [start_date + timedelta(days=i) for i in range(7)]
    with open(csv_name + '.csv', 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(['date', 'new_cases'])
        for date, new_cases in zip(dates, forecast):
            writer.writerow([date.strftime('%Y-%m-%d'), new_cases])

def plot_forecast(forecast, last_reported_date, plot_name):
    start_date = last_reported_date + timedelta(days=1)
    dates = [start_date + timedelta(days=i) for i in range(7)]
    plt.plot(dates, forecast)
    plt.gca().xaxis.set_major_formatter(mdates.DateFormatter('%m-%d'))
    plt.gca().xaxis.set_major_locator(mdates.DayLocator())
    plt.yticks(np.arange(min(forecast), max(forecast) + 1 , 1))
    plt.xlabel('Date')
    plt.ylabel('New Cases')
    plt.title('Forecast of New Cases')
    plt.savefig(plot_name + '.png')

data_path = 'data.json'
file_name = 'forecast'
daily_cases, last_reported_date = get_data(data_path)
forecast = get_forecast(daily_cases)
csv_forecast(forecast, last_reported_date, file_name)
plot_forecast(forecast, last_reported_date, file_name)

fake_data_path = 'fake_data.json'
fake_file_name = 'fake_forecast'
fake_daily_cases, fake_last_reported_date = get_data(fake_data_path)
fake_forecast = get_forecast(fake_daily_cases)
csv_forecast(fake_forecast, fake_last_reported_date, fake_file_name)
plot_forecast(fake_forecast, fake_last_reported_date, fake_file_name)

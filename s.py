import json
import numpy as np
from statsmodels.tsa.arima.model import ARIMA


with open('data.json', 'r') as f:
    data = json.load(f)

data = [x for x in data if x is not None]
data = np.array(data)


model = ARIMA(data, order=(1,1,0))
model_fit = model.fit()

forecast = model_fit.forecast(steps=7)

print(forecast)

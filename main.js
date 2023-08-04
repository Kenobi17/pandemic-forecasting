function loadForecast(endpoint, tableId, chartId) {
    fetch(endpoint)
        .then(response => response.json())
        .then(data => {
            const tableBody = document.querySelector(`#${tableId} tbody`);
            const ctx = document.querySelector(`#${chartId}`).getContext('2d');

            const labels = data.dates;
            const dataSet = data.daily_cases;

            labels.forEach((date, i) => {
                const row = document.createElement('tr');
                const dateCell = document.createElement('td');
                const casesCell = document.createElement('td');
                dateCell.textContent = date;
                casesCell.textContent = dataSet[i];
                dateCell.className = 'px-4 py-2 text-center'; 
                casesCell.className = 'px-4 py-2 text-center'; 
                row.appendChild(dateCell);
                row.appendChild(casesCell);
                tableBody.appendChild(row);
            });

            new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Daily Cases',
                        data: dataSet,
                        fill: false,
                        borderColor: 'rgb(75, 192, 192)',
                        tension: 0.1
                    }]
                }
            });
        });
}

loadForecast('/api/forecast', 'forecastTable', 'forecastChart');
loadForecast('/api/forecast/fake', 'fakeForecastTable', 'fakeForecastChart');
